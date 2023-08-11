package autowire

import (
	"github.com/non1996/go-jsonobj/stream"
)

// 组件容器，包含组件工厂对象factory和组件实例化后的instance
type component struct {
	factory  Factory
	instance any
}

// 组件集合
type components struct {
	list           []*component
	typeIndex      map[string][]int // 类型名 -> 组件下标列表
	nameIndex      map[string]int   // 组件名称 -> 组件下标
	configurations []int            // configuration类的下标列表
}

func newComponents() components {
	return components{
		typeIndex: map[string][]int{},
		nameIndex: map[string]int{},
	}
}

// 注册工厂对象
func (c *components) add(f Factory) {
	c.list = append(c.list, &component{factory: f})

	idx := len(c.list) - 1
	c.nameIndex[f.name()] = idx

	typeName := getTypeNameT(f.cType())
	c.typeIndex[typeName] = append(c.typeIndex[typeName], idx)

	impls := f.impl()
	for _, impl := range impls {
		implName := getTypeNameT(impl)
		c.typeIndex[implName] = append(c.typeIndex[implName], idx)
	}
}

// 根据组件名获取组件
func (c *components) getByName(name string) *component {
	idx, exist := c.nameIndex[name]
	if exist {
		return c.list[idx]
	}

	return nil
}

// 根据类型名获取该类型的所有组件
func (c *components) listByTypeName(typeName string) []*component {
	idxes := c.typeIndex[typeName]

	return stream.Map(idxes, func(idx int) *component { return c.list[idx] })
}

func getInstance[T any](ctx *AppContext, component *component) T {
	if component.instance == nil {
		component.instance = component.factory.build(ctx)
	}

	return cast[T](component.instance)
}
