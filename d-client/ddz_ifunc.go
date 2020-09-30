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
	case cm.StageGrabLandlord:
		GrabLandlord(dc, true)
	case cm.StagePlayPoker:
		dc.ShowMessage(cm.ClientLevel, "操作无效")
	default:
		dc.ShowMessage(cm.ClientLevel, "还未轮到操作")
	}
}

func (dc *DdzClient) NoCommand(val string) {
	switch dc.stage {
	case cm.StageGrabLandlord:
		GrabLandlord(dc, false)
	case cm.StagePlayPoker:
		dc.ShowMessage(cm.ClientLevel, "操作无效")
	default:
		dc.ShowMessage(cm.ClientLevel, "还未轮到操作")
	}
}

func (dc *DdzClient) PlayPoker(val string) {
	if len(val) == 0 {
		if dc.prevPoker == nil || dc.lastPlay == dc.userName {
			// 对局第一次出牌或出牌后无人压制
			dc.ShowMessage(cm.ClientLevel, "无法跳过本次出牌")
		} else {
			// 跳过出牌
			gm := GameMessage{"", cm.GamePlayPokerSkip}
			err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
			if err != nil {
				log.Fatal("PlayPoker error:", err)
			}
		}
		return
	}
	var tempPks []cm.Poker
	var pkIdx []int
	val = strings.ReplaceAll(val, "10", "0")
	// 记录已匹配poker
	cPkMap := make(map[int]byte)
	for i := 0; i < len(val); i++ {
		pkLevel := strings.ToUpper(val[i : i+1])
		if pkLevel == "0" {
			pkLevel = "10"
		}
		for j, pk := range dc.pokerSlice {
			if _, ok := cPkMap[j]; !ok && pk.Level == pkLevel {
				tempPks = append(tempPks, pk)
				pkIdx = append(pkIdx, j)
				cPkMap[j] = 0
				break
			}
		}
	}
	if len(pkIdx) == len(val) {
		pkt := cm.GetPokerType(tempPks)
		if pkt.PkType != cm.Invalid {
			if (dc.prevPoker != nil && cm.ComparePoker(tempPks, dc.prevPoker) == 0) ||
				dc.prevPoker == nil || dc.lastPlay == dc.userName {
				gm := GameMessage{cm.StructToJsonString(pkIdx), cm.GamePlayPoker}
				err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
				if err != nil {
					log.Fatal("PlayPoker error:", err)
				}
			} else {
				dc.ShowMessage(cm.ClientLevel, "当前出牌小于上家")
			}
		} else {
			dc.ShowMessage(cm.ClientLevel, "当前出牌类型无效")
		}
	} else {
		dc.ShowMessage(cm.ClientLevel, "当前出牌未通过检测")
	}
}

func (dc *DdzClient) ShowData(val string) {
	switch val {
	case "l":
		ShowLandlordData(dc)
	case "p":
		ShowSelfPokerData(dc)
	}
}

func GrabLandlord(dc *DdzClient, val bool) {
	gm := GameMessage{fmt.Sprint(val), cm.GameGrabLandlord}
	err := dc.conn.WriteJSON(SendMessage{cm.RoomLevel, cm.RoomGameMessage, cm.StructToJsonString(gm)})
	if err != nil {
		log.Fatal("GrabLandlord error:", err)
	}
}

func ShowLandlordData(dc *DdzClient) {
	dc.ShowMessage(cm.ClientLevel, fmt.Sprintf("当前房主为[%s]", dc.landlord))
}

func ShowSelfPokerData(dc *DdzClient) {
	dc.ShowSelfPoker()
}
