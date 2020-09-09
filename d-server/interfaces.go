package main

type BaseRoom interface {
	RoomId() RoomId
	MessageChan() chan RoomMessage
	ClientMap() map[ClientId]*Client
	RoomSize() uint
	Broadcast(MessageType)
	BroadcastM(string, MessageType)
	BroadcastL(string, MessageType, MessageLevel)
	RemoveClient(ClientId)
	Homeowner() *Client
	UpdateHomeowner(*Client)
	FuncMap() map[MessageType]RoomMessageFunc
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
