// generated by codegen, do not edit
import codeql.rust.elements
import TestUtils

from ParenPat x, string hasPat
where
  toBeTested(x) and
  not x.isUnknown() and
  if x.hasPat() then hasPat = "yes" else hasPat = "no"
select x, "hasPat:", hasPat
