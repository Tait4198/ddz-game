package main

import cm "com.github/gc-common"

type MessageFunc func(ClientMessage)

type CommandFunc func(string)

// 客户端发出消息
type ClientMessage struct {
	Level   cm.MessageLevel `json:"level"`
	Type    cm.MessageType  `json:"type"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
}

// 服务端接收消息
type SendMessage struct {
	Level   cm.MessageLevel `json:"level"`
	Type    cm.MessageType  `json:"type"`
	Message string          `json:"message"`
}

type GameMessage struct {
	Message     string            `json:"message"`
	MessageType cm.DdzMessageType `json:"type"`
}
