package main

import (
	cm "com.github/gc-common"
	"fmt"
	"log"
	"strings"
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

func (dc *DdzClient) PlayPoker(val string) {
	if len(val) == 0 {
		log.Println("Pass")
	}
	var tempPks []cm.Poker
	var pkIdx []int
	playIdx := 0
	isValid := false
	val = strings.ReplaceAll(val, "10", "0")
	for i, pk := range dc.pokerSlice {
		pkLevel := strings.ToUpper(val[playIdx : playIdx+1])
		if pkLevel == "0" {
			pkLevel = "10"
		}
		if pk.Level == pkLevel {
			tempPks = append(tempPks, pk)
			pkIdx = append(pkIdx, i)
			playIdx++
			if playIdx >= len(val) {
				isValid = true
				break
			}
		}
	}
	if isValid {
		pkt := cm.GetPokerType(tempPks)
		if pkt.PkType != cm.Invalid {
			if dc.prevPoker != nil && cm.ComparePoker(tempPks, dc.prevPoker) == 0 || dc.prevPoker == nil {
				gm := GameMessage{cm.StructToJsonString(pkIdx), cm.GamePlayPoker}
				err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
				if err != nil {
					log.Fatal("PlayPoker error:", err)
				}
			} else {
				log.Println("")
			}
		}
	} else {
		log.Println("无效出牌")
	}
}

func GrabLandlord(dc *DdzClient, val bool) {
	gm := GameMessage{fmt.Sprint(val), cm.GameGrabLandlord}
	err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
	if err != nil {
		log.Fatal("GrabLandlord error:", err)
	}
}
