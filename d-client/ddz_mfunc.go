package main

import "fmt"

func (dc *DdzClient) GameNewLandlord(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("当前地主[%s]", cm.Message))
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
	dc.isReady = true
}

func (dc *DdzClient) RoomCancelReady(cm ClientMessage) {
	dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]取消准备", cm.Message))
	dc.isReady = false
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
