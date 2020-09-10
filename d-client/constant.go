package main

const (
	CenterLevel = "center"
	RoomLevel   = "room"
	GameLevel   = "game"
	ClientLevel = "client"
)

const (
	StageGrabLandlord GameStage = 1 //抢地主阶段
)

const (
	ClientRegister   MessageType = 1000
	ClientUnregister             = 1001
	RoomCreate                   = 1100
	RoomDisband                  = 1101
	RoomJoin                     = 1102
	RoomQuit                     = 1103
	RoomClose                    = 1104
	RoomGameMessage              = 1105
	RoomReady                    = 1106
	RoomCancelReady              = 1107
	RoomUnableCreate             = 1108
	RoomAlreadyIn                = 1109
	RoomFull                     = 1110
	RoomInvalid                  = 1111
	RoomUnableExit               = 1112
	RoomRun                      = 1113
	RoomMissUser                 = 1114
	RoomSomeoneQuit              = 1115
	RoomNewHomeowner             = 1116
)

const (
	GameStop             MessageType = 2000
	GameStart                        = 2001
	GameRestart                      = 2002
	GameCountdown                    = 2003
	GameInvalidOps                   = 2004
	GameOpsTimeout                   = 2005
	GameNextUserOps                  = 2006
	GameExe                          = 2010
	GameDealPoker                    = 2100
	GameGrabLandlord                 = 2101
	GameNGrabLandlord                = 2102
	GameGrabLandlordEnd              = 2103
	GameGrabHostingOps               = 2104
	GameNoGrabLandlord               = 2105
	GameNewLandlord                  = 2105
	GameWaitGrabLandlord             = 2106
)
