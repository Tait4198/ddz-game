package main

import (
	cm "com.github/gc-common"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type DdzStageFunc func(bool, *DdzRoom)

type DdzRoomMessageFunc func(DdzRoomMessage)

type DdzRoomMessage struct {
	client      *DdzClient
	Message     string            `json:"message"`
	MessageType cm.DdzMessageType `json:"type"`
}

type DdzClient struct {
	*Client
	Prev       *DdzClient
	Next       *DdzClient
	PokerSlice []cm.Poker
}

func (c *DdzClient) NextClient() *DdzClient {
	c.currentRoom.BroadcastL(c.Next.userName, cm.GameNextUserOps, cm.GameLevel)
	log.Printf("轮到[%s]操作", c.Next.userName)
	return c.Next
}

// 回合时间
const RoundTimeVal = 60
const SuitStr = "♠|♥|♣|♦"
const LevelStr = "3|4|5|6|7|8|9|10|J|Q|K|A|2"
const JokerS = "S"
const JokerX = "X"

type DdzRoom struct {
	BaseRoom
	iFuncMap     map[cm.DdzMessageType]DdzRoomMessageFunc
	iMessageChan chan DdzRoomMessage
	ddzStart     bool
	stageFuncMap map[uint]DdzStageFunc

	landlord    *DdzClient
	roundClient *DdzClient
	roundTime   uint
	waitUserOps bool
	ddzClients  []*DdzClient

	// 阶段计数
	stageIndex uint

	holePokers []cm.Poker

	// 最后抢地主client
	lastGrab *DdzClient
	// 抢地主计数
	grabIndex uint
	// 默认地主是否抢地主
	landlordGrab bool
	// 其他client是否抢地主
	otherGrab bool

	// 前一位出牌
	prevPokers []cm.Poker
	// 最后出牌者
	lastPlay *DdzClient

	// 回合统计
	gameRound uint
}

func newDdzClient(client *Client) *DdzClient {
	dc := &DdzClient{
		Client: client,
	}
	return dc
}

func newDdzRoom(client *Client, center *Center) BaseRoom {
	ddzRoom := &DdzRoom{
		BaseRoom:     newRoom(center),
		iFuncMap:     make(map[cm.DdzMessageType]DdzRoomMessageFunc),
		stageFuncMap: make(map[uint]DdzStageFunc),
	}
	ddzRoom.UpdateHomeowner(client)
	// ddz实现
	ddzRoom.iFuncMap[cm.GameGrabLandlord] = ddzRoom.GameGrabLandlord
	ddzRoom.iFuncMap[cm.GamePlayPoker] = ddzRoom.GamePlayPoker
	ddzRoom.iFuncMap[cm.GameExe] = ddzRoom.GameExe

	// 阶段方法
	ddzRoom.stageFuncMap[0] = stageGrab
	ddzRoom.stageFuncMap[1] = stagePlay
	return ddzRoom
}

func (r *DdzRoom) GameExe(msg DdzRoomMessage) {
	if !r.ddzStart {
		r.GameStart(true)
	}

	if r.roundTime > 0 {
		r.roundTime -= 1
		if r.roundTime%10 == 0 && r.roundTime < RoundTimeVal/3 {
			r.roundClient.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameCountdown),
				true, fmt.Sprint(r.roundTime / 2)}
			log.Printf("用户[%s]还剩%d秒操作时间", r.roundClient.userName, r.roundTime/2)
		}
	} else {
		r.roundClient.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameOpsTimeout),
			true, ""}
		log.Printf("用户[%s]操作超时", r.roundClient.userName)
	}

	if r.waitUserOps && r.roundTime > 0 {
		return
	}
	// 如果roundTime大于0则代表着主动操作结束
	r.stageFuncMap[r.stageIndex](!(r.roundTime > 0), r)
}

func stageGrab(auto bool, r *DdzRoom) {
	if auto {
		log.Printf("用户[%s]不抢地主(托管操作)", r.roundClient.userName)
		r.BroadcastL(r.roundClient.userName, cm.GameGrabHostingOps, cm.GameLevel)
		r.GameGrabLandlord(DdzRoomMessage{r.roundClient, "false", cm.GameGrabLandlord})
	}
	if r.grabIndex < 4 {
		if r.roundClient == nil {
			r.roundClient = r.landlord
		} else {
			r.roundClient = r.roundClient.NextClient()
		}
		log.Printf("等待用户[%s]抢地主", r.roundClient.userName)
		r.roundClient.messageChan <- ClientMessage{cm.GameLevel, cm.GameWaitGrabLandlord, true, ""}
		r.waitUserOps = true
		r.roundTime = RoundTimeVal
	}
}

