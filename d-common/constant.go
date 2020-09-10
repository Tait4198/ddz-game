package d_common

const (
	CenterLevel MessageLevel = "center"
	RoomLevel                = "room"
	GameLevel                = "game"
	ClientLevel              = "client"
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
