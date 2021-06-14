package token

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	auth := Auth{
		Id:     1,
		RoleId: 1000,
		ExpAt:  time.Now().Add(time.Hour * 24).Unix(),
	}
	t.Log(auth.Encode("key"))
}
func TestVerify(t *testing.T) {
	raw := "eyJpZCI6MSwicmlkIjoxMDAwLCJlYXQiOjE2MjA5Nzc3NTl9.8H_3zTm7B6RJFbXtQ5CVZxPkoag"
	Auth, err := Verify(raw, "key")
	t.Log(Auth, err)
}
func BenchmarkEncode(b *testing.B) {
	auth := Auth{
		Id:     1,
		RoleId: 1000,
		ExpAt:  time.Now().Add(time.Hour * 2).Unix(),
	}
	for i := 0; i < b.N; i++ {
		auth.Encode("key")
	}
}
func BenchmarkVerify(b *testing.B) {
	raw := "eyJpZCI6MSwicmlkIjoxMDAwLCJlYXQiOjE2MjA5Nzc3NTl9.8H_3zTm7B6RJFbXtQ5CVZxPkoag"
	for i := 0; i < b.N; i++ {
		Verify(raw, "key")
	}

}
