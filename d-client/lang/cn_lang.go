package lang

type CnLang struct {
	langMap map[LanguageKey]string
}

func NewCn() Lang {
	lang := &CnLang{
		langMap: make(map[LanguageKey]string),
	}
	//ddz_client
	lang.langMap[InvalidOperation] = "操作无效"
	lang.langMap[LobbyMessage] = "[大厅消息] %s"
	lang.langMap[ClientMessage] = "[客户端消息] %s"
	lang.langMap[GameMessage] = "[游戏消息] %s"
	lang.langMap[RoomMessage] = "[房间消息] %s"
	lang.langMap[Hand] = "手牌如下："
	lang.langMap[HelpInfo] = "输入 h 查看帮助信息"
	lang.langMap[CurrentUser] = "当前用户 %s"
	lang.langMap[LinkSuccess] = "已连接至服务器 %s:%d"
	lang.langMap[InvalidInput] = "输入无效"
	lang.langMap[InvalidInput] = "命令无效"
	//ddz_mfunc
	lang.langMap[CurrentLandlord] = "当前地主[%s]"
	lang.langMap[NotElectLandlord] = "未选出地主"
	lang.langMap[EstablishRoom] = "房间[%s]已创建"
	lang.langMap[IntoRoom] = "用户[%s]加入房间"
	lang.langMap[NotFindRoom] = "未找到对应房间"
	lang.langMap[QuitRoom] = "退出房间"
	lang.langMap[PreparationFailed] = "无法准备(未知原因)"
	lang.langMap[UserPrepare] = "用户[%s]已准备"
	lang.langMap[CancelPrepared] = "用户[%s]取消准备"
	lang.langMap[UserQuitRoom] = "用户[%s]退出房间"
	lang.langMap[Homeowner] = "用户[%s]成为房主"
	lang.langMap[IsFull] = "房间已满"
	lang.langMap[RepeatIntoRoom] = "已在房间内无法加入房间"
	lang.langMap[CanNotEstablishRoom] = "已在房间内无法创建新房间"
	lang.langMap[QuitRoomFailed] = "不在任何房间内无法退出"
	lang.langMap[GameStart] = "对局开始"
	lang.langMap[RoomClosed] = "房间关闭"
	lang.langMap[OperateTimeOut] = "操作时间用尽"
	lang.langMap[TimeRemain] = "操作时间还剩%s秒"
	lang.langMap[Operate] = "***请操作***"
	lang.langMap[TurnToUser] = "轮到[%s]操作"
	lang.langMap[WhetherSeizeLandlord] = "是否抢地主 (y/n)"
	lang.langMap[TrustSeizeLandlord] = "[%s]托管操作不抢地主"
	lang.langMap[SeizeLandlord] = "用户[%s]抢地主"
	lang.langMap[NotSeizeLandlord] = "用户[%s]不抢地主"
	lang.langMap[LandlordUser] = "***地主用户[%s]***"
	lang.langMap[HoleCard] = "底牌:"
	lang.langMap[TurnToDiscard] = "***[%s]出牌***"
	lang.langMap[SkipToDiscard] = "[%s]跳过出牌"
	lang.langMap[AttentionRemain] = "***注意[%s]还剩%d张手牌***"
	lang.langMap[UserWin] = "***注意[%s]还剩%d张手牌***"
	lang.langMap[WinnerLandlord] = "***[%s]为优胜者成为地主***"
	lang.langMap[TrustToDiscard] = "[%s]托管操作出牌"
	lang.langMap[UserOperateTimeOut] = "[%s]操作超时将由系统托管操作"
	lang.langMap[GameOver] = "对局结束"
	lang.langMap[NextGame] = "可进行下一场对局"
	lang.langMap[CurrentRoomInfo] = "当前房间信息\n"
	lang.langMap[AllRoomInfo] = "所有房间信息\n"
	lang.langMap[UserCardRemain] = "用户剩余卡牌\n"
	lang.langMap[CardNum] = "%s[%d]张 "
	return lang
}

func (l *CnLang) Get(key LanguageKey) string {
	return l.langMap[key]
}
