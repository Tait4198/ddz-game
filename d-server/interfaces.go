package main

import cm "com.github/gc-common"

type BaseRoom interface {
	RoomId() RoomId
	MessageChan() chan RoomMessage
	ClientMap() map[ClientId]*Client
	RoomSize() uint
	Broadcast(cm.MessageType)
	BroadcastM(string, cm.MessageType)
	BroadcastL(string, cm.MessageType, cm.MessageLevel)
	RemoveClient(ClientId)
	Homeowner() *Client
	UpdateHomeowner(*Client)
	FuncMap() map[cm.MessageType]RoomMessageFunc
	Start(RoomMessage)
	Stop()
	Quit(*Client)
	Join(*Client)
	GameMessage(RoomMessage)
	Ready(RoomMessage)
	CancelReady(RoomMessage)
	ResetReady()
	IsRun() bool
	Run()
}
