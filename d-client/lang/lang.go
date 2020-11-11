package lang

type Lang interface {
	Get(key LanguageKey) string
}

type LanguageKey int

const (
	InvalidOperation LanguageKey = iota
	LobbyMessage
	RoomMessage
	GameMessage
	ClientMessage
)
