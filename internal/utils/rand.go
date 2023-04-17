package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	digit     = "0123456789"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	chars     = digit + lowercase
	charsLen  = len(chars)
	digitLen  = len(digit)
	charsMask = 1<<6 - 1
	digitMask = 1<<4 - 1
)

var rng = rand.NewSource(time.Now().UnixNano())

// RandStr 返回指定长度的随机字符串
// 包含数字 小写字母
func Str(ln int) string {
	/* chars 36个字符
	 * rng.Int63() 每次产出64bit的随机数,每次我们使用6bit(2^6=64) 可以使用10次
	 */
	buf := make([]byte, ln)
	for idx, cache, remain := 0, rng.Int63(), 10; idx < ln; {
		if remain == 0 {
			cache, remain = rng.Int63(), 10
		}
		buf[idx] = chars[int(cache&charsMask)%charsLen]
		cache >>= 6
		remain--
		idx++
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// RandDigitStr 返回指定长度的随机字符串
// 只包含数字
func Digit(ln int) string {
	/* digits 10个字符
	 * fastRand() 每次产出34bit的随机数,每次我们使用4bit(2^4=16) 可以使用8次
	 */
	buf := make([]byte, ln)
	for idx, cache, remain := 0, fastRand(), 8; idx < ln; {
		if remain == 0 {
			cache, remain = fastRand(), 8
		}
		buf[idx] = chars[int(cache&digitMask)%digitLen]
		cache >>= 4
		remain--
		idx++
	}
	return *(*string)(unsafe.Pointer(&buf))
}
