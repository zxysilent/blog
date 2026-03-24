package rndm

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}
func BenchmarkGUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GUID()
	}
}
func BenchmarkGUID1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(UUID(), "-", "")
	}
}
func TestUUID(t *testing.T) {
	t.Log(UUID())
}
func TestGUID(t *testing.T) {
	t.Log(GUID())
}

func TestMUID(t *testing.T) {
	t.Log(MUID())
}

func BenchmarkSUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SUID()
	}
}
func TestSUID(t *testing.T) {
	t.Log(SUID())
}

func TestBits(t *testing.T) {
	fmt.Printf("%b\n", 253402246800000/50)
}
func TestIsUUID(t *testing.T) {
	t.Log(IsUUID("123e4567-e89b-12d3-a456-426655440000"))
}
