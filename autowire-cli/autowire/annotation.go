package autowire

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/non1996/go-jsonobj/stream"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

func (p *Package) parseAnnoFile() {
	astFile := p.ast.Files[filepath.Join(p.AbsolutePath, p.AutowireFileName)]
	if astFile == nil {
		return
	}

	imports := parseImport(astFile.Imports)
	i := imports.GetByPath("github.com/non1996/go-autowire/a")
	if i == nil {
		return
	}

	p.AnnoFile = &AnnoFile{Package: p}

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		p.AnnoFile.Annotations = append(p.AnnoFile.Annotations, annotation.Parse(i.Alias, genDecl)...)
	}
}

func parseAnnoComponentScans(app *Application, a annotation.SecondaryAnnotation) {
	app.ComponentScan = stream.Map(a.GetParam("Value").(*ast.CompositeLit).Elts,
		func(t1 ast.Expr) string {
			return mustStringLit(t1)
		})
}

func parseAnnoAlias(component *Component, a annotation.SecondaryAnnotation) {
	component.Alias = a.GetStringParam("Value")
}

func parseAnnoValueType(component *Component, _ annotation.SecondaryAnnotation) {
	component.Ptr = false
}

func parseAnnoImplement(component *Component, a annotation.SecondaryAnnotation) {
	assert.Assert(len(a.Generics) == 1, "annotation <Implement> should has only one generic type")

	implType := component.NewType(a.Generics[0])

	assert.Assert(implType.notPointer(), "annotation <Implement> should not have pointer generic type")

	component.Implements = append(component.Implements, implType)
}

func parseAnnoConfiguration(component *Component, _ annotation.SecondaryAnnotation) {
	component.IsConfiguration = true
}

func parseAnnoPrimary(component *Component, _ annotation.SecondaryAnnotation) {
	component.Primary = true
}

func parseAnnoConditionalOnProperty(component *Component, a annotation.SecondaryAnnotation) {
	component.Condition = &Condition{
		Value: a.GetStringParam("Value"),
		Scope: a.GetStringParam("Scope"),
		Key:   a.GetStringParam("Key"),
	}
}

func parseAnnoAutowired(component *Component, a annotation.SecondaryAnnotation) {
	component.AddInjector(&ComponentInjector{
		BaseInjector: parseBaseInjector(component, a),
		Qualifier:    a.GetStringParam("Qualifier", ""),
	})
}

func parseAnnoValue(component *Component, a annotation.SecondaryAnnotation) {
	component.AddInjector(&ValueInjector{
		BaseInjector: parseBaseInjector(component, a),
		Scope:        a.GetStringParam("Scope"),
		Key:          a.GetStringParam("Key"),
	})
}

func parseAnnoEnv(component *Component, a annotation.SecondaryAnnotation) {
	component.AddInjector(&EnvInjector{
		BaseInjector: parseBaseInjector(component, a),
		Key:          a.GetStringParam("Key"),
		Default:      a.GetStringParam("Default", ""),
	})
}

func parseAnnoPostConstruct(component *Component, a annotation.SecondaryAnnotation) {
	component.PostConstruct = &PostConstruct{
		MethodName: a.GetStringParam("Value"),
	}
}

func parseAnnoBean(component *Component, a annotation.SecondaryAnnotation) {
	component.Beans = append(component.Beans, &Bean{
		Alias:  a.GetStringParam("Alias", ""),
		Method: a.GetStringParam("Method"),
	})
}

func parseAnnoPropertyProvider(component *Component, a annotation.SecondaryAnnotation) {
	component.Properties = append(component.Properties, &PropertyProvider{
		Value: a.GetStringParam("Value"),
		Scope: a.GetStringParam("Scope"),
	})
}

func parseBaseInjector(component *Component, a annotation.SecondaryAnnotation) BaseInjector {
	return BaseInjector{
		Value:    a.GetStringParam("Value"),
		CompType: component.Type,
		Required: a.GetBoolParam("Required", true),
	}
}

func errInvalidAnnotation(anno string) error {
	return fmt.Errorf("invalid component annotation <%s>", anno)
}
