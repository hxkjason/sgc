package rand_service

import (
	"math/rand"
	"time"
)

const (
	LetterNumericCharList = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberCharList        = "0123456789"
	LetterCharList        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenRandomString 随机生成字符串
func GenRandomString(charsetList string, length int) string {

	b := make([]byte, length, length)
	charLen := len(charsetList)
	if charLen <= 0 {
		panic("待选字符不能为空")
	}
	for i := range b {
		b[i] = charsetList[seededRand.Intn(charLen)]
	}
	return string(b)
}
