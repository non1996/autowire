package autowire

import (
	"path/filepath"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Application struct {
	Component
	ComponentScan []string
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
	"PostConstruct": func(app *Application, a annotation.SecondaryAnnotation) {
		parseAnnoPostConstruct(&app.Component, a)
	},
	"ComponentScans": parseAnnoComponentScans,
}

func (p *Package) evaluateApplication(annotation annotation.PrimaryAnnotation) (a *Application) {
	assert.Assert(len(annotation.Generics) == 1)

	a = &Application{}
	a.Package = p
	a.Type = p.NewType(annotation.Generics[0])

	for _, child := range annotation.Childrens {
		name := child.GetName()
		parser, ok := applicationAnnoParser[child.GetName()]
		if !ok {
			panic(errInvalidAnnotation(name))
		}
		parser(a, child)
	}

	for _, cs := range a.ComponentScan {
		dirs := traversalDir(filepath.Join(p.Root, cs))

		for _, dir := range dirs {
			pkg := genContext.getPackage(dir)

			if pkg == nil || pkg.AnnoFile == nil {
				continue
			}

			pkgPath := filepath.Join(pkg.Module, pkg.PackagePath)

			if p.GenFile.Imports.GetByPath(pkgPath) == nil {
				p.GenFile.Imports.Add(&Import{
					Alias:         "_",
					PackagePath:   pkgPath,
					ExplicitAlias: true,
				})
			}
		}
	}

	return a
}
