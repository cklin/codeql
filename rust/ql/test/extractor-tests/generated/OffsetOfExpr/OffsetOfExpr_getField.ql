// generated by codegen
import codeql.rust.elements
import TestUtils

from OffsetOfExpr x, int index
where toBeTested(x) and not x.isUnknown()
select x, index, x.getField(index)
