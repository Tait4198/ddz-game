package main

import (
	cm "com.github/gc-common"
	"sort"
)

type PkScoreSize struct {
	Score uint
	Size  uint
}

type PkScoreSizeWrapper struct {
	ss []PkScoreSize
	by func(p, q *PkScoreSize) bool
}

type SortBy func(p, q *PkScoreSize) bool

func (pw PkScoreSizeWrapper) Len() int {
	return len(pw.ss)
}

func (pw PkScoreSizeWrapper) Swap(i, j int) {
	pw.ss[i], pw.ss[j] = pw.ss[j], pw.ss[i]
}

func (pw PkScoreSizeWrapper) Less(i, j int) bool {
	return pw.by(&pw.ss[i], &pw.ss[j])
}

func SortPkScoreSource(pks []PkScoreSize, by SortBy) {
	sort.Sort(PkScoreSizeWrapper{pks, by})
}

func PkAutoPlay(prevPks, pks []cm.Poker) []int {
	var resultIdx []int
	pkMap := cm.GetPkMap(pks)
	var ss []PkScoreSize
	for score, size := range pkMap {
		ss = append(ss, PkScoreSize{Score: score, Size: size})
	}
	SortPkScoreSource(ss, func(p, q *PkScoreSize) bool {
		return p.Score < q.Score
	})

	if len(prevPks) == 0 || len(pks) == 0 {
		if len(pks) > 0 {
			return simpleCase(1, 0, ss, pks)
		} else {
			return resultIdx
		}
	}
	dpr := cm.GetPokerType(prevPks)
	if dpr.PkType == cm.KingBomb {
		return resultIdx
	}
	if len(prevPks) <= len(pks) {
		switch dpr.PkType {
		case cm.Single, cm.Double, cm.Three, cm.Bomb:
			resultIdx = simpleCase(uint(len(prevPks)), dpr.Score, ss, pks)
		case cm.ThreeWithOne:
			resultIdx = withCase(3, 1, dpr.Score, ss, pks)
		case cm.FourWithTwo:
			resultIdx = withCase(4, 2, dpr.Score, ss, pks)
		case cm.Straight:
			resultIdx = contCase(1, false, dpr, ss, pks)
		case cm.ContDouble:
			resultIdx = contCase(2, false, dpr, ss, pks)
		case cm.ContThree:
			resultIdx = contCase(3, false, dpr, ss, pks)
		case cm.Aircraft:
			resultIdx = contCase(3, true, dpr, ss, pks)
		}
	}
	if len(resultIdx) == 0 {
		if dpr.PkType != cm.Bomb {
			resultIdx = simpleCase(4, 0, ss, pks)
		}
	}
	if len(resultIdx) == 0 {
		_, sOk := pkMap[cm.SScore]
		_, xOk := pkMap[cm.XScore]
		if sOk && xOk {
			for i, pk := range pks {
				if pk.Score == cm.SScore || pk.Score == cm.XScore {
					resultIdx = append(resultIdx, i)
				}
			}
		}
	}
	return resultIdx
}

func simpleCase(pkSize, dprScore uint, ss []PkScoreSize, pks []cm.Poker) []int {
	var resultIdx []int
	var useScore = matchScore(pkSize, dprScore, ss)
	if useScore > 0 {
		cPkMap := make(map[int]byte)
		for i, pk := range pks {
			if _, ok := cPkMap[i]; !ok && pk.Score == useScore {
				resultIdx = append(resultIdx, i)
				cPkMap[i] = 0
				if len(cPkMap) == int(pkSize) {
					break
				}
			}
		}
	}
	return resultIdx
}

func withCase(pkSize, withSize, dprScore uint, ss []PkScoreSize, pks []cm.Poker) []int {
	var resultIdx []int
	fourScore := matchAccScore(pkSize, dprScore, ss, true)
	if fourScore == 0 {
		return resultIdx
	}
	cPkMap := make(map[int]byte)
	for i, pk := range pks {
		if _, ok := cPkMap[i]; !ok && pk.Score == fourScore {
			resultIdx = append(resultIdx, i)
			cPkMap[i] = 0
			if uint(len(cPkMap)) == pkSize {
				break
			}
		}
	}
	var index uint = 0
	for i := range pks {
		if _, ok := cPkMap[i]; !ok {
			resultIdx = append(resultIdx, i)
			index++
			if index == withSize {
				break
			}
		}
	}
	return resultIdx
}

func contCase(pkSize uint, with bool, dpr cm.DdzPokerResult, ss []PkScoreSize, pks []cm.Poker) []int {
	var resultIdx []int
	var beforeScore uint = 0
	var ssMap = make(map[uint]byte)
	for _, s := range ss {
		// 小于S
		if s.Score > dpr.Score && s.Score < cm.SScore && s.Size >= pkSize {
			if beforeScore == 0 || s.Score-beforeScore == 1 {
				beforeScore = s.Score
				ssMap[s.Score] = 0
				if len(ssMap) == int(dpr.Len) {
					break
				}
			}
		}
	}
	if len(ssMap) != int(dpr.Len) {
		return []int{}
	}

	iPkMap := make(map[uint]uint)
	cPkMap := make(map[int]uint)
	for i, pk := range pks {
		if _, ok := ssMap[pk.Score]; !ok {
			continue
		}
		var vs uint = 0
		if v, ok := iPkMap[pk.Score]; (ok && v < pkSize) || !ok {
			vs = v
		} else {
			continue
		}
		if _, ok := cPkMap[i]; !ok {
			resultIdx = append(resultIdx, i)
			iPkMap[pk.Score] = vs + 1
			cPkMap[i] = 0
		}
	}

	if with {
		var index uint = 0
		for i := range pks {
			if _, ok := cPkMap[i]; !ok {
				resultIdx = append(resultIdx, i)
				index++
				if index == dpr.Len {
					break
				}
			}
		}
	}

	return resultIdx
}

// 匹配最接近的评分 (允许拆分)
func matchAccScore(size, minScore uint, ss []PkScoreSize, accurate bool) uint {
	var score uint = 0
	for _, s := range ss {
		if s.Size == size && s.Score > minScore {
			score = s.Score
			break
		}
	}
	if score == 0 && !accurate {
		for _, s := range ss {
			if s.Size >= size && s.Score > minScore {
				score = s.Score
				break
			}
		}
	}
	return score
}

// 匹配最接近的评分
func matchScore(size, minScore uint, ss []PkScoreSize) uint {
	return matchAccScore(size, minScore, ss, false)
}
