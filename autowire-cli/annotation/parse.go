package annotation

import (
	"go/ast"
)

// Parse 解析注解列表
// annotations 都是如下形式
// var _ = a.Annotations{...}
func Parse(pkgAlias string, genDecl *ast.GenDecl) []PrimaryAnnotation {
	if len(genDecl.Specs) != 1 {
		return nil
	}

	value, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok {
		return nil
	}
	if len(value.Values) != 1 {
		return nil
	}

	// anno.Annotations{} 这一元素
	annoValue, ok := value.Values[0].(*ast.CompositeLit)
	if !ok {
		return nil
	}
	annoType, ok := annoValue.Type.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	if annoType.X == nil || annoType.Sel == nil {
		return nil
	}
	annoType2, ok := annoType.X.(*ast.Ident)
	if !ok {
		return nil
	}
	if annoType2.Name != pkgAlias && annoType.Sel.Name != "Annotations" {
		return nil
	}

	var annotations []PrimaryAnnotation
	for _, annoElem := range annoValue.Elts {
		compositeLit := annoElem.(*ast.CompositeLit)
		annotation := PrimaryAnnotation{
			BaseAnnotation: parseBaseAnnotation(compositeLit.Type),
		}

		for _, subAnnoElem := range compositeLit.Elts {
			subCompositeLit := subAnnoElem.(*ast.CompositeLit)
			subAnnotation := SecondaryAnnotation{
				BaseAnnotation: parseBaseAnnotation(subCompositeLit.Type),
			}

			for _, elt := range subCompositeLit.Elts {
				kv := elt.(*ast.KeyValueExpr)
				subAnnotation.Params = append(subAnnotation.Params, AnnotationParam{
					Key:   kv.Key.(*ast.Ident).Name,
					Value: kv.Value,
				})
			}

			annotation.Childrens = append(annotation.Childrens, subAnnotation)
		}

		annotations = append(annotations, annotation)
	}

	return annotations
}

func parseBaseAnnotation(expr ast.Expr) BaseAnnotation {
	var (
		name     string
		generics []ast.Expr
	)

	switch typ := expr.(type) {
	case *ast.IndexExpr: // 单泛型参数
		name = parseBaseAnnotationName(typ.X)
		generics = []ast.Expr{typ.Index}
	case *ast.IndexListExpr: // 多泛型参数
		name = parseBaseAnnotationName(typ.X)
		generics = typ.Indices
	case *ast.SelectorExpr: // 不带泛型参数
		name = parseBaseAnnotationName(typ.Sel)
	default:
		panic(errInvalidAnnoType(expr))
	}

	return BaseAnnotation{
		Name:     name,
		Generics: generics,
	}
}

// 解析注解名称
// 从语法 yyy.XXX[...]，拿到 yyy.XXX 这部分
func parseBaseAnnotationName(name ast.Expr) string {
	switch n := name.(type) {
	case *ast.SelectorExpr:
		return n.Sel.Name
	case *ast.Ident:
		return n.Name
	default:
		panic(errInvalidAnnoNameType(name))
	}
}
