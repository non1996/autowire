package autowire

import (
	"os"
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
	Qualifier string
	Required  bool
	InjectFn  func(*C, D)
}

func (f ComponentInjector[C, D]) inject(ctx *AppContext, comp *C) {
	var dep D

	if f.Qualifier != "" {
		dep = GetComponentByName[D](ctx, f.Qualifier, f.Required)
	} else {
		dep = GetComponent[D](ctx, f.Required)
	}

	f.InjectFn(comp, dep)
}

// ValueInjector 值注入器
type ValueInjector[C any] struct {
	Scope    string
	Key      string
	Required bool
	InjectFn func(*C, any)
}

func (i ValueInjector[C]) inject(ctx *AppContext, comp *C) {
	value, exist := ctx.properties.get(i.Scope, i.Key)
	if !exist && i.Required {
		panic("")
	}

	i.InjectFn(comp, value)
}

// EnvInjector 环境变量注入器
type EnvInjector[C any] struct {
	Key          string
	Required     bool
	DefaultValue string
	InjectFn     func(*C, string)
}

func (i EnvInjector[C]) inject(ctx *AppContext, comp *C) {
	ev, exist := os.LookupEnv(i.Key)
	if !exist && i.Required {
		panic("")
	}

	if i.DefaultValue != "" {
		i.InjectFn(comp, i.DefaultValue)
	}

	i.InjectFn(comp, ev)
}
