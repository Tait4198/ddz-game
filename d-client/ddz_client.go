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
	conn      *websocket.Conn
	host      string
	port      int
	userName  string
	password  string
	simplify  bool
	mFuncMap  map[cm.MessageType]MessageFunc
	dmFuncMap map[cm.DdzMessageType]MessageFunc
	iFuncMap  map[string]CommandFunc

	landlord   string
	roundUser  string
	isReady    bool
	stage      cm.GameStage
	pokerSlice []cm.Poker

	prevPoker []cm.Poker
	lastPlay  string
}

func (dc *DdzClient) DcReset() {
	dc.roundUser = ""
	dc.isReady = false
	dc.pokerSlice = nil
	dc.prevPoker = nil
	dc.lastPlay = ""
	dc.landlord = ""
	dc.stage = cm.StageWait
}

func (*DdzClient) ShowMessage(level cm.MessageLevel, message string) {
	switch level {
	case cm.CenterLevel:
		log.Printf("[lobby message]%s", message)
	case cm.RoomLevel:
		log.Printf("[room message]%s", message)
	case cm.GameLevel:
		log.Printf("[game message]%s", message)
	case cm.ClientLevel:
		log.Printf("[client message]%s", message)
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

func ShowPoker(title string, pks []cm.Poker, showIndex, simplify bool) {
	// ┌ └ ┐ ┘ ─ │ ├ ┤ ┬ ┴ ┼
	if pks != nil {
		var line0 = "┌"
		var line1 = "│"
		if simplify {
			line1 = ""
		}
		var line2 = "│"
		var line3 = "│"
		var line4 = "└"
		for i, pk := range pks {
			line0 += "──"
			line1 += fmtSprint(pk.Level)
			line2 += fmtSprint(pk.Suit)
			line3 += fmtSprint(fmt.Sprint(i))
			line4 += "──"
			if i < len(pks)-1 {
				line0 += "┬"
				if simplify {
					line1 += " "
				} else {
					line1 += "│"
				}
				line2 += "│"
				line3 += "│"
				line4 += "┴"
			}
		}
		line0 += "┐"
		if simplify {
			line1 += ""
		} else {
			line1 += "│"
		}
		line2 += "│"
		line3 += "│"
		line4 += "┘"

		var showPks string
		if showIndex {
			showPks = title + "\n" + line0 + "\n" + line1 + "\n" + line2 + "\n" + line3 + "\n" + line4
		} else if simplify {
			showPks = title + "\n" + line1
		} else {
			showPks = title + "\n" + line0 + "\n" + line1 + "\n" + line2 + "\n" + line4
		}
		log.Println(showPks)
	}
}

func (dc *DdzClient) ShowSelfPoker() {
	ShowPoker("手牌如下:", dc.pokerSlice, false, dc.simplify)
}

func NewDdzClient(usr, host string, port int, simplify bool) *DdzClient {
	dc := &DdzClient{
		userName:  usr,
		password:  "123456",
		host:      host,
		port:      port,
		simplify:  simplify,
		mFuncMap:  make(map[cm.MessageType]MessageFunc),
		dmFuncMap: make(map[cm.DdzMessageType]MessageFunc),
		iFuncMap:  make(map[string]CommandFunc),
	}
	// 房间创建
	dc.iFuncMap["c"] = dc.CreateRoom
	// 退出房间
	dc.iFuncMap["q"] = dc.QuitRoom
	// 加入房间
	dc.iFuncMap["j"] = dc.JoinRoom
	// 准备或取消准备
	dc.iFuncMap["r"] = dc.ReadyOrCancelRoom
	// 确定操作
	dc.iFuncMap["y"] = dc.YesCommand
	// 取消操作
	dc.iFuncMap["n"] = dc.NoCommand
	// 出牌或跳过出牌
	dc.iFuncMap["p"] = dc.PlayPoker
	// [s p] 显示手牌 [s l] 显示地主 [s s] 显示剩余手牌数量
	// [s r] 显示房间 [s cr] 显示当前房间
	dc.iFuncMap["s"] = dc.ShowData
	// 帮助
	dc.iFuncMap["h"] = dc.ShowHelp

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
	dc.mFuncMap[cm.RoomUnableCreate] = dc.RoomUnableCreate
	dc.mFuncMap[cm.RoomAlreadyIn] = dc.RoomAlreadyIn
	dc.mFuncMap[cm.RoomFull] = dc.RoomFull
	dc.mFuncMap[cm.RoomUnableExit] = dc.RoomUnableExit
	dc.mFuncMap[cm.RoomRun] = dc.RoomRun
	dc.mFuncMap[cm.RoomClose] = dc.RoomClose
	dc.mFuncMap[cm.GetAllRoomInfo] = dc.GetAllRoomInfo
	dc.mFuncMap[cm.GetCurRoomInfo] = dc.GetCurRoomInfo

	dc.dmFuncMap[cm.GameNewLandlord] = dc.GameNewLandlord
	dc.dmFuncMap[cm.GameStart] = dc.GameStart
	dc.dmFuncMap[cm.GameRestart] = dc.GameRestart
	dc.dmFuncMap[cm.GameCountdown] = dc.GameCountdown
	dc.dmFuncMap[cm.GameNextUserOps] = dc.GameNextUserOps
	dc.dmFuncMap[cm.GameWaitGrabLandlord] = dc.GameWaitGrabLandlord
	dc.dmFuncMap[cm.GameGrabHostingOps] = dc.GameGrabHostingOps
	dc.dmFuncMap[cm.GameGrabLandlord] = dc.GameGrabLandlord
	dc.dmFuncMap[cm.GameNGrabLandlord] = dc.GameNGrabLandlord
	dc.dmFuncMap[cm.GameGrabLandlordEnd] = dc.GameGrabLandlordEnd
	dc.dmFuncMap[cm.GameDealPoker] = dc.GameDealPoker
	dc.dmFuncMap[cm.GameDealHolePokers] = dc.GameDealHolePokers
	dc.dmFuncMap[cm.GameShowHolePokers] = dc.GameShowHolePokers
	dc.dmFuncMap[cm.GamePlayPoker] = dc.GamePlayPoker
	dc.dmFuncMap[cm.GamePlayPokerUpdate] = dc.GamePlayPokerUpdate
	dc.dmFuncMap[cm.GamePlayPokerSkip] = dc.GamePlayPokerSkip
	dc.dmFuncMap[cm.GamePlayPokerRemaining] = dc.GamePlayPokerRemaining
	dc.dmFuncMap[cm.GameSettlement] = dc.GameSettlement
	dc.dmFuncMap[cm.GamePlayPokerHostingOps] = dc.GamePlayPokerHostingOps
	dc.dmFuncMap[cm.GameOpsTimeout] = dc.GameOpsTimeout
	dc.dmFuncMap[cm.GameStop] = dc.GameStop
	dc.dmFuncMap[cm.GamePokerRemaining] = dc.GamePokerRemaining

	return dc
}

func (dc *DdzClient) Run() {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", dc.host, dc.port), Path: "/ws"}
	q := u.Query()
	q.Set("usr", dc.userName)
	q.Set("pwd", dc.password)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	} else {
		dc.ShowMessage(cm.ClientLevel, fmt.Sprintf("已连接至服务器 %s:%d", dc.host, dc.port))
		dc.ShowMessage(cm.CenterLevel, fmt.Sprintf("当前用户 %s", dc.userName))
		dc.ShowMessage(cm.CenterLevel, "输入 h 查看帮助信息")
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
			var cMsg ClientMessage
			err = json.Unmarshal(message, &cMsg)
			if err != nil {
				log.Fatal("unmarshal:", err)
				return
			}
			if mFunc, ok := dc.mFuncMap[cMsg.Type]; ok {
				go mFunc(cMsg)
			} else if dmFunc, ok := dc.dmFuncMap[cm.DdzMessageType(cMsg.Type)]; ok {
				go dmFunc(cMsg)
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
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\n", "")
		if match, _ := regexp.MatchString("^(\\w [\\w\\d]+)|(\\w+ ?)$", text); !match {
			dc.ShowMessage(cm.ClientLevel, "输入无效")
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
			dc.ShowMessage(cm.ClientLevel, "命令无效")
		}
	}
}
