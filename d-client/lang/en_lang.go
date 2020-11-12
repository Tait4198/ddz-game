package lang

type EnLang struct {
	langMap map[LanguageKey]string
}

func NewEn() Lang {
	lang := &EnLang{
		langMap: make(map[LanguageKey]string),
	}
	//ddz_client
	lang.langMap[InvalidOperation] = "Invalid Operation"
	lang.langMap[LobbyMessage] = "[Lobby Message] %s"
	lang.langMap[ClientMessage] = "[Room Message] %s"
	lang.langMap[GameMessage] = "[Game Message] %s"
	lang.langMap[RoomMessage] = "[Room Message] %s"
	lang.langMap[Hand] = "Hand："
	lang.langMap[HelpInfo] = "input 'h' for checking help-info"
	lang.langMap[CurrentUser] = "current user %s"
	lang.langMap[LinkSuccess] = "Server connection successful %s:%d"
	lang.langMap[InvalidInput] = "invalid input"
	lang.langMap[InvalidCommand] = "invalid command"
	//ddz_mfunc
	lang.langMap[CurrentLandlord] = "current landlord[%s]"
	lang.langMap[NotElectLandlord] = "landlord were not elected"
	lang.langMap[EstablishRoom] = "room [%s] was established"
	lang.langMap[IntoRoom] = "user [%s] get into the room"
	lang.langMap[NotFindRoom] = "didn't find target room"
	lang.langMap[QuitRoom] = "quit room"
	lang.langMap[PreparationFailed] = "preparation failed"
	lang.langMap[UserPrepare] = "user [%s] prepared"
	lang.langMap[CancelPrepared] = "usr [%s] cancel prepared"
	lang.langMap[UserQuitRoom] = "user [%s] quit room"
	lang.langMap[Homeowner] = "user [%s] is Homeowner now"
	lang.langMap[IsFull] = "room is full"
	lang.langMap[RepeatIntoRoom] = "can't get into other room when you stay in a room"
	lang.langMap[CanNotEstablishRoom] = "can't establish room when you stay in a room"
	lang.langMap[QuitRoomFailed] = "can't quit room when you are not in a room"
	lang.langMap[GameStart] = "game start"
	lang.langMap[RoomClosed] = "room closed"
	lang.langMap[OperateTimeOut] = "operate time out"
	lang.langMap[TimeRemain] = "%s second remain"
	lang.langMap[Operate] = "*** make your step ***"
	lang.langMap[TurnToUser] = "turn [%s] to operate"
	lang.langMap[WhetherSeizeLandlord] = "do you want seize landlord?(y/n)"
	lang.langMap[TrustSeizeLandlord] = "[%s] trust not to seize landlord"
	lang.langMap[SeizeLandlord] = "user [%s] seize landlord"
	lang.langMap[NotSeizeLandlord] = "user [%s] didn't seize landlord"
	lang.langMap[LandlordUser] = "*** landlord user: [%s] ***"
	lang.langMap[HoleCard] = "hole card: "
	lang.langMap[TurnToDiscard] = "*** turn [%s] to discard ***"
	lang.langMap[SkipToDiscard] = "[%s] skip discarding"
	lang.langMap[AttentionRemain] = "*** Attention! user[%s] %d card remain ***"
	lang.langMap[UserWin] = "*** user [%s] win ***"
	lang.langMap[WinnerLandlord] = "*** winner[%s] will become landlord next game ***"
	lang.langMap[TrustToDiscard] = "user[%s] trust to discard"
	lang.langMap[UserOperateTimeOut] = "user[%s] operate time out, trust to system"
	lang.langMap[GameOver] = "game over"
	lang.langMap[NextGame] = "could play next game"
	lang.langMap[CurrentRoomInfo] = "current room info \n"
	lang.langMap[AllRoomInfo] = "all room info \n"
	lang.langMap[UserCardRemain] = "users card remain \n"
	lang.langMap[CardNum] = "user[%s] %d "
	return lang
}

func (l *EnLang) Get(key LanguageKey) string {
	return l.langMap[key]
}
