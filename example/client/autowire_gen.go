package client

import (
	"github.com/non1996/autowire/autowire"
)

var (
	DBFactory = autowire.ComponentFactory[DB]{
		Name:           "",
		Implements:     nil,
		Primary:        true,
		Ptr:            false,
		Condition:      nil,
		FieldInjectors: nil,
		PostConstruct:  nil,
	}
)
