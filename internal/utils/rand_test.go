package utils

import (
	"testing"
)

// 性能测试
func BenchmarkStr8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Str(8)
	}
}

func BenchmarkStr10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Str(10)
	}
}
func BenchmarkStr16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Str(16)
	}
}
func BenchmarkRandDigitStr8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Digit(8)
	}
}

func BenchmarkRandDigitStr16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Digit(16)
	}
}

func TestStr(t *testing.T) {
	Str(16)
}
