package autowire

import (
	"fmt"
	"reflect"

	"github.com/non1996/go-jsonobj/container"
	"github.com/non1996/go-jsonobj/function"

	"github.com/non1996/autowire/cli/model"
)

type typ struct {
	name    string
	primary bool
}

type Component struct {
	instance  any
	name      string
	typ       typ
	impls     []typ
	condition *model.Condition
}

type AppContext struct {
	factories  []Factory
	properties map[string]string
	container  map[string][]*Component
	nameMap    map[string]*Component
}

func NewAppContext(factories ...Factory) *AppContext {
	return &AppContext{
		factories:  factories,
		properties: map[string]string{},
		container:  map[string][]*Component{},
		nameMap:    map[string]*Component{},
	}
}

func Wire(ctx *AppContext) {
	for _, factory := range ctx.factories {
		factory.Build(ctx)
	}
}

func AddComponent(
	ctx *AppContext,
	component *Component,
) {
	if _, exist := ctx.nameMap[component.name]; exist {
		panic("")
	}
	ctx.nameMap[component.name] = component

	add(ctx, component.typ, component)
	for _, impl := range component.impls {
		add(ctx, impl, component)
	}
}

func add(ctx *AppContext, typ typ, component *Component) {
	comps := ctx.container[typ.name]

	if len(comps) != 0 && typ.primary && comps[0].typ.primary {
		panic("")
	}

	if typ.primary {
		comps = container.SliceConcatM([]*Component{component}, comps)
	} else {
		comps = append(comps, component)
	}

	ctx.container[typ.name] = comps
}

func GetComponent[T any](ctx *AppContext, require ...bool) T {
	typeName := getTypeName[T]()

	components := ctx.container[typeName]

	var (
		primary      *Component
		otherMatches []*Component
	)

	for _, comp := range components {
		if comp.typ.primary {
			primary = comp
		} else if match(ctx, comp.condition) {
			otherMatches = append(otherMatches, comp)
		}
	}

	if primary != nil {
		return primary.instance.(T)
	}

	if len(otherMatches) == 1 {
		return otherMatches[0].instance.(T)
	}

	if len(otherMatches) > 0 {

	}

	if len(require) != 0 && require[0] {
		panic("")
	}

	return function.Zero[T]()
}

func GetComponentByAlias[T any](ctx *AppContext, name string, require ...bool) T {
	// TODO
	return function.Zero[T]()
}

func match(ctx *AppContext, cond *model.Condition) bool {
	if cond == nil {
		return false
	}

	v, exist := ctx.properties[cond.Key]
	return exist && cond.Value == v
}

func getTypeName[T any]() string {
	return getTypeNameByType(reflect.TypeOf((*T)(nil)))
}

func getTypeNameByType(typ reflect.Type) string {
	return fmt.Sprintf("%s.%s", typ.Elem().PkgPath(), typ.Elem().Name())
}

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil))
}
