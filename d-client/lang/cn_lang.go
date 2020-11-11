package lang

type CnLang struct {
	langMap map[LanguageKey]string
}

func NewCn() Lang {
	lang := &CnLang{
		langMap: make(map[LanguageKey]string),
	}
	lang.langMap[InvalidOperation] = "操作无效"
	lang.langMap[LobbyMessage] = "[大厅消息] %s"
	lang.langMap[ClientMessage] = "[客户端消息] %s"
	lang.langMap[GameMessage] = "[游戏消息] %s"
	lang.langMap[RoomMessage] = "[房间消息] %s"
	return lang
}

func (l *CnLang) Get(key LanguageKey) string {
	return l.langMap[key]
}
