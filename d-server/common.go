package main

import (
	"math/rand"
	"time"
)

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func RandIntMap(min, max, size int) map[int]bool {
	m := make(map[int]bool)
	for {
		i := RandInt(min, max)
		if _, ok := m[i]; !ok {
			m[i] = true
		}
		if len(m) == size {
			break
		}
	}
	return m
}
