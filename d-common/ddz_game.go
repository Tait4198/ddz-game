package d_common

import "sort"

const (
	Invalid DdzPokerType = iota
	Single
	Double
	Three
	ThreeWithOne
	FourWithTwo
	Straight
	ContDouble
	ContThree
	Aircraft
	Bomb
	KingBomb
)

// 0 a 1 b -1 无效
func ComparePoker(aPks, bPks []Poker) int {
	ar := GetPokerType(aPks)
	br := GetPokerType(bPks)
	if ar.PkType == Invalid || br.PkType == Invalid {
		return -1
	}
	if ar.PkType == KingBomb {
		return 0
	} else if br.PkType == KingBomb {
		return 1
	}
	if ar.PkType == Bomb && br.PkType == Bomb {
		if ar.Score > br.Score {
			return 0
		} else if br.Score > ar.Score {
			return 1
		}
	}
	if ar.PkType == Bomb {
		return 0
	} else if br.PkType == Bomb {
		return 1
	}
	if ar.PkType == br.PkType && ar.Len == br.Len {
		if ar.Score > br.Score {
			return 0
		} else if br.Score > ar.Score {
			return 1
		}
	}
	return -1
}

func GetPokerType(pks []Poker) DdzPokerResult {
	pkLen := len(pks)
	switch pkLen {
	case 1:
		return DdzPokerResult{Single, pks[0].Score, 1}
	case 2:
		if pks[0].Score == pks[1].Score {
			return DdzPokerResult{Double, pks[0].Score, 1}
		} else {
			kf := true
			for _, pk := range pks {
				if pk.Score != 13 && pk.Score != 14 {
					kf = false
				}
			}
			if kf {
				return DdzPokerResult{KingBomb, 0, 1}
			}
		}
	case 3:
		three := true
		tempScore := pks[0].Score
		for _, pk := range pks {
			if pk.Score != tempScore {
				three = false
			}
		}
		if three {
			return DdzPokerResult{Three, pks[0].Score, 1}
		}
	case 4:
		pkMap := getPkMap(pks)
		if len(pkMap) == 1 {
			for _, v := range pkMap {
				if v == 4 {
					return DdzPokerResult{Bomb, pks[0].Score, 1}
				}
			}
		} else if len(pkMap) == 2 {
			one := false
			three := false
			threeK := uint(0)
			for k, v := range pkMap {
				if v == 1 {
					one = true
				} else if v == 3 {
					three = true
					threeK = k
				}
			}
			if one && three {
				return DdzPokerResult{ThreeWithOne, threeK, 1}
			}
		}
	default:
		SortPoker(pks, SortByScore)
		pkMap := getPkMap(pks)
		if r := checkFourWithTwo(pkMap, pks); r.PkType != Invalid {
			return r
		}
		if r := checkStraight(pkMap, pks); r.PkType != Invalid {
			return r
		}
		if r := checkContDouble(pkMap, pks); r.PkType != Invalid {
			return r
		}
		if r := checkContThree(pkMap, pks); r.PkType != Invalid {
			return r
		}
		if r := checkAircraft(pkMap, pks); r.PkType != Invalid {
			return r
		}
	}
	return DdzPokerResult{Invalid, 0, 0}
}

func getPkMap(pks []Poker) map[uint]uint {
	pkMap := make(map[uint]uint)
	for _, pk := range pks {
		if v, ok := pkMap[pk.Score]; ok {
			pkMap[pk.Score] = v + 1
		} else {
			pkMap[pk.Score] = 1
		}
	}
	return pkMap
}

func checkFourWithTwo(pkMap map[uint]uint, pks []Poker) DdzPokerResult {
	pkMapLen := len(pkMap)
	// AAAABC / AAAABB
	if pkMapLen == 2 || pkMapLen == 3 {
		for k, v := range pkMap {
			if v == 4 {
				return DdzPokerResult{FourWithTwo, k, 1}
			}
		}
	}
	return DdzPokerResult{Invalid, 0, 0}
}

func checkStraight(pkMap map[uint]uint, pks []Poker) DdzPokerResult {
	pkLen := len(pks)
	invalid := DdzPokerResult{Invalid, 0, 0}
	if pkLen >= 5 && pkLen <= 13 {
		for _, v := range pkMap {
			if v != 1 {
				return invalid
			}
		}
		for i := 0; i < pkLen-1; i++ {
			if pks[i].Score+1 != pks[i+1].Score {
				return invalid
			}
		}
		return DdzPokerResult{Straight, pks[0].Score, uint(pkLen)}
	}
	return invalid
}

func checkContDouble(pkMap map[uint]uint, pks []Poker) DdzPokerResult {
	pkLen := len(pks)
	invalid := DdzPokerResult{Invalid, 0, 0}
	if pkLen%2 == 0 && pkLen >= 6 && pkLen <= 20 {
		for _, v := range pkMap {
			if v != 2 {
				return invalid
			}
		}
		for i := 0; i < pkLen-2; i += 2 {
			if pks[i].Score != pks[i+1].Score || pks[i].Score+1 != pks[i+2].Score {
				return invalid
			}
		}
		return DdzPokerResult{ContDouble, pks[0].Score, uint(pkLen / 2)}
	}
	return invalid
}

func checkContThree(pkMap map[uint]uint, pks []Poker) DdzPokerResult {
	pkLen := len(pks)
	invalid := DdzPokerResult{Invalid, 0, 0}
	if pkLen%3 == 0 && pkLen >= 9 && pkLen <= 18 {
		for _, v := range pkMap {
			if v != 3 {
				return invalid
			}
		}
		for i := 0; i < pkLen-3; i += 3 {
			if pks[i].Score != pks[i+1].Score || pks[i].Score != pks[i+2].Score ||
				pks[i+1].Score != pks[i+2].Score || pks[i].Score+1 != pks[i+3].Score {
				return invalid
			}
		}
		return DdzPokerResult{ContThree, pks[0].Score, uint(pkLen / 3)}
	}
	return invalid
}

func checkAircraft(pkMap map[uint]uint, pks []Poker) DdzPokerResult {
	pkLen := len(pks)
	if pkLen%4 == 0 && pkLen >= 12 && pkLen <= 20 {
		aSize := pkLen / 4
		bSize := 0
		var threeSl []int
		for k, v := range pkMap {
			if v == 3 {
				bSize++
				threeSl = append(threeSl, int(k))
			}
		}
		sort.Ints(threeSl)
		if bSize >= aSize {
			threeNum := 1
			lastThreeNum := 0
			minStart := 9999
			for i := 0; i < len(threeSl)-1; i++ {
				if threeSl[i]+1 != threeSl[i+1] {
					if lastThreeNum > threeNum {
						lastThreeNum = threeNum
					}
					threeNum = 1
					minStart = 9999
				} else {
					if threeSl[i] <= minStart {
						minStart = threeSl[i]
					}
					threeNum++
				}
			}
			if threeNum == aSize || lastThreeNum == aSize {
				return DdzPokerResult{Aircraft, uint(minStart), uint(aSize)}
			} else if threeNum > aSize && threeSl != nil {
				return DdzPokerResult{Aircraft, uint(threeSl[threeNum-aSize]), uint(aSize)}
			}
		}
	}
	return DdzPokerResult{Invalid, 0, 0}
}
