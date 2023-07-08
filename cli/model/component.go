package model

type Component struct {
	Name          string      // 组件名，默认为全类型名
	TypeName      Identifier  // 类型名
	Implement     *Identifier // 实现的接口
	Primary       bool        // 是否是主要
	Condition     Condition
	Injectors     map[string]FieldInjector
	PostConstruct string
}

type FieldInjector struct {
	Type      string
	FieldName string
	Qualifier string
	Key       string
	Required  bool
}
