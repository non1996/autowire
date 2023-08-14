package autowire

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
	return appFactory.build(ctx)
}

func (ctx *AppContext) getComponent(typ Type, require ...bool) any {
	typeName := getTypeNameT(typ)

	comps := ctx.components.listByTypeName(typeName)
	if len(comps) == 0 && required(require) {
		panic(errComponentNotFound(typeName))
	}
	if len(comps) == 1 {
		return ctx.getInstance(comps[0])
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
		return ctx.getInstance(primary)
	}

	if len(otherMatches) == 1 {
		return ctx.getInstance(otherMatches[0])
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

	return ctx.getInstance(comp)
}

func (ctx *AppContext) match(cond *Condition) bool {
	if cond == nil {
		return false
	}

	v, exist := ctx.properties.get(cond.Scope, cond.Key)
	if !exist {
		return false
	}

	s, ok := v.(string)
	if !ok {
		return false
	}
	return exist && cond.Value == s
}

func (ctx *AppContext) getInstance(component *component) any {
	if component.instance == nil {
		component.instance = component.factory.build(ctx)
	}

	return component.instance
}
