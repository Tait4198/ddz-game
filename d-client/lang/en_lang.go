package lang

type EnLang struct {
	langMap map[LanguageKey]string
}

func NewEn() Lang {
	lang := &EnLang{
		langMap: make(map[LanguageKey]string),
	}
	lang.langMap[InvalidOperation] = "Invalid Operation"
	lang.langMap[LobbyMessage] = "[Lobby Message] %s"
	lang.langMap[ClientMessage] = "[Room Message] %s"
	lang.langMap[GameMessage] = "[Game Message] %s"
	lang.langMap[RoomMessage] = "[Room Message] %s"
	return lang
}

func (l *EnLang) Get(key LanguageKey) string {
	return l.langMap[key]
}
