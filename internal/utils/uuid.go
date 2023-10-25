package utils

import (
	"unsafe"
)

// go clean -testcache // Delete all cached test results

const (
	suuidStr  = "23456789abcdefghijkmnpqrstuvwxyz" //32
	suuidMask = 1<<5 - 1                           //11111
	_suid     = 8                                  //LEQ 12
	_muid     = 16                                 //LEQ 16

	uuidStr  = "0123456789abcdef" //16
	uuidMask = 1<<4 - 1           //1111

)

// SUID SUID生成
func SUID() string {
	buf := make([]byte, _suid)
	for idx, cache := 0, fastRand(); idx < _suid; {
		buf[idx] = suuidStr[cache&suuidMask]
		cache >>= 5
		idx++
	}
	return unsafe.String(&buf[0], len(buf))
}

// MUID MUID生成
func MUID() string {
	buf := make([]byte, _muid)
	for idx, cache := 0, fastRand(); idx < _muid; {
		buf[idx] = uuidStr[cache&uuidMask]
		cache >>= 4
		idx++
	}
	return unsafe.String(&buf[0], len(buf))
}

// UUID
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func UUID() string {
	buf := make([]byte, 36)
	for idx, cache, remain := 0, fastRand(), 16; idx < 32; {
		if remain == 0 {
			cache, remain = fastRand(), 16
		}
		buf[idx] = uuidStr[cache&uuidMask]
		cache >>= 4
		remain--
		idx++
	}
	buf[32] = buf[8]
	buf[33] = buf[13]
	buf[34] = buf[18]
	buf[35] = buf[23]
	buf[8] = '-'
	buf[13] = '-'
	buf[14] = '4'
	buf[18] = '-'
	buf[23] = '-'
	// return *(*string)(unsafe.Pointer(&buf))
	return unsafe.String(&buf[0], len(buf))
}

//go:linkname fastRand runtime.fastrand64
func fastRand() uint64

//go:noescape
//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)

// for runtime.walltime
func Now() int64 {
	sec, nsec, _ := now()
	return sec*1e3 + int64(nsec)/1e6
	// return time.Now().Unix()
}
