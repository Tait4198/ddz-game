package main

// 用于指定调用方法
type MessageType uint

type MessageFunc func(ClientMessage)

type InstructionFunc func(string)

// 客户端发出消息
type ClientMessage struct {
	Level   string      `json:"level"`
	Type    MessageType `json:"type"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
}

// 服务端接收消息
type SendMessage struct {
	Level   string      `json:"level"`
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}
