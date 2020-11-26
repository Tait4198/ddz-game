package d_common

import (
	"fmt"
	"sort"
)

type DdzMessageType uint

type DdzPokerType uint

// 游戏阶段
type GameStage uint

type DdzPokerResult struct {
	PkType DdzPokerType
	Score  uint
	Len    uint
}

type UserPlayInfo struct {
	Pokers    []Poker `json:"pokers"`
	Remaining int     `json:"remaining"`
	Name      string  `json:"name"`
}

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

func SortByScore(p, q *Poker) bool {
	return p.Score < q.Score
}

func (p *Poker) String() string {
	return fmt.Sprintf("%s-%s", p.Level, p.Suit)
}
