package main

import (
	gcm "com.github/gc-common"
	"encoding/json"
	"fmt"
)

func (dc *DdzClient) GameNewLandlord(cm ClientMessage) {
	if cm.Message != "" {
		dc.ShowMessage(cm.Level, fmt.Sprintf("当前地主[%s]", cm.Message))
	} else {
		dc.ShowMessage(cm.Level, "未选出地主")
	}
}

func (dc *DdzClient) RoomCreate(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("房间[%s]已创建", cm.Message))
}

func (dc *DdzClient) RoomJoin(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]加入房间", cm.Message))
}

func (dc *DdzClient) RoomInvalid(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "未找到对应房间")
}

func (dc *DdzClient) RoomQuit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "退出房间")
	// 重置状态
	dc.DcReset()
}

func (dc *DdzClient) RoomMissUser(cm ClientMessage) {
	if cm.Message != "" {
		dc.ShowMessage(cm.Level, cm.Message)
	} else {
		dc.ShowMessage(cm.Level, "无法准备(未知原因)")
	}
}

func (dc *DdzClient) RoomReady(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]已准备", cm.Message))
	if cm.Message == dc.userName {
		dc.isReady = true
	}
}

func (dc *DdzClient) RoomCancelReady(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]取消准备", cm.Message))
	if cm.Message == dc.userName {
		dc.isReady = false
	}
}

func (dc *DdzClient) RoomSomeoneQuit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]退出房间", cm.Message))
}

func (dc *DdzClient) RoomNewHomeowner(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]成为房主", cm.Message))
}

func (dc *DdzClient) RoomFull(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "房间已满")
}

func (dc *DdzClient) RoomAlreadyIn(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "已在房间内无法加入房间")
}

func (dc *DdzClient) RoomUnableCreate(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "已在房间内无法创建新房间")
}

func (dc *DdzClient) RoomUnableExit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "不在任何房间内无法退出")
}

func (dc *DdzClient) RoomRun(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "对局开始")
}

func (dc *DdzClient) RoomClose(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "房间关闭")
}

func (dc *DdzClient) GameStart(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "游戏开始")
}

func (dc *DdzClient) GameRestart(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "游戏重新开始")
}

func (dc *DdzClient) GameCountdown(cm ClientMessage) {
	if cm.Message == "0" {
		dc.ShowMessage(cm.Level, "操作时间用尽")
	} else {
		dc.ShowMessage(cm.Level, fmt.Sprintf("操作时间还剩%s秒", cm.Message))
	}
}

func (dc *DdzClient) GameNextUserOps(cm ClientMessage) {
	dc.roundUser = cm.Message
	if cm.Message == dc.userName {
		dc.ShowMessage(gcm.ClientLevel, fmt.Sprintf("***请操作***"))
		if dc.stage == gcm.StagePlayPoker {
			dc.ShowSelfPoker()
		}

	} else {
		dc.ShowMessage(cm.Level, fmt.Sprintf("轮到[%s]操作", cm.Message))
	}
}

func (dc *DdzClient) GameWaitGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "是否抢地主 (y/n)")
	dc.stage = gcm.StageGrabLandlord
}

func (dc *DdzClient) GameGrabHostingOps(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("[%s]托管操作不抢地主", cm.Message))
}

func (dc *DdzClient) GameGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]抢地主", cm.Message))
}

func (dc *DdzClient) GameNGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]不抢地主", cm.Message))
}

func (dc *DdzClient) GameGrabLandlordEnd(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("***地主用户[%s]***", cm.Message))

	dc.landlord = cm.Message
	dc.stage = gcm.StagePlayPoker
}

func (dc *DdzClient) GameDealPoker(cm ClientMessage) {
	pks := convertPokers(cm.Message)
	dc.pokerSlice = pks
	dc.ShowSelfPoker()
}

