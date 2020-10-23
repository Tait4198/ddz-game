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
	if dc.stage == cm.StageWait {
		err := dc.conn.WriteJSON(SendMessage{cm.CenterLevel, cm.RoomQuit, val})
		if err != nil {
			log.Fatal("QuitRoom error:", err)
		}
	} else {
		dc.ShowMessage(cm.ClientLevel, "操作无效")
	}
}

func (dc *DdzClient) ReadyOrCancelRoom(val string) {
	if dc.stage == cm.StageWait {
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
	} else {
		dc.ShowMessage(cm.ClientLevel, "操作无效")
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
	if dc.stage == cm.StagePlayPoker && dc.roundUser == dc.userName {
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
	} else {
		dc.ShowMessage(cm.ClientLevel, "操作无效")
	}
}

func (dc *DdzClient) ShowHelp(val string) {
	helpStr := "\n***游戏帮助***\n"
	var tipHelpSl []string
	tipHelpSl = append(tipHelpSl, "流程:创建/加入房间 -> 准备操作 -> 开始游戏")
	tipHelpSl = append(tipHelpSl, "特殊牌型 S = 小王 / X = 大王 / 0或10 = 10")
	for i, s := range tipHelpSl {
		helpStr += fmt.Sprintf("%d. %s\n", i+1, s)
	}
	helpStr += "***命令帮助***\n"
	var cmdHelpSl []string
	cmdHelpSl = append(cmdHelpSl, "创建房间 -> c")
	cmdHelpSl = append(cmdHelpSl, "退出房间 -> q")
	cmdHelpSl = append(cmdHelpSl, "加入指定房间 -> j 房间数字id")
	cmdHelpSl = append(cmdHelpSl, "房间内准备或取消 -> r")
	cmdHelpSl = append(cmdHelpSl, "游戏内出牌 -> p 牌型对应数字或字母,如 p 334455 / p 90jqk (忽略大小写)")
	cmdHelpSl = append(cmdHelpSl, "游戏内跳过出牌 -> p (仅输入p)")
	cmdHelpSl = append(cmdHelpSl, "显示当前手牌信息 -> s p")
	cmdHelpSl = append(cmdHelpSl, "显示当前地主信息 -> s l")
	cmdHelpSl = append(cmdHelpSl, "显示所有房间信息 -> s r")
	for i, s := range cmdHelpSl {
		helpStr += fmt.Sprintf("%d. %s\n", i+1, s)
	}
	dc.ShowMessage(cm.ClientLevel, helpStr)
}

func (dc *DdzClient) ShowData(val string) {
	switch val {
	case "l":
		ShowLandlordData(dc)
	case "p":
		ShowSelfPokerData(dc)
	case "r":
		ShowAllRoomData(dc)
	default:
		dc.ShowMessage(cm.ClientLevel, "无效命令")
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
	if dc.landlord != "" {
		dc.ShowMessage(cm.ClientLevel, fmt.Sprintf("当前地主为[%s]", dc.landlord))
	} else {
		dc.ShowMessage(cm.ClientLevel, "该命令仅在对局选出地主后有效")
	}
}

func ShowSelfPokerData(dc *DdzClient) {
	if dc.stage == cm.StagePlayPoker {
		dc.ShowSelfPoker()
	} else {
		dc.ShowMessage(cm.ClientLevel, "当前不在游戏阶段")
	}
}

func ShowAllRoomData(dc *DdzClient) {
	err := dc.conn.WriteJSON(SendMessage{cm.CenterLevel, cm.GetRoomInfo, ""})
	if err != nil {
		log.Fatal("GetRoomInfo error:", err)
	}
}
