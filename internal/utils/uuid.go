package utils

import (
	"unsafe"
)

// go clean -testcache // Delete all cached test results

const (
	suuidStr  = "23456789abcdefghijkmnpqrstuvwxyz" //32
	suuidMask = 1<<5 - 1                           //11111

	uuidStr  = "0123456789abcdef" //16
	uuidMask = 1<<4 - 1           //1111

)

// UUID UUID生成-16bit
// 10-times '-' 6-random string '-' 6-random string
func SUID() string {
	buf := make([]byte, 16)
	idx := 15
	// 15-begin 32*32*32*32*32*32==1073741824
	rand := fastRand() //6
	for ; idx > 9; idx-- {
		buf[idx] = suuidStr[rand&suuidMask]
		rand >>= 5
	}
	now := Now() //10
	// 50bit
	// 10000-01-01 01:00:00 =>  253402246800000 ms
	// 111001100111011111001110111001111110001010000000 48bit
	for idx = 9; idx >= 0; idx-- {
		buf[idx] = suuidStr[now&suuidMask]
		now >>= 5
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// UUID
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func UUID() string {
	buf := make([]byte, 36)
	for idx, cache, remain := 12, fastRand(), 8; idx < 36; {
		if remain == 0 {
			cache, remain = fastRand(), 8
		}
		buf[idx] = uuidStr[cache&uuidMask]
		cache >>= 4
		remain--
		idx++
	}
	// 50bit
	// 10000-01-01 01:00:00 =>  253402246800000 ms
	// 111001100111011111001110111001111110001010000000 48bit
	now := Now() >> 4
	for idx := 9; idx >= 0; idx-- {
		buf[idx] = uuidStr[now&uuidMask]
		now >>= 4
	}
	buf[8] = '-'
	buf[10] = buf[13]
	buf[11] = buf[14]
	buf[13] = '-'
	buf[14] = '4'
	buf[18] = '-'
	buf[23] = '-'
	// return *(*string)(unsafe.Pointer(&buf))
	return unsafe.String(&buf[0], len(buf))
}

//go:linkname fastRand runtime.fastrand
func fastRand() uint32

//go:noescape
//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)

// for runtime.walltime
func Now() int64 {
	sec, nsec, _ := now()
	return sec*1e3 + int64(nsec)/1e6
	// return time.Now().Unix()
}
