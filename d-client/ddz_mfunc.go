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
}

func (dc *DdzClient) RoomMissUser(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "缺少用户或存在未准备用户")
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
		dc.isReady = true
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
	dc.ShowMessage(cm.Level, fmt.Sprintf("轮到[%s]操作", cm.Message))
	dc.roundUser = cm.Message
}

func (dc *DdzClient) GameWaitGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "是否抢地主 (y/n)")
	dc.stage = StageGrabLandlord
}

func (dc *DdzClient) GameGrabHostingOps(cm ClientMessage) {
	dc.ShowMessage(cm.Level, "托管操作不抢地主")
}

func (dc *DdzClient) GameGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]抢地主", cm.Message))
}

func (dc *DdzClient) GameNGrabLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]不抢地主", cm.Message))
}

func (dc *DdzClient) GameGrabLandlordEnd(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("地主用户[%s]", cm.Message))
	dc.landlord = cm.Message
}

func (dc *DdzClient) GameDealPoker(cm ClientMessage) {
	var pks []gcm.Poker
	if err := json.Unmarshal([]byte(cm.Message), &pks); err != nil {
		panic(err)
	}
	dc.pokerSlice = pks

	gcm.SortPoker(dc.pokerSlice, func(p, q *gcm.Poker) bool {
		return p.Score < q.Score
	})
	dc.ShowPoker()
}
