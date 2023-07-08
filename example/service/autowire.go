package service

import (
	"github.com/non1996/autowire/anno"
)

var _ = anno.Annotations{
	anno.Component[TestServiceImpl]{
		anno.Name{Value: "TestServiceImpl"},
		anno.Implement[TestService]{},
		anno.Primary{},
		anno.ConditionalOnConfig{Key: "MockConfig.Mock", Value: "1"},
		anno.Autowired{Field: "ADao"},
		anno.Autowired{Field: "BDao", Required: false},
		anno.Autowired{Field: "MQ", Qualifier: "TestMQ"},
		anno.Value{Field: "A", Key: "GlobalConfig.A"},
		anno.PostConstruct{Value: "Construct"},
	},
}
