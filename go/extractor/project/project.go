package project

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/github/codeql-go/extractor/diagnostics"
	"github.com/github/codeql-go/extractor/util"
	"golang.org/x/mod/semver"
)

func checkDirsNested(inputDirs []string) (string, bool) {
	// replace "." with "" so that we can check if all the paths are nested
	dirs := make([]string, len(inputDirs))
	for i, inputDir := range inputDirs {
		if inputDir == "." {
			dirs[i] = ""
		} else {
			dirs[i] = inputDir
		}
	}
	// the paths were generated by a depth-first search so I think they might
	// be sorted, but we sort them just in case
	sort.Strings(dirs)
	for _, dir := range dirs {
		if !strings.HasPrefix(dir, dirs[0]) {
			return "", false
		}
	}
	return dirs[0], true
}

// Returns the directory to run the go build in and whether to use a go.mod
// file.
func getBuildRoot(emitDiagnostics bool) (baseDir string, useGoMod bool) {
	goModPaths := util.FindAllFilesWithName(".", "go.mod", "vendor")
	if len(goModPaths) == 0 {
		baseDir = "."
		useGoMod = false
		return
	}
	goModDirs := util.GetParentDirs(goModPaths)
	if util.AnyGoFilesOutsideDirs(".", goModDirs...) {
		if emitDiagnostics {
			diagnostics.EmitGoFilesOutsideGoModules(goModPaths)
		}
		baseDir = "."
		useGoMod = false
		return
	}
	if len(goModPaths) > 1 {
		// currently not supported
		baseDir = "."
		commonRoot, nested := checkDirsNested(goModDirs)
		if nested && commonRoot == "" {
			useGoMod = true
		} else {
			useGoMod = false
		}
		if emitDiagnostics {
			if nested {
				diagnostics.EmitMultipleGoModFoundNested(goModPaths)
			} else {
				diagnostics.EmitMultipleGoModFoundNotNested(goModPaths)
			}
		}
		return
	}
	if emitDiagnostics {
		if goModDirs[0] == "." {
			diagnostics.EmitSingleRootGoModFound(goModPaths[0])
		} else {
			diagnostics.EmitSingleNonRootGoModFound(goModPaths[0])
		}
	}
	baseDir = goModDirs[0]
	useGoMod = true
	return
}

// DependencyInstallerMode is an enum describing how dependencies should be installed
type DependencyInstallerMode int

const (
	// GoGetNoModules represents dependency installation using `go get` without modules
	GoGetNoModules DependencyInstallerMode = iota
	// GoGetWithModules represents dependency installation using `go get` with modules
	GoGetWithModules
	// Dep represent dependency installation using `dep ensure`
	Dep
	// Glide represents dependency installation using `glide install`
	Glide
)

// Returns the appropriate DependencyInstallerMode for the current project
func getDepMode(emitDiagnostics bool) (DependencyInstallerMode, string) {
	bazelPaths := util.FindAllFilesWithName(".", "BUILD", "vendor")
	bazelPaths = append(bazelPaths, util.FindAllFilesWithName(".", "BUILD.bazel", "vendor")...)
	if len(bazelPaths) > 0 {
		// currently not supported
		if emitDiagnostics {
			diagnostics.EmitBazelBuildFilesFound(bazelPaths)
		}
	}

	goWorkPaths := util.FindAllFilesWithName(".", "go.work", "vendor")
	if len(goWorkPaths) > 0 {
		// currently not supported
		if emitDiagnostics {
			diagnostics.EmitGoWorkFound(goWorkPaths)
		}
	}

	baseDir, useGoMod := getBuildRoot(emitDiagnostics)
	if useGoMod {
		log.Println("Found go.mod, enabling go modules")
		return GoGetWithModules, baseDir
	}

	if util.FileExists("Gopkg.toml") {
		if emitDiagnostics {
			diagnostics.EmitGopkgTomlFound()
		}
		log.Println("Found Gopkg.toml, using dep instead of go get")
		return Dep, "."
	}

	if util.FileExists("glide.yaml") {
		if emitDiagnostics {
			diagnostics.EmitGlideYamlFound()
		}
		log.Println("Found glide.yaml, using Glide instead of go get")
		return Glide, "."
	}
	return GoGetNoModules, "."
}

// ModMode corresponds to the possible values of the -mod flag for the Go compiler
type ModMode int

const (
	ModUnset ModMode = iota
	ModReadonly
	ModMod
	ModVendor
)

// argsForGoVersion returns the arguments to pass to the Go compiler for the given `ModMode` and
// Go version
func (m ModMode) ArgsForGoVersion(version string) []string {
	switch m {
	case ModUnset:
		return []string{}
	case ModReadonly:
		return []string{"-mod=readonly"}
	case ModMod:
		if !semver.IsValid(version) {
			log.Fatalf("Invalid Go semver: '%s'", version)
		}
		if semver.Compare(version, "v1.14") < 0 {
			return []string{} // -mod=mod is the default behaviour for go <= 1.13, and is not accepted as an argument
		} else {
			return []string{"-mod=mod"}
		}
	case ModVendor:
		return []string{"-mod=vendor"}
	}
	return nil
}

// Returns the appropriate ModMode for the current project
func getModMode(depMode DependencyInstallerMode, baseDir string) ModMode {
	if depMode == GoGetWithModules {
		// if a vendor/modules.txt file exists, we assume that there are vendored Go dependencies, and
		// skip the dependency installation step and run the extractor with `-mod=vendor`
		if util.FileExists(filepath.Join(baseDir, "vendor", "modules.txt")) {
			return ModVendor
		} else if util.DirExists(filepath.Join(baseDir, "vendor")) {
			return ModMod
		}
	}
	return ModUnset
}

type BuildInfo struct {
	DepMode DependencyInstallerMode
	ModMode ModMode
	BaseDir string
}

func GetBuildInfo(emitDiagnostics bool) BuildInfo {
	depMode, baseDir := getDepMode(true)
	modMode := getModMode(depMode, baseDir)
	return BuildInfo{depMode, modMode, baseDir}
}

type GoVersionInfo struct {
	// The version string, if any
	Version string
	// A value indicating whether a version string was found
	Found bool
}

// Tries to open `go.mod` and read a go directive, returning the version and whether it was found.
func TryReadGoDirective(buildInfo BuildInfo) GoVersionInfo {
	if buildInfo.DepMode == GoGetWithModules {
		versionRe := regexp.MustCompile(`(?m)^go[ \t\r]+([0-9]+\.[0-9]+(\.[0-9]+)?)$`)
		goMod, err := os.ReadFile(filepath.Join(buildInfo.BaseDir, "go.mod"))
		if err != nil {
			log.Println("Failed to read go.mod to check for missing Go version")
		} else {
			matches := versionRe.FindSubmatch(goMod)
			if matches != nil {
				if len(matches) > 1 {
					return GoVersionInfo{string(matches[1]), true}
				}
			}
		}
	}
	return GoVersionInfo{"", false}
}
