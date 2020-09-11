package main

import (
	"bufio"
	cm "com.github/gc-common"
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
	mFuncMap map[cm.MessageType]MessageFunc
	iFuncMap map[string]CommandFunc

	landlord   string
	roundUser  string
	isReady    bool
	stage      GameStage
	pokerSlice []cm.Poker
}

func (*DdzClient) ShowMessage(level cm.MessageLevel, message string) {
	switch level {
	case cm.CenterLevel:
		log.Printf("[大厅消息]%s", message)
	case cm.RoomLevel:
		log.Printf("[房间消息]%s", message)
	case cm.GameLevel:
		log.Printf("[游戏消息]%s", message)
	}
}

func fmtSprint(str string) string {
	if len(str) == 1 {
		return fmt.Sprintf(" %s", str)
	} else if ok, err := regexp.Match("[♠♥♣♦]", []byte(str)); err == nil && ok {
		return fmt.Sprintf(" %s", str)
	} else {
		return fmt.Sprintf("%s", str)
	}
}

func (dc *DdzClient) ShowPoker() {
	// ┌ └ ┐ ┘ ─ │ ├ ┤ ┬ ┴ ┼
	if dc.pokerSlice != nil {
		var line0 = "┌"
		var line1 = "│"
		var line2 = "│"
		var line3 = "│"
		var line4 = "└"
		for i, pk := range dc.pokerSlice {
			line0 += "──"
			line1 += fmtSprint(pk.Level)
			line2 += fmtSprint(pk.Suit)
			line3 += fmtSprint(fmt.Sprint(i))
			line4 += "──"
			if i < len(dc.pokerSlice)-1 {
				line0 += "┬"
				line1 += "│"
				line2 += "│"
				line3 += "│"
				line4 += "┴"
			}
		}
		line0 += "┐"
		line1 += "│"
		line2 += "│"
		line3 += "│"
		line4 += "┘"

		log.Println(line0)
		log.Println(line1)
		log.Println(line2)
		log.Println(line3)
		log.Println(line4)
	}
}

func NewDdzClient(usr, pwd string) *DdzClient {
	dc := &DdzClient{
		userName: usr,
		password: pwd,
		mFuncMap: make(map[cm.MessageType]MessageFunc),
		iFuncMap: make(map[string]CommandFunc),
	}
	// 房间创建
	dc.iFuncMap["c"] = dc.CreateRoom
	// 退出房间
	dc.iFuncMap["q"] = dc.QuitRoom
	// 加入房间
	dc.iFuncMap["j"] = dc.JoinRoom
	// 准备或取消准备
	dc.iFuncMap["r"] = dc.ReadyOrCancelRoom
	// 准备或取消准备
	dc.iFuncMap["y"] = dc.YesCommand
	// 准备或取消准备
	dc.iFuncMap["n"] = dc.NoCommand

	// 消息监听
	dc.mFuncMap[cm.RoomCreate] = dc.RoomCreate
	dc.mFuncMap[cm.RoomJoin] = dc.RoomJoin
	dc.mFuncMap[cm.RoomInvalid] = dc.RoomInvalid
	dc.mFuncMap[cm.RoomQuit] = dc.RoomQuit
	dc.mFuncMap[cm.RoomReady] = dc.RoomReady
	dc.mFuncMap[cm.RoomCancelReady] = dc.RoomCancelReady
	dc.mFuncMap[cm.RoomSomeoneQuit] = dc.RoomSomeoneQuit
	dc.mFuncMap[cm.RoomMissUser] = dc.RoomMissUser
	dc.mFuncMap[cm.RoomNewHomeowner] = dc.RoomNewHomeowner
	dc.mFuncMap[cm.GameNewLandlord] = dc.GameNewLandlord
	dc.mFuncMap[cm.RoomUnableCreate] = dc.RoomUnableCreate
	dc.mFuncMap[cm.RoomAlreadyIn] = dc.RoomAlreadyIn
	dc.mFuncMap[cm.RoomFull] = dc.RoomFull
	dc.mFuncMap[cm.RoomUnableExit] = dc.RoomUnableExit
	dc.mFuncMap[cm.RoomRun] = dc.RoomRun
	dc.mFuncMap[cm.RoomClose] = dc.RoomClose

	dc.mFuncMap[cm.GameStart] = dc.GameStart
	dc.mFuncMap[cm.GameRestart] = dc.GameRestart
	dc.mFuncMap[cm.GameCountdown] = dc.GameCountdown
	dc.mFuncMap[cm.GameNextUserOps] = dc.GameNextUserOps
	dc.mFuncMap[cm.GameWaitGrabLandlord] = dc.GameWaitGrabLandlord
	dc.mFuncMap[cm.GameGrabHostingOps] = dc.GameGrabHostingOps
	dc.mFuncMap[cm.GameGrabLandlord] = dc.GameGrabLandlord
	dc.mFuncMap[cm.GameNGrabLandlord] = dc.GameNGrabLandlord
	dc.mFuncMap[cm.GameGrabLandlordEnd] = dc.GameGrabLandlordEnd
	dc.mFuncMap[cm.GameDealPoker] = dc.GameDealPoker
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
		} else {
			log.Println("无效命令")
		}
	}
}
