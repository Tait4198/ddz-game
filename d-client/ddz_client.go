package main

import (
	"bufio"
	"encoding/json"
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
	isReady  bool
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
	// 退出房间
	dc.iFuncMap["q"] = dc.QuitRoom
	// 加入房间
	dc.iFuncMap["j"] = dc.JoinRoom
	// 准备或取消准备
	dc.iFuncMap["r"] = dc.ReadyOrCancelRoom

	// 消息监听
	dc.mFuncMap[RoomCreate] = dc.RoomCreate
	dc.mFuncMap[RoomJoin] = dc.RoomJoin
	dc.mFuncMap[RoomInvalid] = dc.RoomInvalid
	dc.mFuncMap[RoomQuit] = dc.RoomQuit
	dc.mFuncMap[RoomReady] = dc.RoomReady
	dc.mFuncMap[RoomCancelReady] = dc.RoomCancelReady
	dc.mFuncMap[RoomSomeoneQuit] = dc.RoomSomeoneQuit
	dc.mFuncMap[RoomMissUser] = dc.RoomMissUser
	dc.mFuncMap[RoomNewHomeowner] = dc.RoomNewHomeowner
	dc.mFuncMap[GameNewLandlord] = dc.GameNewLandlord
	dc.mFuncMap[RoomUnableCreate] = dc.RoomUnableCreate
	dc.mFuncMap[RoomAlreadyIn] = dc.RoomAlreadyIn
	dc.mFuncMap[RoomFull] = dc.RoomFull
	dc.mFuncMap[RoomUnableExit] = dc.RoomUnableExit
	dc.mFuncMap[RoomRun] = dc.RoomRun
	dc.mFuncMap[RoomClose] = dc.RoomClose

	dc.mFuncMap[GameStart] = dc.GameStart
	dc.mFuncMap[GameRestart] = dc.GameRestart
	dc.mFuncMap[GameCountdown] = dc.GameCountdown
	dc.mFuncMap[GameNextUserOps] = dc.GameNextUserOps
	dc.mFuncMap[GameWaitGrabLandlord] = dc.GameWaitGrabLandlord
	return dc
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

	// 服务端消息监听
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
			} else {
				// 如果未找到处理方法直接输出
				log.Printf("recv: %s", string(message))
			}
		}
	}()

	// 输入监听
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
