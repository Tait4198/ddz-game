package main

import "log"

func (dc *DdzClient) CreateRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomCreate, val})
	if err != nil {
		log.Fatal("CreateRoom error:", err)
	}
}

func (dc *DdzClient) JoinRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomJoin, val})
	if err != nil {
		log.Fatal("JoinRoom error:", err)
	}
}

func (dc *DdzClient) QuitRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomQuit, val})
	if err != nil {
		log.Fatal("QuitRoom error:", err)
	}
}

func (dc *DdzClient) ReadyOrCancelRoom(val string) {
	var sendType MessageType
	if dc.isReady {
		sendType = RoomCancelReady
	} else {
		sendType = RoomReady
	}
	err := dc.conn.WriteJSON(SendMessage{RoomLevel, sendType, val})
	if err != nil {
		log.Fatal("ReadyRoom error:", err)
	}
}
