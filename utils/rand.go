package utils

import (
	"math/rand"
	"time"
)

func getRander() *rand.Rand {
	seed := time.Now().UnixNano()
	return rand.New(rand.NewSource(seed))
}

// GetRandIdList 从id列表list中随机选取number个元素
func GetRandIdList(list []int64, number int) []int64 {
	r := getRander()
	if number > len(list) {
		return list
	}
	r.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list[:number]
}
