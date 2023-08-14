package autowire

var defaultAppContext = NewAppContext()

func Context() *AppContext {
	return defaultAppContext
}

func GetComponent[T any](ctx *AppContext, require ...bool) T {
	c := ctx.getComponent(TypeOf[T](), require...)
	return cast[T](c)
}

func GetComponentByName[T any](ctx *AppContext, name string, require ...bool) T {
	c := ctx.getComponentByName(name, require...)
	return cast[T](c)
}

func Register(factories ...Factory) any {
	for _, factory := range factories {
		defaultAppContext.components.add(factory)
		factory.onRegister(defaultAppContext)
	}
	return nil
}
