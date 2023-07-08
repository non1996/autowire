package dal

import (
	"fmt"

	"github.com/non1996/autowire/example/client"
)

type ADao interface {
	A()
}

type ADaoImpl struct {
	db client.DB
}

func (A *ADaoImpl) A() {
	fmt.Println("[ADaoImpl]")
}

type ADaoMock struct {
}

func (A *ADaoMock) A() {
	fmt.Println("[ADaoMock]")
}

type BDao interface {
	B()
}

type BDaoImpl struct {
	db client.DB
}

func (B *BDaoImpl) B() {
	fmt.Println("[BDaoImpl]")
}
