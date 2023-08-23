package autowire

import (
	"github.com/non1996/go-jsonobj/stream"
)

// Injector 注入器
// 注意：不管组件是以值类型注册还是以指针类型注册，泛型参数 C 都是组件的值类型
type Injector[C any] interface {
	inject(*AppContext, *C)
}

// ComponentInjector 组件注入器
// 注意：
// * 不管组件是以值类型注册还是以指针类型注册，泛型参数 C 都是组件的值类型
// * D 是被依赖组件的类型，可能是指针
type ComponentInjector[C any, D any] struct {
	Qualifier     string
	Required      bool
	IsSlice       bool
	InjectFn      func(*C, D)
	InjectSliceFn func(*C, []D)
}

func (f ComponentInjector[C, D]) inject(ctx *AppContext, comp *C) {
	log.Debug("[ComponentInjector] type <%s>, qualifier: %s, required: %t", getTypeName[D](), f.Qualifier, f.Required)

	if f.IsSlice {
		deps := ctx.listComponent(TypeOf[D]())

		f.InjectSliceFn(comp, stream.Map(deps, func(d any) D { return cast[D](d) }))
	} else {
		var dep any

		if f.Qualifier != "" {
			dep = ctx.getComponentByName(f.Qualifier, f.Required)
		} else {
			dep = ctx.getComponent(TypeOf[D](), f.Required)
		}

		f.InjectFn(comp, cast[D](dep))
	}

}

// ValueInjector 值注入器
type ValueInjector[C any] struct {
	Scope    string
	Key      string
	Required bool
	InjectFn func(*C, any)
}

func (i ValueInjector[C]) inject(ctx *AppContext, comp *C) {
	log.Debug("[ValueInjector] scope %s, key %s, required %t", i.Scope, i.Key, i.Required)

	value, exist := ctx.properties.get(i.Scope, i.Key)
	if !exist && i.Required {
		panic(errValueNotFound(i.Scope, i.Key))
	}

	i.InjectFn(comp, value)
}

// EnvInjector 环境变量注入器
type EnvInjector[C any] struct {
	Key          string
	DefaultValue string
	Required     bool
	InjectFn     func(*C, string)
}

func (i EnvInjector[C]) inject(ctx *AppContext, comp *C) {
	log.Debug("[EnvInjector] key %s, default value %s, required %t", i.Key, i.DefaultValue, i.Required)

	ev := ctx.environmentVariables.get(i.Key, i.DefaultValue, i.Required)

	i.InjectFn(comp, ev)
}
