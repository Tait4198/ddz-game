package d_common

import "fmt"

// 用于指定调用方法
type MessageType uint

//消息发往 center/room
type MessageLevel string

type ResultRoom struct {
	Id        uint           `json:"id"`
	Homeowner string         `json:"homeowner"`
	IsRun     bool           `json:"is_run"`
	Clients   []ResultClient `json:"clients"`
}

func (rr *ResultRoom) String() string {
	isRun := "正在等待"
	if rr.IsRun {
		isRun = "正在对局"
	}
	clientNames := ""
	for _, rc := range rr.Clients {
		clientNames += rc.String()
	}
	return fmt.Sprintf("Id: %d %s 房主: %s 用户: %s\n", rr.Id, isRun, rr.Homeowner, clientNames)
}

type ResultClient struct {
	Ready    bool   `json:"ready"`
	Username string `json:"username"`
}

func (rc *ResultClient) String() string {
	if rc.Ready {
		return rc.Username + " [已准备] "
	} else {
		return rc.Username + " [未准备] "
	}
}