func stagePlay(auto bool, r *DdzRoom) {
	if r.gameRound == 0 {
		r.roundClient = r.landlord
		r.BroadcastL(r.roundClient.userName, cm.GameNextUserOps, cm.GameLevel)
		log.Printf("地主[%s]准备出牌", r.roundClient.userName)
	} else {
		r.roundClient = r.roundClient.NextClient()
		log.Printf("[%s]准备出牌", r.roundClient.userName)
	}
	r.roundClient.messageChan <- ClientMessage{cm.GameLevel, 0, true, "准备出牌"}
	r.waitUserOps = true
	r.roundTime = RoundTimeVal
	r.gameRound += 1
}

func (r *DdzRoom) GameGrabLandlord(msg DdzRoomMessage) {
	b, e := strconv.ParseBool(msg.Message)
	if e == nil && b {
		r.lastGrab = r.roundClient
		log.Printf("用户[%s]抢地主", r.lastGrab.userName)
		r.BroadcastL(r.lastGrab.userName, cm.GameGrabLandlord, cm.GameLevel)
		if r.grabIndex == 0 {
			r.landlordGrab = true
		} else if r.grabIndex < 3 {
			r.otherGrab = true
		}
	} else {
		log.Printf("用户[%s]不抢地主", r.roundClient.userName)
		r.BroadcastL(r.roundClient.userName, cm.GameNGrabLandlord, cm.GameLevel)
	}
	r.grabIndex += 1
	r.waitUserOps = false
	if r.grabIndex == 4 || (r.grabIndex == 3 && (!r.landlordGrab || !r.otherGrab)) {
		if r.lastGrab == nil {
			// 重新开局'
			log.Println("未选出地主重新开局")
			r.BroadcastL("", cm.GameNoGrabLandlord, cm.GameLevel)
			r.GameStart(false)
			return
		}
		// 阶段结束
		r.landlord = r.lastGrab
		r.roundClient = r.lastGrab
		r.stageIndex += 1
		r.BroadcastL(r.landlord.userName, cm.GameGrabLandlordEnd, cm.GameLevel)
		log.Printf("地主用户[%s]", r.landlord.userName)
		r.landlord.PokerSlice = append(r.landlord.PokerSlice, r.holePokers...)
		cm.SortPoker(r.landlord.PokerSlice, cm.SortByScore)
		holePokersJson, _ := json.Marshal(r.holePokers)
		// 公示底牌
		r.BroadcastL(string(holePokersJson), cm.GameShowHolePokers, cm.GameLevel)
		log.Println("底牌已公示")
		// 发送底牌给地主
		r.landlord.messageChan <- ClientMessage{cm.GameLevel, cm.GameDealHolePokers,
			true, string(holePokersJson)}
		log.Println("底牌已发送")
	}
}

func (r *DdzRoom) GamePlayPoker(msg DdzRoomMessage) {
	var pkIdx []int
	err := json.Unmarshal([]byte(msg.Message), &pkIdx)
	if err == nil {
		pkSlice := msg.client.PokerSlice
		if len(pkIdx) > len(pkSlice) {
			// error
		}
		var tempPks []cm.Poker
		for _, idx := range pkIdx {
			if idx < len(msg.client.PokerSlice) {
				tempPks = append(tempPks, pkSlice[idx])
			} else {
				// error
			}
		}
		msg.client.PokerSlice = cm.PokerRemove(msg.client.PokerSlice, pkIdx)

		if r.prevPokers == nil {
			r.prevPokers = tempPks
		}

		log.Printf("%s出牌", msg.client.userName)
		for _, pk := range tempPks {
			log.Println(pk)
		}

		msg.client.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GamePlayPokerUpdate),
			true, cm.StructToJsonString(msg.client.PokerSlice)}

		r.BroadcastL(cm.StructToJsonString(tempPks), cm.GamePlayPoker, cm.GameLevel)
	} else {
		panic(err)
	}
}

func (r *DdzRoom) GameStart(reRl bool) {
	log.Printf("房间[%d]对局开始", r.RoomId())
	// 重新关联代表着新game开始
	if reRl {
		r.BroadcastL("", cm.GameStart, cm.GameLevel)
	} else {
		r.BroadcastL("", cm.GameRestart, cm.GameLevel)
	}
	r.ddzStart = true
	r.waitUserOps = false
	r.roundClient = nil
	r.roundTime = RoundTimeVal

	r.stageIndex = 0
	r.holePokers = nil

	r.lastGrab = nil
	r.grabIndex = 0
	r.landlordGrab = false
	r.otherGrab = false

	r.lastPlay = nil
	r.prevPokers = nil

	r.gameRound = 0

	// 对局开始调用
	// 建立关联
	if reRl {
		for i, dc := range r.ddzClients {
			dc.Next = r.ddzClients[r.nextIndex(i, -1)]
			dc.Prev = r.ddzClients[r.nextIndex(i, 1)]
		}
	}

	for _, dc := range r.ddzClients {
		dc.PokerSlice = nil
	}

	log.Printf("房间[%d]开始发牌", r.RoomId())
	pokers := r.RandomPokerSlice()
	m := cm.RandIntMap(0, len(pokers)-1, 3)
	pc := r.landlord
	for i, p := range pokers {
		if _, ok := m[i]; ok {
			r.holePokers = append(r.holePokers, p)
		} else {
			pc.PokerSlice = append(pc.PokerSlice, p)
			pc = pc.Next
		}
	}

	log.Println("底牌")
	for _, p := range r.holePokers {
		log.Println(p.String())
	}

	for _, dc := range r.ddzClients {
		cm.SortPoker(dc.PokerSlice, cm.SortByScore)
		pokerJson, _ := json.Marshal(dc.PokerSlice)
		dc.messageChan <- ClientMessage{cm.GameLevel, cm.GameDealPoker, true, string(pokerJson)}
	}

}

