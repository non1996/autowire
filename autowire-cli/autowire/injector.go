package autowire

type InjectorKind string

const (
	InjectorKindComponent InjectorKind = "component"
	InjectorKindValue     InjectorKind = "value"
	InjectorKindEnv       InjectorKind = "env"
)

type Injector interface {
	Base() *BaseInjector
	Kind() InjectorKind
}

type BaseInjector struct {
	Value    string // 字段名称/函数名称
	CompType Type   // 组件类型
	IsMethod bool   // 是否使用函数进行注入
	IsSlice  bool   // 是否注入的是一个列表
	Required bool
}

type ComponentInjector struct {
	BaseInjector
	Type      Type
	Qualifier string
}

func (i *ComponentInjector) Base() *BaseInjector {
	return &i.BaseInjector
}

func (i *ComponentInjector) Kind() InjectorKind {
	return InjectorKindComponent
}

type ValueInjector struct {
	BaseInjector
	Scope string
	Key   string
}

func (i *ValueInjector) Base() *BaseInjector {
	return &i.BaseInjector
}

func (i *ValueInjector) Kind() InjectorKind {
	return InjectorKindValue
}

type EnvInjector struct {
	BaseInjector
	Key     string
	Default string
}

func (i *EnvInjector) Base() *BaseInjector {
	return &i.BaseInjector
}

func (i *EnvInjector) Kind() InjectorKind {
	return InjectorKindEnv
}
