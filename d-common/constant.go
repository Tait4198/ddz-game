package d_common

const (
	CenterLevel MessageLevel = "center"
	RoomLevel   MessageLevel = "room"
	GameLevel   MessageLevel = "game"
	ClientLevel MessageLevel = "client"
)

const (
	ClientRegister   MessageType = 1000
	ClientUnregister MessageType = 1001
	RoomCreate       MessageType = 1100
	RoomDisband      MessageType = 1101
	RoomJoin         MessageType = 1102
	RoomQuit         MessageType = 1103
	RoomClose        MessageType = 1104
	RoomGameMessage  MessageType = 1105
	RoomReady        MessageType = 1106
	RoomCancelReady  MessageType = 1107
	RoomUnableCreate MessageType = 1108
	RoomAlreadyIn    MessageType = 1109
	RoomFull         MessageType = 1110
	RoomInvalid      MessageType = 1111
	RoomUnableExit   MessageType = 1112
	RoomRun          MessageType = 1113
	RoomMissUser     MessageType = 1114
	RoomSomeoneQuit  MessageType = 1115
	RoomNewHomeowner MessageType = 1116
	GetAllRoomInfo   MessageType = 1200
	GetCurRoomInfo   MessageType = 1201
)
