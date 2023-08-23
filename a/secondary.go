package a

type ComponentScans struct {
	Value []string
}

func (c ComponentScans) isSecondary() {
}

type Alias struct {
	Value string
}

func (a Alias) isSecondary() {
}

type ValueType struct {
}

func (v ValueType) isSecondary() {
}

type Implement[_ any] struct {
}

func (i Implement[_]) isSecondary() {
}

type Configuration struct {
}

func (c Configuration) isSecondary() {
}

type Primary struct {
}

func (p Primary) isSecondary() {
}

type ConditionalOnProperty struct {
	Value string
	Scope string
	Key   string
}

func (c ConditionalOnProperty) isSecondary() {
}

type Autowired struct {
	Value     any
	Qualifier string
	Required  bool
}

func (a Autowired) isSecondary() {
}

type Value struct {
	Value    string
	Scope    string
	Key      string
	Required bool
}

func (v Value) isSecondary() {
}

type Env struct {
	Value    any
	Key      string
	Default  string
	Required bool
}

func (e Env) isSecondary() {
}

type PostConstruct struct {
	Value any
}

func (p PostConstruct) isSecondary() {
}

type Bean struct {
	Alias  string
	Method any
}

func (b Bean) isSecondary() {
}

type PropertyProvider struct {
	Value string
	Scope string
}

func (p PropertyProvider) isSecondary() {
}
