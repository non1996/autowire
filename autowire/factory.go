package autowire

import (
	"reflect"

	"github.com/non1996/go-jsonobj/stream"

	"github.com/non1996/autowire/cli/model"
)

type PostConstructFunc[C any] func(*C) error

type Factory interface {
	Build(ctx *AppContext)
}

type Type struct {
	Type    reflect.Type
	Primary bool
}

type ComponentFactory[C any] struct {
	Name           string
	Implements     []Type
	Primary        bool
	Ptr            bool
	Condition      *model.Condition
	FieldInjectors []FieldInjector[C]
	PostConstruct  PostConstructFunc[C]
}

func (f ComponentFactory[C]) Build(ctx *AppContext) {
	comp := new(C)

	for _, fieldInjector := range f.FieldInjectors {
		fieldInjector.Inject(ctx, comp)
	}

	if f.PostConstruct != nil {
		err := f.PostConstruct(comp)
		if err != nil {
			panic(err)
		}
	}

	impls := stream.Map(f.Implements, func(t Type) typ {
		return typ{
			name:    getTypeNameByType(t.Type),
			primary: t.Primary,
		}
	})

	typeName := getTypeName[C]()
	if f.Name == "" {
		f.Name = typeName
	}

	if f.Ptr {
		AddComponent(ctx, &Component{
			instance:  comp,
			name:      f.Name,
			typ:       typ{name: typeName, primary: f.Primary},
			impls:     impls,
			condition: f.Condition,
		})
	} else {
		AddComponent(ctx, &Component{
			instance:  *comp,
			name:      f.Name,
			typ:       typ{name: typeName, primary: f.Primary},
			impls:     impls,
			condition: f.Condition,
		})

	}
}
