package main

import (
	cm "com.github/gc-common"
	"fmt"
	"log"
	"strconv"
)

// 中心
type Center struct {
	// 客户端id[客户端对象指针]
	clientMap map[ClientId]*Client
	// 房间id[具体房间指针]
	roomMap map[RoomId]BaseRoom
	// 消息channel
	messageChan chan ServerMessage
	funcMap     map[cm.MessageType]CenterMessageFunc
	// 最后房间id
	lastRoomId RoomId
}

func (c *Center) run() {
	for {
		select {
		// 消息监听
		case cMsg := <-c.messageChan:
			if cFunc, ok := c.funcMap[cMsg.messageType]; ok {
				go cFunc(cMsg)
			}
		}
	}
}

// 创建大厅
func newCenter() *Center {
	// 创建大厅实体
	center := &Center{
		clientMap:   make(map[ClientId]*Client),
		roomMap:     make(map[RoomId]BaseRoom),
		messageChan: make(chan ServerMessage),
		funcMap:     make(map[cm.MessageType]CenterMessageFunc),
		lastRoomId:  RoomId(0),
	}

	// 客户端连接
	center.funcMap[cm.ClientRegister] = center.clientRegister
	// 客户端断开连接
	center.funcMap[cm.ClientUnregister] = center.clientUnregister

	// 退出房间
	center.funcMap[cm.RoomQuit] = center.roomQuit
	// 解散房间
	center.funcMap[cm.RoomDisband] = center.roomDisband
	// 创建房间
	center.funcMap[cm.RoomCreate] = center.roomCreate
	// 加入房间
	center.funcMap[cm.RoomJoin] = center.roomJoin

	// 大厅运行监听客户端消息
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
			c.messageChan <- ServerMessage{messageType: cm.RoomQuit, client: client}
		}
	}
	client.center = nil
}

// 创建房间时调用
// 监听room消息
func (c *Center) roomStart(room BaseRoom) {
	defer func() {
		close(room.MessageChan())
		// 从大厅的roomMap移除room
		delete(c.roomMap, room.RoomId())
		log.Printf("房间[%d]完成解散", room.RoomId())
	}()
	for {
		select {
		// 监听房间消息
		case msg := <-room.MessageChan():
			if msg.messageType == cm.RoomClose {
				// 关闭房间
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
		msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomUnableCreate,
			false, ""}
		log.Printf("已在房间[%d]内,无法创建新房间", msg.client.currentRoom.RoomId())
		return
	}
	room := newDdzRoom(msg.client, c)
	// Room Message
	// 在这里添加是因为可以使用具体实现room的方法
	room.FuncMap()[cm.RoomReady] = room.Ready
	room.FuncMap()[cm.RoomCancelReady] = room.CancelReady
	room.FuncMap()[cm.RoomGameMessage] = room.GameMessage

	// 大厅roomMap添加新的room
	c.roomMap[room.RoomId()] = room
	// 房间clientMap添加client
	room.ClientMap()[msg.client.id] = msg.client
	room.ClientReadyMap()[msg.client.id] = false
	// client当前房间更新
	msg.client.currentRoom = room
	// 执行加入房间
	go room.Join(msg.client)
	// 房间运行
	go c.roomStart(room)
	// 广播房间创建 (其实只有创建者看到)
	room.BroadcastL(fmt.Sprint(room.RoomId()), cm.RoomCreate, cm.CenterLevel)
	log.Printf("房间创建[%d]", room.RoomId())
}

func (c *Center) roomJoin(msg ServerMessage) {
	// 转化为数字roomId
	roomId, err := strconv.ParseUint(msg.message, 10, 32)
	if err == nil {
		// roomId有效
		if room, ok := c.roomMap[RoomId(roomId)]; ok {
			// 获取当前消息的client
			client := msg.client
			// 检查client是否在room内
			if client.currentRoom != nil {
				msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomAlreadyIn, false, ""}
				log.Printf("已在房间[%d]内", client.currentRoom.RoomId())
				return
			}
			// 检查client是否在当前room内
			if _, ok := room.ClientMap()[client.id]; ok {
				msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomAlreadyIn, false, ""}
				log.Printf("已在房间[%d]内", room.RoomId())
				return
			}
			// 检查当前room下client是否已满
			if room.RoomSize() <= uint(len(room.ClientMap())) {
				msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomFull, false, ""}
				log.Printf("房间[%d]人员已满", room.RoomId())
				return
			}
			// client关联至房间
			room.ClientMap()[client.id] = client
			room.ClientReadyMap()[client.id] = false
			// client当前房间更新
			client.currentRoom = room
			// 执行加入房间
			go room.Join(client)
			// 广播有新client加入房间
			room.BroadcastM(client.userName, cm.RoomJoin)
			log.Printf("用户[%s]加入房间[%d]", client.userName, room.RoomId())
			var userNames []string
			for _, inClient := range room.ClientMap() {
				userNames = append(userNames, inClient.userName)
			}
			log.Printf("房间当前用户:%v", userNames)
			return
		}
	}
	msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomInvalid, false, ""}
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
		client.messageChan <- ClientMessage{cm.RoomLevel, cm.RoomQuit, true, ""}
		// room成员为0
		if len(client.currentRoom.ClientMap()) == 0 {
			// 发送解散房间
			c.messageChan <- ServerMessage{messageType: cm.RoomDisband, room: client.currentRoom}
		} else {
			// 向room内其他client广播退出信息
			client.currentRoom.BroadcastL(client.userName, cm.RoomSomeoneQuit, cm.RoomLevel)
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
		msg.client.messageChan <- ClientMessage{cm.CenterLevel, cm.RoomUnableExit, false, ""}
		log.Println("不在任何房间内无法退出")
	}
}

func (c *Center) roomDisband(msg ServerMessage) {
	if msg.room != nil && len(msg.room.ClientMap()) == 0 {
		log.Printf("房间[%d]即将解散\n", msg.room.RoomId())
		// 向房间发送close消息
		msg.room.MessageChan() <- RoomMessage{messageType: cm.RoomClose}
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
