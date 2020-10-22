package d_common

// 用于指定调用方法
type MessageType uint

//消息发往 center/room
type MessageLevel string

type ResultRoom struct {
	Id        uint     `json:"id"`
	Homeowner string   `json:"homeowner"`
	IsRun     bool     `json:"is_run"`
	Clients   []string `json:"clients"`
}
