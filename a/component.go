package a

// Component 标明是一个依赖注入组件
type Component[_ any] []any

type Alias struct {
	Value string
}

type ValueType struct{}

type Implement[_ any] struct{}

type Primary struct{}

type ConditionalOnProperty struct {
	Scope string
	Key   string
	Value string
}

type Autowired[_ any] struct {
	Field     any
	Qualifier string
	Required  bool
}

type Value struct {
	Field    string
	Scope    string
	Key      string
	Required bool
}

type PostConstruct struct {
	Value any
}

type Configuration struct {
	Format string
}

type Bean[T any] struct {
	Alias  string
	Method any
}

type PropertyProvider struct {
	Field string
	Scope string
}

type Env struct {
	Field    any
	Key      string
	Required bool
	Default  string
}
