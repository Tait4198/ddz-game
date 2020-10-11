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
	c.currentRoom.BroadcastL(c.Next.userName, cm.MessageType(cm.GameNextUserOps), cm.GameLevel)
	log.Printf("轮到[%s]操作", c.Next.userName)
	return c.Next
}

func (c *DdzClient) SelfClient() *DdzClient {
	c.currentRoom.BroadcastL(c.userName, cm.MessageType(cm.GameNextUserOps), cm.GameLevel)
	log.Printf("轮到[%s]操作", c.userName)
	return c
}

// 回合时间
const RoundTimeVal = 30 * 2
const SuitStr = "♠|♥|♣|♦"
const LevelStr = "3|4|5|6|7|8|9|10|J|Q|K|A|2"
const JokerS = "S"
const JokerX = "X"

type DdzRoom struct {
	BaseRoom
	iFuncMap     map[cm.DdzMessageType]DdzRoomMessageFunc
	iMessageChan chan DdzRoomMessage
	ddzStart     bool
	stageFuncMap map[cm.GameStage]DdzStageFunc
	closeCode    string

	landlord    *DdzClient
	roundClient *DdzClient
	roundTime   int
	waitUserOps bool
	ddzClients  []*DdzClient

	// 游戏阶段
	stage cm.GameStage

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
		stageFuncMap: make(map[cm.GameStage]DdzStageFunc),
	}
	ddzRoom.UpdateHomeowner(client)
	// ddz实现
	ddzRoom.iFuncMap[cm.GameGrabLandlord] = ddzRoom.GameGrabLandlord
	ddzRoom.iFuncMap[cm.GamePlayPoker] = ddzRoom.GamePlayPoker
	ddzRoom.iFuncMap[cm.GamePlayPokerSkip] = ddzRoom.GamePlayPokerSkip
	ddzRoom.iFuncMap[cm.GameExe] = ddzRoom.GameExe

	// 阶段方法
	ddzRoom.stageFuncMap[cm.StageGrabLandlord] = stageGrab
	ddzRoom.stageFuncMap[cm.StagePlayPoker] = stagePlay
	ddzRoom.stageFuncMap[cm.StageSettlement] = stageSettlement
	return ddzRoom
}

func (r *DdzRoom) GameExe(msg DdzRoomMessage) {
	if !r.ddzStart {
		r.GameStart(true)
	}
	r.roundTime -= 1
	if r.roundTime > 0 {
		if r.roundTime%20 == 0 && r.roundTime <= RoundTimeVal/3 {
			r.roundClient.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameCountdown),
				true, fmt.Sprint(r.roundTime / 2)}
			log.Printf("用户[%s]还剩%d秒操作时间", r.roundClient.userName, r.roundTime/2)
		}
	} else {
		r.roundClient.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameOpsTimeout),
			true, r.roundClient.userName}
		log.Printf("用户[%s]操作超时", r.roundClient.userName)
	}
	if r.waitUserOps && r.roundTime > 0 {
		return
	}
	// 如果roundTime大于0则代表着主动操作结束
	r.stageFuncMap[r.stage](!(r.roundTime > 0), r)
}

func stageGrab(auto bool, r *DdzRoom) {
	if auto {
		log.Printf("用户[%s]不抢地主(托管操作)", r.roundClient.userName)
		r.BroadcastL(r.roundClient.userName, cm.MessageType(cm.GameGrabHostingOps), cm.GameLevel)
		r.GameGrabLandlord(DdzRoomMessage{r.roundClient, "false", cm.GameGrabLandlord})
	}
	if r.grabIndex < 4 && r.stage == cm.StageGrabLandlord {
		if r.roundClient == nil {
			r.roundClient = r.landlord.SelfClient()
		} else {
			r.roundClient = r.roundClient.NextClient()
		}
		log.Printf("等待用户[%s]抢地主", r.roundClient.userName)
		r.roundClient.messageChan <- ClientMessage{cm.GameLevel,
			cm.MessageType(cm.GameWaitGrabLandlord), true, ""}
		r.waitUserOps = true
		r.roundTime = RoundTimeVal
	}
}