func (r *DdzRoom) GameEnd() {
	log.Printf("房间[%d]对局结束", r.RoomId())
	r.ddzStart = false

	close(r.iMessageChan)

	r.ResetReady()
	r.BroadcastL("", cm.MessageType(cm.GameStop), cm.GameLevel)
	// 取消关联
	for _, dc := range r.ddzClients {
		dc.Prev = nil
		dc.Next = nil
	}
}

func (r *DdzRoom) Quit(c *Client) {
	rmIdx := -1
	for i, dc := range r.ddzClients {
		if dc.Client == c {
			rmIdx = i
			break
		}
	}
	r.ddzClients = sliceRemove(r.ddzClients, rmIdx)
	if r.ddzStart {
		r.iMessageChan <- DdzRoomMessage{MessageType: cm.GameStop}
	}
	if r.landlord.Client == c {
		for _, nextClient := range r.ddzClients {
			r.UpdateLandlord(nextClient)
			break
		}
	}
}

func (r *DdzRoom) Join(c *Client) {
	dc := newDdzClient(c)
	r.ddzClients = append(r.ddzClients, dc)
	if r.landlord == nil {
		r.UpdateLandlord(dc)
	}
}

func (*DdzRoom) RoomSize() uint {
	return 3
}

func (r *DdzRoom) GameMessage(msg RoomMessage) {
	if r.IsRun() && r.roundClient.Client == msg.client {
		drm := DdzRoomMessage{}
		err := json.Unmarshal([]byte(msg.message), &drm)
		if err == nil {
			for _, dc := range r.ddzClients {
				if dc.Client == msg.client {
					drm.client = dc
					break
				}
			}
			r.iMessageChan <- drm
		}
	} else if r.roundClient.Client != msg.client {
		msg.client.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameInvalidOps),
			false, ""}
		log.Printf("用户[%s]无效操作", r.roundClient.userName)
	}
}

func (r *DdzRoom) Stop() {
	if r.IsRun() {
		// 增加随机close识别
		r.iMessageChan <- DdzRoomMessage{MessageType: cm.GameStop}
	}
}

func (r *DdzRoom) Run() {
	ticker := time.NewTicker(500 * time.Millisecond)
	r.iMessageChan = make(chan DdzRoomMessage)
	defer func() {
		r.GameEnd()
		ticker.Stop()
	}()
	go func() {
		for {
			select {
			case <-ticker.C:
				r.iMessageChan <- DdzRoomMessage{MessageType: cm.GameExe}
			}
		}
	}()
	for {
		select {
		case msg := <-r.iMessageChan:
			if msg.MessageType == cm.GameStop {
				return
			}
			if cFunc, ok := r.iFuncMap[msg.MessageType]; ok {
				cFunc(msg)
			} else {
				log.Printf("无效消息消息[%d]", msg.MessageType)
			}
		}
	}
}

func (r *DdzRoom) nextIndex(cur, add int) uint {
	if cur+add < 0 {
		return r.RoomSize() - 1
	} else {
		return uint(cur+add) % r.RoomSize()
	}
}

func (r *DdzRoom) RandomPokerSlice() []cm.Poker {
	var pokerSlice []cm.Poker
	suits := strings.Split(SuitStr, "|")
	levels := strings.Split(LevelStr, "|")
	for i, level := range levels {
		for _, suit := range suits {
			poker := cm.Poker{Level: level, Score: uint(i), Suit: suit}
			pokerSlice = append(pokerSlice, poker)
		}
	}
	pokerSlice = append(pokerSlice, cm.Poker{Level: JokerS, Score: uint(len(levels)), Suit: " "})
	pokerSlice = append(pokerSlice, cm.Poker{Level: JokerX, Score: uint(len(levels) + 1), Suit: " "})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pokerSlice), func(i, j int) {
		pokerSlice[i], pokerSlice[j] = pokerSlice[j], pokerSlice[i]
	})
	return pokerSlice
}

func (r *DdzRoom) UpdateLandlord(dc *DdzClient) {
	r.landlord = dc
	r.BroadcastL(dc.userName, cm.GameNewLandlord, cm.GameLevel)
}

func sliceRemove(slice []*DdzClient, s int) []*DdzClient {
	return append(slice[:s], slice[s+1:]...)
}
