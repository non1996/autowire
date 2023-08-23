package autowire

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type Type struct {
	*token.FileSet
	ast.Expr
	i string
}

func (p *Package) NewType(expr ast.Expr) Type {
	return Type{
		FileSet: p.fset,
		Expr:    expr,
	}
}

func (t *Type) SetImport(i string) *Type {
	t.i = i
	return t
}

func (t *Type) TypeName() string {
	var expr = t.Expr

	if t.isPointer() {
		expr = typeDeref(t.Expr)
	}

	if t.isSlice() {
		expr = typeSliceElem(t.Expr)
	}

	return outputExpr(t.FileSet, expr)
}

func (t *Type) TypeNameComplete() string {
	return outputExpr(t.FileSet, t.Expr)
}

func (t *Type) isPointer() bool {
	_, ok := t.Expr.(*ast.StarExpr)
	return ok
}

func (t *Type) notPointer() bool {
	return !t.isPointer()
}

func (t *Type) getPackage() string {
	expr := typeDeref(t.Expr)
	selectorExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	return selectorExpr.X.(*ast.Ident).Name
}

func (t *Type) isThisPackage() bool {
	_, ok := typeDeref(t.Expr).(*ast.SelectorExpr)
	return !ok
}

func (t *Type) setPackage(alias string) {
	typeDeref(t.Expr).(*ast.SelectorExpr).X.(*ast.Ident).Name = alias
}

func (t *Type) isSlice() bool {
	_, ok := typeDeref(t.Expr).(*ast.ArrayType)
	return ok
}

func (t *Type) sliceElem() Type {
	return Type{
		FileSet: t.FileSet,
		Expr:    t.Expr.(*ast.ArrayType).Elt,
		i:       t.i,
	}
}

func outputExpr(fset *token.FileSet, expr ast.Expr) string {
	buf := bytes.NewBuffer(nil)
	err := format.Node(buf, fset, expr)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func typeDeref(expr ast.Expr) ast.Expr {
	for {
		if t, ok := expr.(*ast.StarExpr); ok {
			expr = t.X
		} else {
			break
		}
	}

	return expr
}

func typeSliceElem(expr ast.Expr) ast.Expr {
	for {
		if t, ok := expr.(*ast.ArrayType); ok {
			expr = t.Elt
		} else {
			break
		}
	}

	return expr
}
