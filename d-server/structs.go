package main

// 客户端id
type ClientId string

// 房间id
type RoomId uint

//消息发往 center/room
type MessageLevel string

// 用于指定调用方法
type MessageType uint

type CenterMessageFunc func(ServerMessage)

type RoomMessageFunc func(RoomMessage)

type ServerMessage struct {
	client      *Client
	room        BaseRoom
	message     string
	messageType MessageType
}

type RoomMessage struct {
	client      *Client
	message     string
	messageType MessageType
}

// 客户端发出消息
type ClientMessage struct {
	Level   MessageLevel `json:"level"`
	Type    MessageType  `json:"type"`
	Status  bool         `json:"status"`
	Message string       `json:"message"`
}

// 服务端接收消息
type ReceiveMessage struct {
	Level   MessageLevel `json:"level"`
	Type    MessageType  `json:"type"`
	Message string       `json:"message"`
}
