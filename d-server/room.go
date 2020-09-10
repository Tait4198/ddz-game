package main

import (
	cm "com.github/gc-common"
	"log"
)

// 房间
type Room struct {
	id          RoomId
	homeowner   *Client
	clientMap   map[ClientId]*Client
	messageChan chan RoomMessage
	center      *Center
	funcMap     map[cm.MessageType]RoomMessageFunc

	roomReady      bool
	clientReadyMap map[ClientId]bool
}

func (r *Room) RoomId() RoomId {
	return r.id
}

func (r *Room) ClientMap() map[ClientId]*Client {
	return r.clientMap
}

func (r *Room) MessageChan() chan RoomMessage {
	return r.messageChan
}

func (r *Room) RoomSize() uint {
	return 1
}

func (r *Room) Broadcast(msgType cm.MessageType) {
	r.BroadcastL("", msgType, cm.RoomLevel)
}

func (r *Room) BroadcastM(msg string, msgType cm.MessageType) {
	r.BroadcastL(msg, msgType, cm.RoomLevel)
}

func (r *Room) BroadcastL(msg string, msgType cm.MessageType, level cm.MessageLevel) {
	for _, client := range r.clientMap {
		client.messageChan <- ClientMessage{level, msgType, true, msg}
	}
}

func (r *Room) RemoveClient(id ClientId) {
	if _, ok := r.clientMap[id]; ok {
		delete(r.clientMap, id)
	}
}

func (r *Room) Homeowner() *Client {
	return r.homeowner
}

func (r *Room) UpdateHomeowner(client *Client) {
	r.homeowner = client
	r.BroadcastL(client.userName, cm.RoomNewHomeowner, cm.RoomLevel)
}

func (r *Room) FuncMap() map[cm.MessageType]RoomMessageFunc {
	return r.funcMap
}

func (r *Room) Start(RoomMessage) {
	log.Printf("房间[%d]任务开始", r.RoomId())
}

func (r *Room) Stop() {
	// 房间全局停止消息接收时调用
	// 可重写用于发送下级停止消息
	log.Printf("房间[%d]任务停止", r.RoomId())
}

func (r *Room) Quit(c *Client) {
	//log.Printf("用户[%s]退出房间[%d]", c.userName, r.RoomId())
}

func (r *Room) Join(c *Client) {
	//log.Printf("用户[%s]加入房间[%d]", c.userName, r.RoomId())
}

func (r *Room) GameMessage(msg RoomMessage) {
	// 接收 game level 消息转换发送
}

func (r *Room) Run() {
	log.Println("Game Start")
}

func (r *Room) Ready(msg RoomMessage) {
	c := msg.client
	if r.homeowner.id == c.id {
		allReady := true
		r.clientReadyMap[c.id] = false
		for ci, ready := range r.clientReadyMap {
			if !ready && ci != r.homeowner.id {
				allReady = false
				break
			}
		}
		if allReady && uint(len(r.clientReadyMap)) == c.currentRoom.RoomSize() {
			r.clientReadyMap[c.id] = true
			r.roomReady = true
			// 广播对局开始
			r.Broadcast(cm.RoomRun)
			log.Printf("房间[%d]对局运行\n", r.id)
			// 对局开始
			// 使用客户端所在到具体实现room进行开局
			c.currentRoom.Run()
		} else {
			log.Println("还存在未准备用户或缺少用户")
			c.messageChan <- ClientMessage{cm.RoomLevel, cm.RoomMissUser, false, ""}
		}
	} else {
		r.clientReadyMap[c.id] = true
		r.BroadcastM(c.userName, cm.RoomReady)
		log.Printf("房间[%d]用户[%s]准备\n", r.id, c.userName)
	}

}

// client取消准备
func (r *Room) CancelReady(msg RoomMessage) {
	c := msg.client
	if !r.roomReady {
		r.clientReadyMap[c.id] = false
		r.BroadcastM(c.userName, cm.RoomCancelReady)
		log.Printf("房间[%d]用户[%s]取消准备\n", r.id, c.userName)
	}
}

// room对局是否正在进行
func (r *Room) IsRun() bool {
	return r.roomReady
}

// client准备
func (r *Room) ResetReady() {
	r.roomReady = false
	for id := range r.clientReadyMap {
		r.clientReadyMap[id] = false
	}
}

func newRoom(center *Center) BaseRoom {
	room := &Room{
		id:             center.nextRoomId(),
		clientMap:      make(map[ClientId]*Client),
		messageChan:    make(chan RoomMessage),
		center:         center,
		funcMap:        make(map[cm.MessageType]RoomMessageFunc),
		roomReady:      false,
		clientReadyMap: make(map[ClientId]bool),
	}
	return room
}
