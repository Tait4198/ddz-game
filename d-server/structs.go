package main

import cm "com.github/gc-common"

// 客户端id
type ClientId string

// 房间id
type RoomId uint

type CenterMessageFunc func(ServerMessage)

type RoomMessageFunc func(RoomMessage)

type ServerMessage struct {
	client      *Client
	room        BaseRoom
	message     string
	messageType cm.MessageType
}

type RoomMessage struct {
	client      *Client
	message     string
	messageType cm.MessageType
}

// 客户端发出消息
type ClientMessage struct {
	Level   cm.MessageLevel `json:"level"`
	Type    cm.MessageType  `json:"type"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
}

// 服务端接收消息
type ReceiveMessage struct {
	Level   cm.MessageLevel `json:"level"`
	Type    cm.MessageType  `json:"type"`
	Message string          `json:"message"`
}
