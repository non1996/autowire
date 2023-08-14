package autowire

import (
	"fmt"
	"go/ast"

	"github.com/non1996/go-jsonobj/stream"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

func parseAnnoAlias(component *Component, a annotation.SecondaryAnnotation) {
	component.Alias = a.GetStringParam("Value")
}

func parseAnnoValueType(component *Component, _ annotation.SecondaryAnnotation) {
	component.Type.Ptr = false
}

func parseAnnoImplement(component *Component, a annotation.SecondaryAnnotation) {
	assert.Assert(len(a.Generics) == 1, "annotation <Implement> should has only one generic type")

	implType := a.Generics[0]

	component.Implements = append(component.Implements, parseType(implType))
}

func parseAnnoConfiguration(component *Component, _ annotation.SecondaryAnnotation) {
	component.IsConfiguration = true
}

func parseAnnoPrimary(component *Component, _ annotation.SecondaryAnnotation) {
	component.Primary = true
}

func parseAnnoCondition(component *Component, a annotation.SecondaryAnnotation) {
	component.Condition = &Condition{
		Scope: a.GetStringParam("Scope"),
		Key:   a.GetStringParam("Key"),
		Value: a.GetStringParam("Value"),
	}
}

func parseAnnoAutowired(component *Component, a annotation.SecondaryAnnotation) {
	assert.Assert(len(a.Generics) == 1)

	component.AddInjector(ComponentInjector{
		BaseInjector: parseBaseInjector(component, a),
		Type:         parseType(a.Generics[0]),
		Qualifier:    a.GetStringParam("Qualifier", ""),
	})
}

func parseAnnoValue(component *Component, a annotation.SecondaryAnnotation) {
	component.AddInjector(ValueInjector{
		BaseInjector: parseBaseInjector(component, a),
		Scope:        a.GetStringParam("Scope"),
		Key:          a.GetStringParam("Key"),
	})
}

func parseAnnoEnv(component *Component, a annotation.SecondaryAnnotation) {
	component.Injectors = append(component.Injectors, EnvInjector{
		BaseInjector: parseBaseInjector(component, a),
		Key:          a.GetStringParam("Key"),
		Default:      a.GetStringParam("Default", ""),
	})
}

func parseAnnoPostConstruct(component *Component, a annotation.SecondaryAnnotation) {
	component.PostConstruct = &PostConstruct{
		IsMethod: true,
		FuncName: a.GetStringParam("Value"),
	}
}

func parseAnnoBean(component *Component, a annotation.SecondaryAnnotation) {
	component.Beans = append(component.Beans, Bean{
		Type:   parseType(a.Generics[0]),
		Alias:  mustStringLit(a.GetParam("Alias")),
		Method: a.GetStringParam("Method"),
	})
}

func parseAnnoConfigurations(app *Application, a annotation.SecondaryAnnotation) {
	app.Configurations = stream.Map(a.GetParam("Value").(*ast.CompositeLit).Elts,
		func(t1 ast.Expr) string {
			return mustStringLit(t1)
		})
}

func parseAnnoConfig(component *Component, a annotation.SecondaryAnnotation) {
	component.Properties = append(component.Properties, PropertyProvider{
		Field: a.GetStringParam("Field"),
		Scope: a.GetStringParam("Scope"),
	})
}

func parseBaseInjector(component *Component, a annotation.SecondaryAnnotation) BaseInjector {
	return BaseInjector{
		FieldName: a.GetStringParam("Field"),
		CompType:  component.Type,
		Required:  a.GetBoolParam("Required", true),
	}
}

func errInvalidAnnotation(anno string) error {
	return fmt.Errorf("invalid component annotation <%s>", anno)
}
