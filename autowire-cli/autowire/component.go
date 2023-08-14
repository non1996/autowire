package autowire

import (
	"path"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Component struct {
	Type            Type   // 类型名
	Alias           string // 组件别名
	Primary         bool   // 是否是主要
	IsConfiguration bool   // 是否是configuration
	Implements      []Type // 实现的接口
	Condition       *Condition
	Injectors       []Injector
	PostConstruct   *PostConstruct
	Beans           []Bean
	Properties      []PropertyProvider
}

func (c *Component) AddInjector(injector Injector) {
	c.Injectors = append(c.Injectors, injector)
}

type Condition struct {
	Scope string
	Key   string
	Value string
}

type Bean struct {
	Alias  string // 别名
	Type   Type   // 类型
	Method string // 方法
}

type PostConstruct struct {
	IsMethod bool
	FuncName string
}

type PropertyProvider struct {
	Scope string
	Field string
}

func parseComponent(rootModule, relativePath string, annotation annotation.PrimaryAnnotation) (c *Component) {
	assert.Assert(len(annotation.Generics) == 1)

	c = &Component{}
	c.Type = parseType(annotation.Generics[0])
	c.Type.Ptr = true

	if c.Type.Package != "" {
		c.Alias = c.Type.Package + "." + c.Type.Name
	} else {
		c.Alias = path.Join(rootModule, relativePath+"."+c.Type.Name)
	}

	for _, child := range annotation.Childrens {
		name := child.GetName()
		parser, ok := componentAnnoParser[child.GetName()]
		if !ok {
			panic(errInvalidAnnotation(name))
		}
		parser(c, child)
	}

	return c
}

var componentAnnoParser = map[string]func(*Component, annotation.SecondaryAnnotation){
	"Alias":                 parseAnnoAlias,
	"ValueType":             parseAnnoValueType,
	"Implement":             parseAnnoImplement,
	"Configuration":         parseAnnoConfiguration,
	"Primary":               parseAnnoPrimary,
	"ConditionalOnProperty": parseAnnoProperty,
	"Autowired":             parseAnnoAutowired,
	"Value":                 parseAnnoValue,
	"Env":                   parseAnnoEnv,
	"PostConstruct":         parseAnnoPostConstruct,
	"Bean":                  parseAnnoBean,
	"PropertyProvider":      parseAnnoPropertyProvider,
}
