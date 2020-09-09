package main

import (
	"fmt"
	"log"
	"strconv"
)

// 中心
type Center struct {
	clientMap map[ClientId]*Client
	roomMap   map[RoomId]BaseRoom
	// 消息channel
	messageChan chan ServerMessage
	funcMap     map[MessageType]CenterMessageFunc
	// 最后房间id
	lastRoomId RoomId
}

func (c *Center) run() {
	for {
		select {
		case cMsg := <-c.messageChan:
			if cFunc, ok := c.funcMap[cMsg.messageType]; ok {
				go cFunc(cMsg)
			}
		}
	}
}

func newCenter() *Center {
	center := &Center{
		clientMap:   make(map[ClientId]*Client),
		roomMap:     make(map[RoomId]BaseRoom),
		messageChan: make(chan ServerMessage),
		funcMap:     make(map[MessageType]CenterMessageFunc),
		lastRoomId:  RoomId(0),
	}

	center.funcMap[ClientRegister] = center.clientRegister
	center.funcMap[ClientUnregister] = center.clientUnregister

	center.funcMap[RoomQuit] = center.roomQuit
	center.funcMap[RoomDisband] = center.roomDisband
	center.funcMap[RoomCreate] = center.roomCreate
	center.funcMap[RoomJoin] = center.roomJoin
	go center.run()
	return center
}

func (c *Center) clientRegister(msg ServerMessage) {
	log.Printf("用户[%s]加入\n", msg.client.userName)
	c.clientMap[msg.client.id] = msg.client
}

func (c *Center) clientUnregister(msg ServerMessage) {
	log.Printf("[%s]用户退出\n", msg.client.userName)
	client := msg.client
	clientId := client.id
	if _, ok := c.clientMap[clientId]; ok {
		delete(c.clientMap, clientId)
		if client.currentRoom != nil {
			c.messageChan <- ServerMessage{messageType: RoomQuit, client: client}
		}
	}
	client.center = nil
}

// 创建房间时调用
// 监听room消息
func (c *Center) roomStart(room BaseRoom) {
	defer func() {
		close(room.MessageChan())
		delete(c.roomMap, room.RoomId())
		log.Printf("房间[%d]完成解散", room.RoomId())
	}()
	for {
		select {
		case msg := <-room.MessageChan():
			if msg.messageType == RoomClose {
				go room.Stop()
				return
			}
			if cFunc, ok := room.FuncMap()[msg.messageType]; ok {
				go cFunc(msg)
			}
		}
	}
}

func (c *Center) roomCreate(msg ServerMessage) {
	if msg.client.currentRoom != nil {
		msg.client.messageChan <- ClientMessage{CenterLevel, RoomUnableCreate,
			false, ""}
		log.Printf("已在房间[%d]内,无法创建新房间", msg.client.currentRoom.RoomId())
		return
	}
	room := newDdzRoom(msg.client, c)
	// Room Message
	room.FuncMap()[RoomReady] = room.Ready
	room.FuncMap()[RoomCancelReady] = room.CancelReady
	room.FuncMap()[RoomGameMessage] = room.GameMessage

	c.roomMap[room.RoomId()] = room
	room.ClientMap()[msg.client.id] = msg.client
	msg.client.currentRoom = room
	go room.Join(msg.client)
	go c.roomStart(room)
	room.BroadcastL(fmt.Sprint(room.RoomId()), RoomCreate, CenterLevel)
	log.Printf("房间创建[%d]", room.RoomId())
}

func (c *Center) roomJoin(msg ServerMessage) {
	roomId, err := strconv.ParseUint(msg.message, 10, 32)
	if err == nil {
		if room, ok := c.roomMap[RoomId(roomId)]; ok {
			client := msg.client
			if _, ok := room.ClientMap()[client.id]; ok {
				msg.client.messageChan <- ClientMessage{CenterLevel, RoomAlreadyIn, false, ""}
				log.Printf("已在房间[%d]内", room.RoomId())
				return
			}
			if room.RoomSize() <= uint(len(room.ClientMap())) {
				msg.client.messageChan <- ClientMessage{CenterLevel, RoomFull, false, ""}
				log.Printf("房间[%d]人员已满", room.RoomId())
				return
			}
			room.ClientMap()[client.id] = client
			client.currentRoom = room
			go room.Join(client)
			room.BroadcastM(client.userName, RoomJoin)
			log.Printf("用户[%s]加入房间[%d]", client.userName, room.RoomId())
			var userNames []string
			for _, inClient := range room.ClientMap() {
				userNames = append(userNames, inClient.userName)
			}
			log.Printf("房间当前用户:%v", userNames)
			return
		}
	}
	msg.client.messageChan <- ClientMessage{CenterLevel, RoomInvalid, false, ""}
	log.Printf("无效房间[%d]", roomId)
}

func (c *Center) roomQuit(msg ServerMessage) {
	if msg.client.currentRoom != nil {
		log.Printf("用户[%s]退出房间[%d]\n", msg.client.userName, msg.client.currentRoom.RoomId())
		client := msg.client
		clientId := client.id
		// 移除room内的client
		client.currentRoom.RemoveClient(clientId)
		// 执行room退出操作
		go client.currentRoom.Quit(client)
		// 发送退出room信息
		client.messageChan <- ClientMessage{RoomLevel, RoomQuit, true, ""}
		// room成员为0
		if len(client.currentRoom.ClientMap()) == 0 {
			// 发送解散房间
			c.messageChan <- ServerMessage{messageType: RoomDisband, room: client.currentRoom}
		} else {
			// 向room内其他client广播退出信息
			client.currentRoom.BroadcastL(client.userName, RoomSomeoneQuit, RoomLevel)
			// 退出client为房主
			if client == client.currentRoom.Homeowner() {
				// 从room内其他client选择新房主
				for _, nextClient := range client.currentRoom.ClientMap() {
					// 更新房主
					client.currentRoom.UpdateHomeowner(nextClient)
					log.Printf("用户[%s]升为房主", nextClient.userName)
					break
				}
			}
		}
		client.currentRoom = nil
	} else {
		msg.client.messageChan <- ClientMessage{CenterLevel, RoomUnableExit, false, ""}
		log.Println("不在任何房间内无法退出")
	}
}

func (c *Center) roomDisband(msg ServerMessage) {
	if msg.room != nil && len(msg.room.ClientMap()) == 0 {
		log.Printf("房间[%d]即将解散\n", msg.room.RoomId())
		// 向房间发送close消息
		msg.room.MessageChan() <- RoomMessage{messageType: RoomClose}
	}
}

func (c *Center) nextRoomId() RoomId {
	c.lastRoomId += 1
	return c.lastRoomId
}

func (c *Center) checkClientId(id ClientId) bool {
	if _, ok := c.clientMap[id]; ok {
		return true
	}
	return false
}
