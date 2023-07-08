package handler

import (
	"github.com/non1996/autowire/anno"
)

var _ = anno.Annotations{
	anno.Component[TestHandler]{
		anno.Autowired{Field: "TestService"},
	},
}
