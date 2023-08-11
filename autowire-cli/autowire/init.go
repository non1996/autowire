package autowire

import (
	"text/template"
)

func init() {
	factoryTemplate = template.Must(template.New("tmplGenFile").Parse(tmplGenFile))
	factoryTemplate = template.Must(factoryTemplate.New("tmplComponentFactory").Parse(tmplComponentFactory))
	factoryTemplate = template.Must(factoryTemplate.New("tmplInjector").Parse(tmplInjector))
	factoryTemplate = template.Must(factoryTemplate.New("tmplApp").Parse(tmplApp))
}
