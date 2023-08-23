package autowire

import (
	"go/token"
	"text/template"
)

var genContext GenerateContext

type GenerateContext struct {
	fSet     *token.FileSet
	packages []*Package
}

func (c *GenerateContext) getPackage(absolutePath string) *Package {
	for _, p := range c.packages {
		if p.AbsolutePath == absolutePath {
			return p
		}
	}

	return nil
}

func init() {
	factoryTemplate = template.Must(template.New("tmplGenFile").Parse(tmplGenFile))
	factoryTemplate = template.Must(factoryTemplate.New("tmplComponentFactory").Parse(tmplComponentFactory))
	factoryTemplate = template.Must(factoryTemplate.New("tmplInjector").Parse(tmplInjector))
	factoryTemplate = template.Must(factoryTemplate.New("tmplApp").Parse(tmplApp))
}
