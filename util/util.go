package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	FormatDateTime = "2006-01-02 15:04:05"

	// FmtyyyyMMdd 年月日
	FmtyyyyMMdd = "20060102"
	chars       = "0123456789abcdefghijklmnopqrstuvwxyz" //ABCDEFGHIJKLMNOPQRSTUVWXYZ
	charsLen    = len(chars)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandStr 返回指定长度的随机字符串
func RandStr(ln int) string {
	buf := strings.Builder{}
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < ln; i++ {
		buf.WriteByte(chars[rand.Intn(charsLen)])
	}
	return buf.String()
}
