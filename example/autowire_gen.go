package main

import (
	"github.com/non1996/autowire/autowire"
	"github.com/non1996/autowire/example/client"
	"github.com/non1996/autowire/example/dal"
)

func (a *Application) autowire() {
	a.AppContext = autowire.NewAppContext(
		client.DBFactory,
		dal.ADaoImplFactory,
		dal.ADaoMockFactory,
		dal.BDaoImplFactory,
	)

	autowire.Wire(a.AppContext)
}
