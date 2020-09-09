package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// 客户端
type Client struct {
	id          ClientId
	userName    string
	currentRoom BaseRoom
	center      *Center
	conn        *websocket.Conn
	messageChan chan ClientMessage
}

func newClientId(usr string, pwd string) ClientId {
	return ClientId(fmt.Sprintf("u:%s_p:%s", usr, pwd))
}

func newClient(userName string, id ClientId, center *Center, conn *websocket.Conn) *Client {
	client := &Client{
		id:          id,
		userName:    userName,
		conn:        conn,
		center:      center,
		messageChan: make(chan ClientMessage),
	}
	client.center.messageChan <- ServerMessage{messageType: ClientRegister, client: client}
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.center.messageChan <- ServerMessage{messageType: ClientUnregister, client: c}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var rMessage ReceiveMessage
		err = json.Unmarshal(message, &rMessage)
		if err == nil {
			switch rMessage.Level {
			case CenterLevel:
				c.center.messageChan <- ServerMessage{message: rMessage.Message, messageType: rMessage.Type, client: c}
			case RoomLevel:
				if c.currentRoom != nil {
					c.currentRoom.MessageChan() <- RoomMessage{message: rMessage.Message, messageType: rMessage.Type, client: c}
				}
			}
		} else {
			log.Printf("Unknown message: %s", string(message))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.messageChan:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
