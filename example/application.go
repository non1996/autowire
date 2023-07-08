package main

import (
	"github.com/non1996/autowire/autowire"
	"github.com/non1996/autowire/example/dal"
)

type Application struct {
	*autowire.AppContext
}

func (a *Application) Run() {
	aDao := autowire.GetComponent[dal.ADao](a.AppContext)
	aDao.A()

	bDao := autowire.GetComponent[dal.BDao](a.AppContext)
	bDao.B()
}
