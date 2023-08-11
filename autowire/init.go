package autowire

import (
	"github.com/non1996/go-jsonobj/function"
)

var defaultAppContext = NewAppContext()

func Context() *AppContext {
	return defaultAppContext
}

func GetComponent[T any](ctx *AppContext, require ...bool) T {
	typeName := getTypeName[T]()

	nodes := ctx.components.listByTypeName(typeName)
	if len(nodes) == 0 && required(require) {
		panic(errComponentNotFound(typeName))
	}
	if len(nodes) == 1 {
		return getInstance[T](ctx, nodes[0])
	}

	var (
		primary      *component
		otherMatches []*component
	)

	for _, node := range nodes {
		if node.factory.isPrimary() {
			primary = node
		} else if match(ctx, node.factory.condition()) {
			otherMatches = append(otherMatches, node)
		}
	}

	if len(otherMatches) == 1 {
		return getInstance[T](ctx, otherMatches[0])
	}

	if primary != nil {
		return getInstance[T](ctx, primary)
	}

	if len(otherMatches) > 1 {
		panic(errMultiMatch)
	}

	if len(otherMatches) == 0 && required(require) {
		panic(errComponentNotFound(typeName))
	}

	return function.Zero[T]()
}

func GetComponentByName[T any](ctx *AppContext, name string, require ...bool) T {
	node := ctx.components.getByName(name)
	if node == nil && required(require) {
		panic(errComponentNotFound(name))
	}

	return getInstance[T](ctx, node)
}
