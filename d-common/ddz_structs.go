package d_common

import (
	"fmt"
	"sort"
)

type DdzMessageType uint

type Poker struct {
	Suit  string
	Level string
	Score uint
}

type PokerWrapper struct {
	pks []Poker
	by  func(p, q *Poker) bool
}

type SortBy func(p, q *Poker) bool

func (pw PokerWrapper) Len() int {
	return len(pw.pks)
}

func (pw PokerWrapper) Swap(i, j int) {
	pw.pks[i], pw.pks[j] = pw.pks[j], pw.pks[i]
}

func (pw PokerWrapper) Less(i, j int) bool {
	return pw.by(&pw.pks[i], &pw.pks[j])
}

func SortPoker(pks []Poker, by SortBy) {
	sort.Sort(PokerWrapper{pks, by})
}

func (p *Poker) String() string {
	return fmt.Sprintf("%s-%s", p.Level, p.Suit)
}
