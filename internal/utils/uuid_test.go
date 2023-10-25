package utils

import (
	"testing"
)

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}
func TestUUID(t *testing.T) {
	t.Log(UUID())
}
func BenchmarkSUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SUID()
	}
}
func TestSUID(t *testing.T) {
	t.Log(SUID())
}

func BenchmarkMUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MUID()
	}
}
func TestMUID(t *testing.T) {
	t.Log(MUID())
}
