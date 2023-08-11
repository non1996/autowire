package autowire

import (
	"fmt"
)

type AppContext struct {
	components components
	properties properties
}

func NewAppContext() *AppContext {
	return &AppContext{
		components: newComponents(),
		properties: newProperties(),
	}
}

func (ctx *AppContext) WithFactory(factory Factory) *AppContext {
	factory.register(ctx)
	return ctx
}

func (ctx *AppContext) WithProperty(scope, key string, value any) *AppContext {
	ctx.properties.set(scope, key, value)
	return ctx
}

func (ctx *AppContext) Inject(appFactory Factory) any {
	return appFactory.build(ctx)
}

func match(ctx *AppContext, cond *Condition) bool {
	if cond == nil {
		return false
	}

	v, exist := ctx.properties.m[cond.Key]
	return exist && cond.Value == v
}

func required(r []bool) bool {
	return len(r) == 0 || r[0]
}

func cast[T any](v any) T {
	v2, ok := v.(T)
	if !ok {
		panic(fmt.Errorf("type cast failed, source type [%T] destination type [%s] are not compatible",
			v, getTypeName[T]()))
	}

	return v2
}
