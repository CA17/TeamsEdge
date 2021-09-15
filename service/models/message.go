package models

// 定义一些消息模型

type EdgeBootstrapMessage struct {
	Eid string `json:"eid"`
}

type EdgeInformMessage struct {
	Eid string `json:"eid"`
}

type HostpingReqMessage struct {
	Eid string `json:"eid"`
}

type HostpingRepMessage struct {
	Eid string `json:"eid"`
}

type EdgeMesage struct {
	Topic string      `json:"topic"`
	Body  interface{} `json:"body"`
}
