package utils

import (
	"math/rand"
	"time"
)

// GetRandInt 获取 startNum - endNum - 1 之间的随机数
func GetRandInt(startNum, endNum int) int {
	if startNum >= endNum {
		return endNum
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(endNum-startNum) + startNum
}
