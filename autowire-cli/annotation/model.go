package annotation

import (
	"go/ast"
	"strconv"
	"strings"

	"github.com/modern-go/reflect2"
)

// BaseAnnotation 注解基类
type BaseAnnotation struct {
	Name     string     // 注解名称
	Generics []ast.Expr // 注解泛型参数列表
}

func (a *BaseAnnotation) GetName() string {
	return a.Name

}

// PrimaryAnnotation 根注解
type PrimaryAnnotation struct {
	BaseAnnotation
	Childrens []SecondaryAnnotation // 子注解列表
}

// SecondaryAnnotation 次级注解
type SecondaryAnnotation struct {
	BaseAnnotation
	Params []AnnotationParam // 注解参数列表
}

func (a *SecondaryAnnotation) GetParam(name string) ast.Expr {
	for _, param := range a.Params {
		if param.Key == name {
			return param.Value
		}
	}

	return nil
}

func (a *SecondaryAnnotation) GetStringParam(name string, defaultValue ...string) string {
	expr := a.GetParam(name)

	if reflect2.IsNil(expr) {
		if len(defaultValue) == 0 {
			panic(errMissingParam(a.Name, name))
		}
		return defaultValue[0]
	}

	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		panic(errParamValueType(a.Name, name, "string", expr))
	}

	return strings.Trim(lit.Value, "\\\"|`")
}

func (a *SecondaryAnnotation) GetBoolParam(name string, defaultValue ...bool) bool {
	expr := a.GetParam(name)

	if reflect2.IsNil(expr) {
		if len(defaultValue) == 0 {
			panic(errMissingParam(a.Name, name))
		}
		return defaultValue[0]
	}

	ident, ok := expr.(*ast.Ident)
	if !ok {
		panic(errParamValueType(a.Name, name, "bool", expr))
	}
	b, err := strconv.ParseBool(ident.Name)
	if err != nil {
		panic(err)
	}

	return b
}

// AnnotationParam 注解参数
type AnnotationParam struct {
	Key   string
	Value ast.Expr
}
