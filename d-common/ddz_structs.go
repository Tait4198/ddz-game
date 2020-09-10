package d_common

import "fmt"

type DdzMessageType uint

type Poker struct {
	Suit  string
	Level string
	Score uint
}

func (p *Poker) String() string {
	return fmt.Sprintf("%s-%s", p.Level, p.Suit)
}
