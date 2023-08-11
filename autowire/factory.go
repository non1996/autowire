package autowire

import (
	"reflect"

	"github.com/non1996/go-jsonobj/function"
)

type ConstructFunc[C any] func(*C) error

//type ConfigProviderFunc[C any] func(*C) any

type ConfigProvider[C any] struct {
	Scope    string
	Provider func(*C) any
}

type Condition struct {
	Key   string
	Value string
}

type Factory interface {
	name() string              // 组件名称
	cType() reflect.Type       // 组件类型
	impl() []reflect.Type      // 组件实现类型
	isPrimary() bool           // 同类型中是否是主要类型
	condition() *Condition     // 组件构造条件
	register(ctx *AppContext)  // 向 app context 注册
	build(ctx *AppContext) any // 构造组件、依赖注入、后置初始化
}

// ComponentFactory 组件工厂
type ComponentFactory[C any] struct {
	Name           string
	Ptr            bool
	Primary        bool
	Configuration  bool
	Implement      []reflect.Type
	Condition      *Condition
	FieldInjectors []Injector[C]
	PostConstruct  ConstructFunc[C]
	Configs        []ConfigProvider[C]
}

func (f ComponentFactory[C]) name() string {
	return f.Name
}

func (f ComponentFactory[C]) cType() reflect.Type {
	return TypeOf[C]()
}

func (f ComponentFactory[C]) impl() []reflect.Type {
	return f.Implement
}

func (f ComponentFactory[C]) isPrimary() bool {
	return f.Primary
}

func (f ComponentFactory[C]) condition() *Condition {
	return f.Condition
}

func (f ComponentFactory[C]) register(ctx *AppContext) {
	ctx.components.add(f)

	if f.Configuration {
		for _, conf := range f.Configs {
			ctx.properties.add(propertyProvider{
				scope: conf.Scope,
				provide: func() any {
					if f.Ptr {
						comp := GetComponent[*C](ctx)
						return conf.Provider(comp)
					} else {
						comp := GetComponent[C](ctx)
						return conf.Provider(&comp)
					}
				},
			})
		}
	}
}

func (f ComponentFactory[C]) build(ctx *AppContext) any {
	comp := new(C)

	// 依赖注入
	for _, fieldInjector := range f.FieldInjectors {
		fieldInjector.inject(ctx, comp)
	}

	// 执行后置操作
	if f.PostConstruct != nil {
		err := f.PostConstruct(comp)
		if err != nil {
			panic(err)
		}
	}

	return function.Ternary[any](f.Ptr, comp, *comp)
}

type BeanFactory[C any, B any] struct {
	Name          string
	ComponentName string
	BuildFunc     func(C) B
}

func (f BeanFactory[C, B]) name() string {
	return f.Name
}

func (f BeanFactory[C, B]) cType() reflect.Type {
	return TypeOf[B]()
}

func (f BeanFactory[C, B]) impl() []reflect.Type {
	return nil
}

func (f BeanFactory[C, B]) isPrimary() bool {
	return false
}

func (f BeanFactory[C, B]) isConfiguration() bool {
	return false
}

func (f BeanFactory[C, B]) condition() *Condition {
	return nil
}

func (f BeanFactory[C, B]) register(ctx *AppContext) {
	ctx.components.add(f)
}

func (f BeanFactory[C, B]) build(ctx *AppContext) any {
	comp := GetComponentByName[C](ctx, f.ComponentName)
	return f.BuildFunc(comp)
}

type ApplicationFactory[A any] struct {
	Factory
	App       *A
	Injectors []Injector[A]
}

func (a ApplicationFactory[A]) build(ctx *AppContext) any {
	// 依赖注入
	for _, fieldInjector := range a.Injectors {
		fieldInjector.inject(ctx, a.App)
	}

	return a.App
}

func Register(factories ...Factory) any {
	for _, factory := range factories {
		factory.register(defaultAppContext)
	}
	return nil
}
