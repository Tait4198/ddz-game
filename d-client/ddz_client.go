package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type DdzClient struct {
	conn     *websocket.Conn
	userName string
	password string
	mFuncMap map[MessageType]MessageFunc
	iFuncMap map[string]InstructionFunc

	landlord string
}

func (*DdzClient) ShowMessage(level, message string) {
	switch level {
	case CenterLevel:
		log.Printf("[大厅消息]%s", message)
	case RoomLevel:
		log.Printf("[房间消息]%s", message)
	case GameLevel:
		log.Printf("[游戏消息]%s", message)
	}
}

func NewDdzClient(usr, pwd string) *DdzClient {
	dc := &DdzClient{
		userName: usr,
		password: pwd,
		mFuncMap: make(map[MessageType]MessageFunc),
		iFuncMap: make(map[string]InstructionFunc),
	}
	// 房间创建
	dc.iFuncMap["c"] = dc.CreateRoom
	dc.iFuncMap["q"] = dc.QuitRoom
	dc.iFuncMap["j"] = dc.JoinRoom

	// 消息监听
	dc.mFuncMap[RoomCreate] = dc.RoomCreate
	dc.mFuncMap[RoomJoin] = dc.RoomJoin
	dc.mFuncMap[GameNewLandlord] = dc.GameNewLandlord
	return dc
}

func (dc *DdzClient) CreateRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomCreate, val})
	if err != nil {
		log.Fatal("CreateRoom error:", err)
	}
}

func (dc *DdzClient) JoinRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomJoin, val})
	if err != nil {
		log.Fatal("JoinRoom error:", err)
	}
}

func (dc *DdzClient) QuitRoom(val string) {
	err := dc.conn.WriteJSON(SendMessage{CenterLevel, RoomQuit, val})
	if err != nil {
		log.Fatal("QuitRoom error:", err)
	}
}

func (dc *DdzClient) GameNewLandlord(cm ClientMessage) {
	if cm.Status {
		dc.ShowMessage(cm.Level, fmt.Sprintf("当前地主[%s]", cm.Message))
	}
}

func (dc *DdzClient) RoomCreate(cm ClientMessage) {
	if cm.Status {
		dc.ShowMessage(cm.Level, fmt.Sprintf("房间[%s]已创建", cm.Message))
	}
}

func (dc *DdzClient) RoomJoin(cm ClientMessage) {
	if cm.Status {
		dc.ShowMessage(cm.Level, fmt.Sprintf("用户[%s]加入房间", cm.Message))
		// 可做额外操作
	}
}

func (dc *DdzClient) Run() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	q := u.Query()
	q.Set("usr", dc.userName)
	q.Set("pwd", dc.password)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	dc.conn = c
	defer dc.conn.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var cm ClientMessage
			err = json.Unmarshal(message, &cm)
			if err != nil {
				log.Fatal("unmarshal:", err)
				return
			}
			if mFunc, ok := dc.mFuncMap[cm.Type]; ok {
				go mFunc(cm)
			}
			log.Printf("recv: %s", string(message))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("write:", err)
			return
		}
		text = strings.ReplaceAll(text, "\n", "")
		if match, _ := regexp.MatchString("^(\\w [\\w\\d]+)|(\\w+)$", text); !match {
			log.Println("无效输入")
			continue
		}
		arr := strings.Split(text, " ")

		if iFunc, ok := dc.iFuncMap[arr[0]]; ok {
			if len(arr) > 1 {
				iFunc(arr[1])
			} else {
				iFunc("")
			}
		}
	}
}