func stagePlay(auto bool, r *DdzRoom) {
	if auto {
		var pkIdx []int
		if len(r.prevPokers) == 0 || r.lastPlay == r.roundClient {
			pkIdx = PkAutoPlay([]cm.Poker{}, r.roundClient.PokerSlice)
		} else {
			pkIdx = PkAutoPlay(r.prevPokers, r.roundClient.PokerSlice)
		}
		if len(pkIdx) > 0 {
			log.Printf("用户[%s]自动出牌(托管操作)", r.roundClient.userName)
			r.GamePlayPoker(DdzRoomMessage{r.roundClient,
				cm.StructToJsonString(pkIdx), cm.GamePlayPoker})
		} else {
			log.Printf("用户[%s]跳过出牌(托管操作)", r.roundClient.userName)
			r.GamePlayPokerSkip(DdzRoomMessage{r.roundClient, "", cm.GamePlayPokerSkip})
		}
		r.BroadcastL(r.roundClient.userName, cm.MessageType(cm.GamePlayPokerHostingOps), cm.GameLevel)
	}
	if r.stage == cm.StagePlayPoker {
		if r.gameRound == 0 {
			r.roundClient = r.landlord.SelfClient()
		} else {
			r.roundClient = r.roundClient.NextClient()
		}
		log.Printf("[%s]准备出牌", r.roundClient.userName)
		r.waitUserOps = true
		r.roundTime = RoundTimeVal
		r.gameRound += 1
	}
}

func stageSettlement(auto bool, r *DdzRoom) {
	if r.stage == cm.StageSettlement {
		var winner string
		if r.lastPlay == r.landlord {
			winner = "地主"
		} else {
			winner = "农民"
		}
		r.BroadcastL(winner, cm.MessageType(cm.GameSettlement), cm.GameLevel)

		// 优胜者成为新地主
		r.UpdateLandlord(r.lastPlay)

		// 清空closeCode结束对局
		r.closeCode = ""
	}
}

func (r *DdzRoom) GameGrabLandlord(msg DdzRoomMessage) {
	b, e := strconv.ParseBool(msg.Message)
	if e == nil && b {
		r.lastGrab = r.roundClient
		log.Printf("用户[%s]抢地主", r.lastGrab.userName)
		r.BroadcastL(r.lastGrab.userName, cm.MessageType(cm.GameGrabLandlord), cm.GameLevel)
		if r.grabIndex == 0 {
			r.landlordGrab = true
		} else if r.grabIndex < 3 {
			r.otherGrab = true
		}
	} else {
		log.Printf("用户[%s]不抢地主", r.roundClient.userName)
		r.BroadcastL(r.roundClient.userName, cm.MessageType(cm.GameNGrabLandlord), cm.GameLevel)
	}
	r.grabIndex += 1
	if r.grabIndex == 4 || (r.grabIndex == 3 && (!r.landlordGrab || !r.otherGrab)) {
		if r.lastGrab == nil {
			// 重新开局'
			log.Println("未选出地主重新开局")
			r.BroadcastL("", cm.MessageType(cm.GameNoGrabLandlord), cm.GameLevel)
			r.GameStart(false)
			return
		}
		r.landlord = r.lastGrab
		r.roundClient = r.lastGrab
		r.BroadcastL(r.landlord.userName, cm.MessageType(cm.GameGrabLandlordEnd), cm.GameLevel)
		log.Printf("地主用户[%s]", r.landlord.userName)
		r.landlord.PokerSlice = append(r.landlord.PokerSlice, r.holePokers...)
		cm.SortPoker(r.landlord.PokerSlice, cm.SortByScore)
		holePokersJson, _ := json.Marshal(r.holePokers)
		// 公示底牌
		r.BroadcastL(string(holePokersJson), cm.MessageType(cm.GameShowHolePokers), cm.GameLevel)
		log.Println("底牌已公示")
		// 发送底牌给地主
		r.landlord.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameDealHolePokers),
			true, string(holePokersJson)}
		log.Println("底牌已发送")

		// 阶段结束
		r.stage = cm.StagePlayPoker
		r.roundTime = RoundTimeVal
		r.roundClient = nil
	}
	r.waitUserOps = false
}

func (r *DdzRoom) GamePlayPokerSkip(msg DdzRoomMessage) {
	// 广播消息
	r.BroadcastL("", cm.MessageType(cm.GamePlayPokerSkip), cm.GameLevel)
	r.waitUserOps = false
}

