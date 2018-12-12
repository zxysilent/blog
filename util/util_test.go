package util

import (
	"testing"
)

// 性能测试
func BenchmarkRandStr6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStr(16)
	}
}
