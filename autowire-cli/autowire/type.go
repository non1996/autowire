package autowire

import (
	"fmt"
	"go/ast"
)

type Type struct {
	Ptr     bool
	Package string
	Name    string
	expr    ast.Expr
}

func (t Type) NameWithPkg() string {
	if t.Package == "" {
		return t.Name
	}

	return t.Package + "." + t.Name
}

func (t Type) NameComplete() string {
	if t.Ptr {
		return "*" + t.NameWithPkg()
	}

	return t.NameWithPkg()
}

func parseType(expr ast.Expr) (t Type) {
	t.expr = expr
	realType := expr

	if t2, ok := realType.(*ast.StarExpr); ok {
		t.Ptr = true
		realType = t2.X
	}

	switch n := realType.(type) {
	case *ast.SelectorExpr:
		t.Package = mustIdent(n.X).Name
		t.Name = n.Sel.Name
	case *ast.Ident:
		t.Name = n.Name
	default:
		panic(errInvalidTypeExpr(expr))
	}

	return t
}

func errInvalidTypeExpr(expr ast.Expr) error {
	return fmt.Errorf("invalid type expr, expect one of <SelectorExpr>, <Ident>, but <%T>", expr)
}