func (r *DdzRoom) GamePlayPoker(msg DdzRoomMessage) {
	var pkIdx []int
	err := json.Unmarshal([]byte(msg.Message), &pkIdx)
	playPokerInvalidMsg := ClientMessage{cm.GameLevel,
		cm.MessageType(cm.GamePlayPokerInvalid), true, ""}
	if err == nil {
		pkSlice := msg.client.PokerSlice
		if len(pkIdx) > len(pkSlice) {
			msg.client.messageChan <- playPokerInvalidMsg
			return
		}
		var playPks []cm.Poker
		for _, idx := range pkIdx {
			if idx < len(msg.client.PokerSlice) {
				playPks = append(playPks, pkSlice[idx])
			} else {
				msg.client.messageChan <- playPokerInvalidMsg
				return
			}
		}

		if r.prevPokers == nil || r.lastPlay == msg.client {
			// 对局第一次出牌或出牌后无人压制
			r.prevPokers = playPks
			r.lastPlay = msg.client
		} else if cm.ComparePoker(playPks, r.prevPokers) == 0 {
			// 出牌大于上家
			r.prevPokers = playPks
			r.lastPlay = msg.client
		} else {
			msg.client.messageChan <- playPokerInvalidMsg
			return
		}

		msg.client.PokerSlice = cm.PokerRemove(msg.client.PokerSlice, pkIdx)
		// 更新手牌
		msg.client.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GamePlayPokerUpdate),
			true, cm.StructToJsonString(msg.client.PokerSlice)}
		// 广播出牌
		upp := cm.UserPlayPoker{Pokers: playPks, Name: msg.client.userName}
		r.BroadcastL(cm.StructToJsonString(upp), cm.MessageType(cm.GamePlayPoker), cm.GameLevel)
		if len(msg.client.PokerSlice) == 0 {
			// 进入结算阶段
			r.stage = cm.StageSettlement
		} else if len(msg.client.PokerSlice) < 5 {
			// 剩余手牌提示
			upr := cm.UserPokerRemaining{Remaining: len(msg.client.PokerSlice), Name: msg.client.userName}
			r.BroadcastL(cm.StructToJsonString(upr), cm.MessageType(cm.GamePlayPokerRemaining), cm.GameLevel)
		}
		// 等待结束
		r.waitUserOps = false
	} else {
		panic(err)
	}
}

func (r *DdzRoom) GameStart(reRl bool) {
	log.Printf("房间[%d]对局开始", r.RoomId())
	// 重新关联代表着新game开始
	if reRl {
		r.BroadcastL("", cm.MessageType(cm.GameStart), cm.GameLevel)
	} else {
		r.BroadcastL("", cm.MessageType(cm.GameRestart), cm.GameLevel)
	}
	r.ddzStart = true
	r.waitUserOps = false
	r.roundClient = nil
	r.roundTime = RoundTimeVal

	r.stage = cm.StageGrabLandlord
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
			dc.Next = r.ddzClients[r.nextIndex(i, 1)]
			dc.Prev = r.ddzClients[r.nextIndex(i, -1)]
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

	for _, dc := range r.ddzClients {
		cm.SortPoker(dc.PokerSlice, cm.SortByScore)
		pokerJson, _ := json.Marshal(dc.PokerSlice)
		dc.messageChan <- ClientMessage{cm.GameLevel, cm.MessageType(cm.GameDealPoker), true, string(pokerJson)}
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
		r.iMessageChan <- DdzRoomMessage{MessageType: cm.GameStop, Message: r.closeCode}
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
		r.iMessageChan <- DdzRoomMessage{MessageType: cm.GameStop, Message: r.closeCode}
	}
}

func (r *DdzRoom) Run() {
	ticker := time.NewTicker(500 * time.Millisecond)
	r.iMessageChan = make(chan DdzRoomMessage)
	r.closeCode = fmt.Sprint(rand.Int31n(999999))
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
			if r.closeCode == "" {
				return
			}
			if msg.MessageType == cm.GameStop {
				if msg.Message == r.closeCode {
					return
				} else {
					log.Printf("CloseCode[%s]无效", msg.Message)
					continue
				}
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
			poker := cm.Poker{Level: level, Score: uint(i + 1), Suit: suit}
			pokerSlice = append(pokerSlice, poker)
		}
	}
	pokerSlice = append(pokerSlice, cm.Poker{Level: JokerS, Score: uint(len(levels) + 1), Suit: " "})
	pokerSlice = append(pokerSlice, cm.Poker{Level: JokerX, Score: uint(len(levels) + 2), Suit: " "})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pokerSlice), func(i, j int) {
		pokerSlice[i], pokerSlice[j] = pokerSlice[j], pokerSlice[i]
	})
	return pokerSlice
}

func (r *DdzRoom) UpdateLandlord(dc *DdzClient) {
	r.landlord = dc
	r.BroadcastL(dc.userName, cm.MessageType(cm.GameNewLandlord), cm.GameLevel)
}

func sliceRemove(slice []*DdzClient, s int) []*DdzClient {
	return append(slice[:s], slice[s+1:]...)
}
