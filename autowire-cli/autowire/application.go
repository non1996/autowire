package autowire

import (
	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Application struct {
	Component
	Configurations []string
}

func parseApplication(annotation annotation.PrimaryAnnotation) (a *Application) {
	assert.Assert(len(annotation.Generics) == 1)
	a = &Application{}
	a.Type = parseType(annotation.Generics[0])

	for _, child := range annotation.Childrens {
		name := child.GetName()
		parser, ok := applicationAnnoParser[child.GetName()]
		if !ok {
			panic(errInvalidAnnotation(name))
		}
		parser(a, child)
	}

	return a
}

var applicationAnnoParser = map[string]func(*Application, annotation.SecondaryAnnotation){
	"Autowired": func(app *Application, a annotation.SecondaryAnnotation) {
		parseAnnoAutowired(&app.Component, a)
	},
	"Value": func(app *Application, a annotation.SecondaryAnnotation) {
		parseAnnoValue(&app.Component, a)
	},
	"Env": func(app *Application, a annotation.SecondaryAnnotation) {
		parseAnnoEnv(&app.Component, a)
	},
	"Configurations": parseAnnoConfigurations,
}
