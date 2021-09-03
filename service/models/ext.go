package models

// EdgeMessage 边缘节点与 ACS 交互消息基础模型
type EdgeMessage struct {
	Command string                 `json:"command"` // 消息命令字
	Code    int                    `json:"code"`    // 状态码
	Reason  string                 `json:"reason"`  // 消息内容
	Message map[string]interface{} `json:"message"` // 扩展属性
}

