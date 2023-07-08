package client

type MQManager struct {
	mqTest string
	mqProd string
}

func (m *MQManager) Init() error {
	return nil
}

func (m *MQManager) TestMQ() *MQ {
	return &MQ{topic: m.mqTest}
}

func (m *MQManager) ProdMQ() *MQ {
	return &MQ{topic: m.mqTest}
}

type MQ struct {
	topic string
}

func (m *MQ) Init() error {
	return nil
}
