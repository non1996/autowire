package dal

import (
	"github.com/non1996/autowire/autowire"
	"github.com/non1996/autowire/cli/model"
	"github.com/non1996/autowire/example/client"
)

var (
	ADaoImplFactory = autowire.ComponentFactory[ADaoImpl]{
		Name: "",
		Implements: []autowire.Type{
			{
				Type:    autowire.TypeOf[ADao](),
				Primary: true,
			},
		},
		Primary:   true,
		Ptr:       true,
		Condition: nil,
		FieldInjectors: []autowire.FieldInjector[ADaoImpl]{
			autowire.FieldInjectorImpl[ADaoImpl, client.DB]{
				Qualifier: "",
				Require:   true,
				InjectFn: func(a *ADaoImpl, db client.DB) {
					a.db = db
				},
			},
		},
		PostConstruct: nil,
	}

	ADaoMockFactory = autowire.ComponentFactory[ADaoMock]{
		Name: "",
		Implements: []autowire.Type{
			{
				Type:    autowire.TypeOf[ADao](),
				Primary: false,
			},
		},
		Primary: false,
		Ptr:     true,
		Condition: &model.Condition{
			Key:   "Mock.Dao",
			Value: "1",
		},
		FieldInjectors: nil,
		PostConstruct:  nil,
	}

	BDaoImplFactory = autowire.ComponentFactory[BDaoImpl]{
		Name: "",
		Implements: []autowire.Type{
			{
				Type:    autowire.TypeOf[BDao](),
				Primary: false,
			},
		},
		Primary:   true,
		Ptr:       true,
		Condition: nil,
		FieldInjectors: []autowire.FieldInjector[BDaoImpl]{
			autowire.FieldInjectorImpl[BDaoImpl, client.DB]{
				Qualifier: "",
				Require:   true,
				InjectFn: func(b *BDaoImpl, db client.DB) {
					b.db = db
				},
			},
		},
		PostConstruct: nil,
	}
)
