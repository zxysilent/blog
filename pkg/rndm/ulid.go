package rndm

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"runtime"
	"sync"
)

// ULID 说明
//
// 128bit = 48bit 毫秒时间戳 + 80bit 随机数
// 编码为 26 位 Crockford Base32（字符集排除 I L O U）
// 字符集 "0123456789abcdefghjkmnpqrstvwxyz" 本身 ASCII 升序，
// 因此字典序 == 时间序，ORDER BY id / Desc("id") 直接等价于按创建时间排序。
//
// 输出恒为小写：本包生成的 ULID 只使用小写字母，全库主键大小写统一，
// 大小写混写会导致唯一索引在 MySQL 默认排序规则下看似重复、在 SQLite 下真的重复。
// 输入侧仍兼容大写（Crockford Base32 本就不区分大小写，见 IsULID）。
//
// 同一毫秒内通过 entropy 自增保证严格单调递增（entropy 溢出则自旋到下一毫秒）。
const (
	ulidChars   = "0123456789abcdefghjkmnpqrstvwxyz" //32 个字符，ASCII 升序
	ulidEncLen  = 26                                 //编码后长度
	ulidTimeLen = 6                                  //48bit 时间戳字节数
	ulidRandLen = 10                                 //80bit 随机数字节数
)

var (
	ulidMu      sync.Mutex
	ulidLastMs  int64             //上一次使用的时间戳(ms)
	ulidEntropy [ulidRandLen]byte //当前毫秒内的 80bit 随机数
	ulidSeeded  bool              //是否已初始化
)

// ULID 返回 26 位 ULID 字符串
// 48bit 毫秒时间戳 + 80bit 随机数，同毫秒内严格单调递增
// 输出恒为小写（Crockford Base32 小写字符集）
func ULID() string {
	ulidMu.Lock()
	ms := Now()
	switch {
	case !ulidSeeded:
		// 首次调用：填充随机种子
		_, _ = crand.Read(ulidEntropy[:])
		ulidSeeded = true
		ulidLastMs = ms
	case ms > ulidLastMs:
		// 进入新的毫秒：重新填充随机数
		_, _ = crand.Read(ulidEntropy[:])
		ulidLastMs = ms
	default:
		// 同毫秒 或 时钟回拨：保持时间戳单调，entropy +1
		if !ulidIncEntropy() {
			// 80bit 随机数溢出（几乎不可能），自旋到下一毫秒
			for ms <= ulidLastMs {
				runtime.Gosched()
				ms = Now()
			}
			_, _ = crand.Read(ulidEntropy[:])
			ulidLastMs = ms
		}
	}
	ts := ulidLastMs
	ent := ulidEntropy
	ulidMu.Unlock()

	var raw [16]byte
	raw[0] = byte(ts >> 40)
	raw[1] = byte(ts >> 32)
	raw[2] = byte(ts >> 24)
	raw[3] = byte(ts >> 16)
	raw[4] = byte(ts >> 8)
	raw[5] = byte(ts)
	copy(raw[ulidTimeLen:], ent[:])

	return ulidEncode(raw)
}

// ULIDFrom 由固定字符串派生稳定的 26 位 ID
//
// 同样的 key 永远得到完全相同的字符串（进程内、跨进程、跨机器、跨库一致），
// 不同 key 落在 128bit 哈希空间内，实际不会碰撞。
// 输出格式与 ULID 完全一致（26 位 Crockford Base32，IsULID 为 true），
// 但不含时间信息，因此字典序没有时间含义，不能用于排序。
//
// 适用于「主键应当由业务键推导」的场景，例如权限点主键由权限标识派生，
// 这样重建库、换环境、回滚版本后主键都不会漂移，外部引用（如角色授权）依然有效。
func ULIDFrom(key string) string {
	sum := sha256.Sum256([]byte(key))
	var raw [16]byte
	copy(raw[:], sum[:16])
	return ulidEncode(raw)
}

// ulidEncode 把 128bit 大端数据编码为 26 位 Crockford Base32
// 130bit 容器中最高 2bit 恒为 0，故首字符一定 <= '7'
// 固定使用小写字符集，保证输出大小写形式唯一
func ulidEncode(raw [16]byte) string {
	hi := binary.BigEndian.Uint64(raw[0:8])  //高 64bit
	lo := binary.BigEndian.Uint64(raw[8:16]) //低 64bit

	var buf [ulidEncLen]byte
	// 128bit 编码为 26 个 5bit 组(130bit)，最高组只占 3bit（首字符 0-7）
	for i := 0; i < ulidEncLen; i++ {
		s := uint(5 * (ulidEncLen - 1 - i))
		var v uint64
		if s >= 64 {
			v = hi >> (s - 64)
		} else {
			v = hi<<(64-s) | lo>>s
		}
		buf[i] = ulidChars[v&0x1f]
	}
	return string(buf[:])
}

// ulidIncEntropy 把 entropy 当作 80bit 大端无符号整数 +1
// 返回 false 表示溢出
// 调用方需持有 ulidMu
func ulidIncEntropy() bool {
	for i := ulidRandLen - 1; i >= 0; i-- {
		ulidEntropy[i]++
		if ulidEntropy[i] != 0 {
			return true
		}
	}
	return false
}

// IsULID 判断是否为合法的 ULID 字符串
// 长度 26、字符集合法、首字符 <= '7'（128bit 编码到 130bit 的最高 2bit 恒为 0）
// Crockford Base32 解码本就不区分大小写，故首字母大小写都可接受
func IsULID(s string) bool {
	if len(s) != ulidEncLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		if ulidCharIndex(s[i]) < 0 {
			return false
		}
	}
	return ulidCharIndex(s[0]) <= 7
}

// ulidCharIndex 返回字符在 Crockford Base32 中的下标，非法返回 -1
// 大写字母先归一化为小写（Crockford 解码大小写不敏感）
func ulidCharIndex(c byte) int {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	for i := 0; i < len(ulidChars); i++ {
		if ulidChars[i] == c {
			return i
		}
	}
	return -1
}
