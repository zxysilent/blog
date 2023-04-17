package utils

const (
	// StdDateTime 标准日期格式化
	StdDateTime = "2006-01-02 15:04:05"
	// StdDate 年月日
	StdDate = "20060102"
	// StdTime 时分秒
	StdTime = "15:04:05"
	// 更加编译平台判断int大小
	intSize = 32 << (^uint(0) >> 63) // 32 or 64
	// 最大Uint
	MaxUint = 1<<intSize - 1
	// 最大Int
	MaxInt = 1<<(intSize-1) - 1
	// 最小Int
	MinInt = -1 << (intSize - 1)
)
