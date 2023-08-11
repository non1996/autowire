package autowire

type Injector interface {
	Kind() string
}

type BaseInjector struct {
	FieldName string // 字段名称
	CompType  Type   // 组件类型
	Required  bool
}

type ComponentInjector struct {
	BaseInjector
	Type      Type
	Qualifier string
}

func (i ComponentInjector) Kind() string {
	return "Component"
}

type ValueInjector struct {
	BaseInjector
	Scope string
	Key   string
}

func (i ValueInjector) Kind() string {
	return "Value"
}

type EnvInjector struct {
	BaseInjector
	Key     string
	Default string
}

func (i EnvInjector) Kind() string {
	return "Env"
}