func (dc *DdzClient) GameShowHolePokers(cm ClientMessage) {
	pks := convertPokers(cm.Message)
	gcm.SortPoker(pks, gcm.SortByScore)
	ShowPoker("底牌:", pks, false, dc.simplify)
}

func (dc *DdzClient) GameDealHolePokers(cm ClientMessage) {
	pks := convertPokers(cm.Message)
	dc.pokerSlice = append(dc.pokerSlice, pks...)
	gcm.SortPoker(dc.pokerSlice, gcm.SortByScore)
}

func (dc *DdzClient) GamePlayPoker(cm ClientMessage) {
	var upp gcm.UserPlayPoker
	if err := json.Unmarshal([]byte(cm.Message), &upp); err != nil {
		panic(err)
	}
	gcm.SortPoker(upp.Pokers, gcm.SortByScore)
	dc.prevPoker = upp.Pokers
	dc.lastPlay = upp.Name
	ShowPoker(fmt.Sprintf("***[%s]出牌***", upp.Name), upp.Pokers, false, dc.simplify)
}

func (dc *DdzClient) GamePlayPokerUpdate(cm ClientMessage) {
	dc.pokerSlice = convertPokers(cm.Message)
}

func (dc *DdzClient) GamePlayPokerSkip(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("[%s]跳过出牌", dc.roundUser))
}

func (dc *DdzClient) GamePlayPokerRemaining(cm ClientMessage) {
	var upr gcm.UserPokerRemaining
	if err := json.Unmarshal([]byte(cm.Message), &upr); err != nil {
		panic(err)
	}
	dc.ShowMessage(cm.Level, fmt.Sprintf("***注意[%s]还剩%d张手牌***", upr.Name, upr.Remaining))
}

func (dc *DdzClient) GameSettlement(cm ClientMessage) {
	dc.stage = gcm.StageSettlement
	dc.ShowMessage(cm.Level, fmt.Sprintf("***[%s]获得对局胜利***", cm.Message))
	dc.ShowMessage(cm.Level, fmt.Sprintf("***[%s]为优胜者成为地主***", dc.roundUser))
}

func (dc *DdzClient) GamePlayPokerHostingOps(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("[%s]托管操作出牌", cm.Message))
}

func (dc *DdzClient) GameOpsTimeout(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("[%s]操作超时将由系统托管操作", cm.Message))
}

func (dc *DdzClient) GameStop(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "对局结束")
	dc.DcReset()
	dc.ShowMessage(gcm.ClientLevel, "可进行下一场对局")
}

func (dc *DdzClient) GetCurRoomInfo(message ClientMessage) {
	var rr gcm.ResultRoom
	if err := json.Unmarshal([]byte(message.Message), &rr); err != nil {
		panic(err)
	}
	roomInfoStr := "当前房间信息\n"
	roomInfoStr += rr.String()
	dc.ShowMessage(message.Level, roomInfoStr)
}

func (dc *DdzClient) GetAllRoomInfo(message ClientMessage) {
	var rrs []gcm.ResultRoom
	if err := json.Unmarshal([]byte(message.Message), &rrs); err != nil {
		panic(err)
	}
	roomInfoStr := "所有房间信息\n"
	for _, rr := range rrs {
		roomInfoStr += rr.String()
	}
	dc.ShowMessage(message.Level, roomInfoStr)
}

func (dc *DdzClient) GamePokerRemaining(message ClientMessage) {
	var us []gcm.UserPokerRemaining
	if err := json.Unmarshal([]byte(message.Message), &us); err != nil {
		panic(err)
	}
	uprStr := "用户剩余卡牌\n"
	for _, upr := range us {
		uprStr += fmt.Sprintf("%s[%d]张 ", upr.Name, upr.Remaining)
	}
	dc.ShowMessage(message.Level, uprStr)
}

func convertPokers(pokerJson string) []gcm.Poker {
	var pks []gcm.Poker
	if err := json.Unmarshal([]byte(pokerJson), &pks); err != nil {
		panic(err)
	}
	return pks
}
