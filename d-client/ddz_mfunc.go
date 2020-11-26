package main

import (
	"com.github/gc-client/lang"
	gcm "com.github/gc-common"
	"encoding/json"
	"fmt"
)

func (dc *DdzClient) GameNewLandlord(cm ClientMessage) {
	if cm.Message != "" {
		dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.CurrentLandlord), cm.Message))
	} else {
		dc.ShowMessage(cm.Level, dc.lang.Get(lang.NotElectLandlord))
	}
}

func (dc *DdzClient) RoomCreate(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.EstablishRoom), cm.Message))
}

func (dc *DdzClient) RoomJoin(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.IntoRoom), cm.Message))
}

func (dc *DdzClient) RoomInvalid(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.NotFindRoom))
}

func (dc *DdzClient) RoomQuit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.QuitRoom))
	// 重置状态
	dc.DcReset()
}

func (dc *DdzClient) RoomMissUser(cm ClientMessage) {
	if cm.Message != "" {
		dc.ShowMessage(cm.Level, cm.Message)
	} else {
		dc.ShowMessage(cm.Level, dc.lang.Get(lang.PreparationFailed))
	}
}

func (dc *DdzClient) RoomReady(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.UserPrepare), cm.Message))
	if cm.Message == dc.userName {
		dc.isReady = true
	}
}

func (dc *DdzClient) RoomCancelReady(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.CancelPrepared), cm.Message))
	if cm.Message == dc.userName {
		dc.isReady = false
	}
}

func (dc *DdzClient) RoomSomeoneQuit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.UserQuitRoom), cm.Message))
}

func (dc *DdzClient) RoomNewHomeowner(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.Homeowner), cm.Message))
}

func (dc *DdzClient) RoomFull(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.IsFull))
}

func (dc *DdzClient) RoomAlreadyIn(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.RepeatIntoRoom))
}

func (dc *DdzClient) RoomUnableCreate(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.CanNotEstablishRoom))
}

func (dc *DdzClient) RoomUnableExit(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.QuitRoomFailed))
}

func (dc *DdzClient) RoomRun(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.GameStart))
}

func (dc *DdzClient) RoomClose(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.RoomClosed))
}

func (dc *DdzClient) GameStart(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "游戏开始")
}

func (dc *DdzClient) GameRestart(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "游戏重新开始")
}

func (dc *DdzClient) GameCountdown(cm ClientMessage) {
	if cm.Message == "0" {
		dc.ShowMessage(cm.Level, dc.lang.Get(lang.OperateTimeOut))
	} else {
		dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.TimeRemain), cm.Message))
	}
}

func (dc *DdzClient) GameNextUserOps(cm ClientMessage) {
	dc.roundUser = cm.Message
	if cm.Message == dc.userName {
		dc.ShowMessage(gcm.ClientLevel, fmt.Sprintf(dc.lang.Get(lang.Operate)))
		if dc.stage == gcm.StagePlayPoker {
			dc.ShowSelfPoker()
		}

	} else {
		dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.TurnToUser), cm.Message))
	}
}

func (dc *DdzClient) GameWaitGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.WhetherSeizeLandlord))
	dc.stage = gcm.StageGrabLandlord
}

func (dc *DdzClient) GameGrabHostingOps(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.TrustSeizeLandlord), cm.Message))
}

func (dc *DdzClient) GameGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.SeizeLandlord), cm.Message))
}

func (dc *DdzClient) GameNGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.NotSeizeLandlord), cm.Message))
}

func (dc *DdzClient) GameGrabLandlordEnd(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.LandlordUser), cm.Message))

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
	ShowPoker(dc.lang.Get(lang.HoleCard), pks, false, dc.simplify)
}

func (dc *DdzClient) GameDealHolePokers(cm ClientMessage) {
	pks := convertPokers(cm.Message)
	dc.pokerSlice = append(dc.pokerSlice, pks...)
	gcm.SortPoker(dc.pokerSlice, gcm.SortByScore)
}

func (dc *DdzClient) GamePlayPoker(cm ClientMessage) {
	var upp gcm.UserPlayInfo
	if err := json.Unmarshal([]byte(cm.Message), &upp); err != nil {
		panic(err)
	}
	gcm.SortPoker(upp.Pokers, gcm.SortByScore)
	dc.prevPoker = upp.Pokers
	dc.lastPlay = upp.Name
	identity := dc.lang.Get(lang.Farmer)
	if upp.Name == dc.landlord {
		identity = dc.lang.Get(lang.Landlord)
	}
	ShowPoker(fmt.Sprintf(dc.lang.Get(lang.TurnToDiscard), identity, upp.Name, upp.Remaining),
		upp.Pokers, false, dc.simplify)
}

func (dc *DdzClient) GamePlayPokerUpdate(cm ClientMessage) {
	dc.pokerSlice = convertPokers(cm.Message)
}

func (dc *DdzClient) GamePlayPokerSkip(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.SkipToDiscard), dc.roundUser))
}

func (dc *DdzClient) GamePlayPokerRemaining(cm ClientMessage) {
	var upr gcm.UserPlayInfo
	if err := json.Unmarshal([]byte(cm.Message), &upr); err != nil {
		panic(err)
	}
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.AttentionRemain), upr.Name, upr.Remaining))
}

func (dc *DdzClient) GameSettlement(cm ClientMessage) {
	dc.stage = gcm.StageSettlement
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.UserWin), cm.Message))
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.WinnerLandlord), dc.roundUser))
}

func (dc *DdzClient) GamePlayPokerHostingOps(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.TrustToDiscard), cm.Message))
}

func (dc *DdzClient) GameOpsTimeout(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf(dc.lang.Get(lang.UserOperateTimeOut), cm.Message))
}

func (dc *DdzClient) GameStop(cm ClientMessage) {
	dc.ShowMessage(cm.Level, dc.lang.Get(lang.GameOver))
	dc.DcReset()
	dc.ShowMessage(gcm.ClientLevel, dc.lang.Get(lang.NextGame))
}

func (dc *DdzClient) GetCurRoomInfo(message ClientMessage) {
	var rr gcm.ResultRoom
	if err := json.Unmarshal([]byte(message.Message), &rr); err != nil {
		panic(err)
	}
	roomInfoStr := dc.lang.Get(lang.CurrentRoomInfo)
	roomInfoStr += rr.String()
	dc.ShowMessage(message.Level, roomInfoStr)
}

func (dc *DdzClient) GetAllRoomInfo(message ClientMessage) {
	var rrs []gcm.ResultRoom
	if err := json.Unmarshal([]byte(message.Message), &rrs); err != nil {
		panic(err)
	}
	roomInfoStr := dc.lang.Get(lang.AllRoomInfo)
	for _, rr := range rrs {
		roomInfoStr += rr.String()
	}
	dc.ShowMessage(message.Level, roomInfoStr)
}

func (dc *DdzClient) GameChat(message ClientMessage) {
	dc.ShowMessage(message.Level, message.Message)
}

func (dc *DdzClient) GamePokerRemaining(message ClientMessage) {
	var us []gcm.UserPlayInfo
	if err := json.Unmarshal([]byte(message.Message), &us); err != nil {
		panic(err)
	}
	uprStr := dc.lang.Get(lang.UserCardRemain)
	for _, upr := range us {
		uprStr += fmt.Sprintf(dc.lang.Get(lang.CardNum), upr.Name, upr.Remaining)
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
