package main

import (
	"github.com/non1996/autowire/anno"
)

var _ = anno.Annotations{
	anno.Application[Application]{
		anno.ComponentScan{
			Value: []string{"example"},
		},
	},
}
