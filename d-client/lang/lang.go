package lang

type Lang interface {
	Get(key LanguageKey) string
}

type LanguageKey int

const (
	InvalidOperation LanguageKey = iota
	//ddz_client
	LobbyMessage
	RoomMessage
	GameMessage
	ClientMessage
	Hand
	HelpInfo
	CurrentUser
	LinkSuccess
	InvalidInput
	InvalidCommand
	//ddz_mfunc
	CurrentLandlord
	NotElectLandlord
	EstablishRoom
	IntoRoom
	NotFindRoom
	QuitRoom
	PreparationFailed
	UserPrepare
	CancelPrepared
	UserQuitRoom
	Homeowner
	IsFull
	RepeatIntoRoom
	CanNotEstablishRoom
	QuitRoomFailed
	GameStart
	RoomClosed
	OperateTimeOut
	TimeRemain
	Operate
	TurnToUser
	WhetherSeizeLandlord
	TrustSeizeLandlord
	SeizeLandlord
	NotSeizeLandlord
	LandlordUser
	HoleCard
	TurnToDiscard
	SkipToDiscard
	AttentionRemain
	UserWin
	WinnerLandlord
	TrustToDiscard
	UserOperateTimeOut
	GameOver
	NextGame
	CurrentRoomInfo
	AllRoomInfo
	UserCardRemain
	CardNum
	Landlord
	Farmer
)
