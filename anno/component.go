package anno

// Component 标明是一个依赖注入组件
type Component[_ any] []any

type Name struct {
	Value string
}

type Implement[_ any] struct{}

type Primary struct{}

type ConditionalOnConfig struct {
	Key   string
	Value string
}

type Autowired struct {
	Field     any
	Qualifier string
	Required  bool
}

type Value struct {
	Field string
	Key   string
}

type PostConstruct struct {
	Value any
}
