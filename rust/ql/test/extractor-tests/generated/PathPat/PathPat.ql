// generated by codegen
import codeql.rust.elements
import TestUtils

from PathPat x, Path getPath
where
  toBeTested(x) and
  not x.isUnknown() and
  getPath = x.getPath()
select x, "getPath:", getPath
