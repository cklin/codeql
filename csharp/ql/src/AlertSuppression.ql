/**
 * @name Alert suppression
 * @description Generates information about alert suppressions.
 * @kind alert-suppression
 * @id cs/alert-suppression
 */

import csharp

/**
 * An alert suppression comment.
 */
class SuppressionComment extends CommentLine {
  string annotation;

  SuppressionComment() {
    // Must be either `// ...` or `/* ... */` on a single line.
    this.getRawText().regexpMatch("//.*|/\\*.*\\*/") and
    exists(string text | text = this.getText() |
      // match `lgtm[...]` anywhere in the comment
      annotation = text.regexpFind("(?i)\\blgtm\\s*\\[[^\\]]*\\]", _, _)
      or
      // match `lgtm` at the start of the comment and after semicolon
      annotation = text.regexpFind("(?i)(?<=^|;)\\s*lgtm(?!\\B|\\s*\\[)", _, _).trim()
    )
  }

  /** Gets the suppression annotation in this comment. */
  string getAnnotation() { result = annotation }

  /**
   * Holds if this comment applies to the range from column `startcolumn` of line `startline`
   * to column `endcolumn` of line `endline` in file `filepath`.
   */
  predicate covers(string filepath, int startline, int startcolumn, int endline, int endcolumn) {
    this.getLocation().hasLocationInfo(filepath, startline, _, endline, endcolumn) and
    startcolumn = 1
  }

  /** Gets the scope of this suppression. */
  SuppressionScope getScope() { this = result.getSuppressionComment() }
}

/**
 * The scope of an alert suppression comment.
 */
class SuppressionScope extends @commentline instanceof SuppressionComment {
  /** Gets a suppression comment with this scope. */
  SuppressionComment getSuppressionComment() { result = this }

  /**
   * Holds if this element is at the specified location.
   * The location spans column `startcolumn` of line `startline` to
   * column `endcolumn` of line `endline` in file `filepath`.
   * For more information, see
   * [Locations](https://codeql.github.com/docs/writing-codeql-queries/providing-locations-in-codeql-queries/).
   */
  predicate hasLocationInfo(
    string filepath, int startline, int startcolumn, int endline, int endcolumn
  ) {
    super.covers(filepath, startline, startcolumn, endline, endcolumn)
  }

  /** Gets a textual representation of this element. */
  string toString() { result = "suppression range" }
}

from SuppressionComment c
select c, // suppression comment
  c.getText(), // text of suppression comment (excluding delimiters)
  c.getAnnotation(), // text of suppression annotation
  c.getScope() // scope of suppression
