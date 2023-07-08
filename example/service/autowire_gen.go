package service

import (
	"github.com/non1996/autowire/example/client"
	"github.com/non1996/autowire/example/dal"
)

func AutowireTestServiceImplInjector(
	aDao dal.ADao,
	bDao dal.BDao,
	mQ *client.MQ,
) {

}
