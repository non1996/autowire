package autowire

import (
	"fmt"

	"github.com/non1996/go-jsonobj/stream"
)

type AppContext struct {
	components           components
	properties           properties
	environmentVariables environmentVariables
}

func NewAppContext() *AppContext {
	return &AppContext{
		components:           newComponents(),
		properties:           newProperties(),
		environmentVariables: newEnvironmentVariables(),
	}
}

func (ctx *AppContext) Inject(appFactory Factory) any {
	log.Debug("[AppContext] inject start")
	return appFactory.build(ctx)
}

func (ctx *AppContext) getComponent(typ Type, require ...bool) any {
	typeName := getTypeNameT(typ)

	comps := ctx.components.listByTypeName(typeName)
	if len(comps) == 0 && required(require) {
		panic(errComponentNotFound(typeName))
	}
	if len(comps) == 1 {
		return comps[0].getInstance(ctx)
	}

	var (
		primary      *component
		otherMatches []*component
	)

	for _, comp := range comps {
		if comp.factory.isPrimary() {
			primary = comp
		} else if ctx.match(comp.factory.condition()) {
			otherMatches = append(otherMatches, comp)
		}
	}

	if primary != nil {
		return primary.getInstance(ctx)
	}

	if len(otherMatches) == 1 {
		return otherMatches[0].getInstance(ctx)
	}

	if len(otherMatches) > 1 {
		panic(errMultiMatch)
	}

	if len(otherMatches) == 0 && required(require) {
		panic(errComponentNotFound(typeName))
	}

	return nil
}

func (ctx *AppContext) getComponentByName(name string, require ...bool) any {
	comp := ctx.components.getByName(name)
	if comp == nil && required(require) {
		panic(errComponentNotFound(name))
	}

	return comp.getInstance(ctx)
}

func (ctx *AppContext) listComponent(typ Type) []any {
	typeName := getTypeNameT(typ)

	comps := ctx.components.listByTypeName(typeName)

	return stream.Map(comps, func(comp *component) any {
		return comp.getInstance(ctx)
	})
}

func (ctx *AppContext) match(cond *Condition) bool {
	if cond == nil {
		return false
	}

	v, exist := ctx.properties.get(cond.Scope, cond.Key)
	if !exist {
		return false
	}

	s := fmt.Sprintf("%+v", v)
	return exist && cond.Value == s
}
