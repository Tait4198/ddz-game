package main

import (
	cm "com.github/gc-common"
	"fmt"
	"log"
)

func (dc *DdzClient) CreateRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{cm.CenterLevel, cm.RoomCreate, val})
	if err != nil {
		log.Fatal("CreateRoom error:", err)
	}
}

func (dc *DdzClient) JoinRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{cm.CenterLevel, cm.RoomJoin, val})
	if err != nil {
		log.Fatal("JoinRoom error:", err)
	}
}

func (dc *DdzClient) QuitRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{cm.CenterLevel, cm.RoomQuit, val})
	if err != nil {
		log.Fatal("QuitRoom error:", err)
	}
}

func (dc *DdzClient) ReadyOrCancelRoom(val string) {
	var sendType cm.MessageType
	if dc.isReady {
		sendType = cm.RoomCancelReady
	} else {
		sendType = cm.RoomReady
	}
	err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, sendType, val})
	if err != nil {
		log.Fatal("ReadyRoom error:", err)
	}
}

func (dc *DdzClient) YesCommand(val string) {
	switch dc.stage {
	case StageGrabLandlord:
		GrabLandlord(dc, true)
	default:
		log.Println("还未轮到操作")
	}
}

func (dc *DdzClient) NoCommand(val string) {
	switch dc.stage {
	case StageGrabLandlord:
		GrabLandlord(dc, false)
	default:
		log.Println("还未轮到操作")
	}
}

func GrabLandlord(dc *DdzClient, val bool) {
	gm := GameMessage{fmt.Sprint(val), cm.GameGrabLandlord}
	err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
	if err != nil {
		log.Fatal("GrabLandlord error:", err)
	}
}
