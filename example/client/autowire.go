package client

import (
	"github.com/non1996/autowire/anno"
)

var _ = anno.Annotations{
	anno.Configuration[MQManager]{
		anno.Value{Field: "mqTest", Key: "Client.MQ.Test"},
		anno.Value{Field: "mqProd", Key: "Client.MQ.Prod"},
		anno.PostConstruct{Value: "Init"},
		anno.Bean{Method: "TestMQ", InitMethod: "Init"},
		anno.Bean{Method: "ProdMQ", InitMethod: "Init"},
	},
}
