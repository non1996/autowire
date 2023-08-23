package autowire

import (
	"path"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Component struct {
	*Package
	Type            Type   // 类型名
	Alias           string // 组件别名
	Ptr             bool   // 是否是指针类型
	Primary         bool   // 是否是主要
	IsConfiguration bool   // 是否是configuration
	Implements      []Type // 实现的接口
	Condition       *Condition
	Injectors       []Injector
	PostConstruct   *PostConstruct
	Beans           []*Bean
	Properties      []*PropertyProvider
}

func (c *Component) AddInjector(injector Injector) {
	c.Injectors = append(c.Injectors, injector)
}

type Condition struct {
	Value string
	Scope string
	Key   string
}

type Bean struct {
	Alias  string // 别名
	Type   Type   // 类型
	Ptr    bool
	Method string // 方法
}

type PostConstruct struct {
	MethodName   string
	HasErrorResp bool
}

type PropertyProvider struct {
	Type     Type
	Scope    string
	Value    string
	IsMethod bool
}

var componentAnnoParser = map[string]func(*Component, annotation.SecondaryAnnotation){
	"Alias":                 parseAnnoAlias,
	"ValueType":             parseAnnoValueType,
	"Implement":             parseAnnoImplement,
	"Configuration":         parseAnnoConfiguration,
	"Primary":               parseAnnoPrimary,
	"ConditionalOnProperty": parseAnnoConditionalOnProperty,
	"Autowired":             parseAnnoAutowired,
	"Value":                 parseAnnoValue,
	"Env":                   parseAnnoEnv,
	"PostConstruct":         parseAnnoPostConstruct,
	"Bean":                  parseAnnoBean,
	"PropertyProvider":      parseAnnoPropertyProvider,
}

func (p *Package) evaluateComponent(annotation annotation.PrimaryAnnotation) (c *Component) {
	// Component注解只能有一个泛型参数
	assert.Assert(len(annotation.Generics) == 1)

	c = &Component{
		Package: p,
		Type:    p.NewType(annotation.Generics[0]),
	}

	// Component中的类型不能是指针类型
	assert.Assert(c.Type.notPointer())
	// Component中的类型必须和注解相同包
	assert.Assert(c.Type.isThisPackage())

	c.Ptr = true

	c.Alias = path.Join(p.Module, p.PackagePath+"."+c.Type.TypeName())

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
