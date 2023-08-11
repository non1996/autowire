package annotation

import (
	"fmt"
	"go/ast"
)

func errInvalidAnnoNameType(expr ast.Expr) error {
	return fmt.Errorf("invalid annotation name's syntax, expect <SelectorExpr>, <Ident>. but <%T>", expr)
}

func errInvalidAnnoType(expr ast.Expr) error {
	return fmt.Errorf("invalid annotation's syntax, expect <IndexListExpr>, <IndexExpr>, <SelectorExpr>. but <%T>", expr)
}

func errMissingParam(annoName, fieldName string) error {
	return fmt.Errorf("[annotation] <%s>, missing field <%s>", annoName, fieldName)
}

func errParamValueType(annoName, fieldName string, typeName string, astExpr ast.Expr) error {
	return fmt.Errorf("[annotation] <%s>, expect field <%s>'s type is %s but ast expr is %T",
		annoName, fieldName, typeName, astExpr)
}
