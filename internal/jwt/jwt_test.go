package jwt

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	auth := JwtAuth{
		Id:    1,
		Role:  1000,
		Num:   "账号",
		Name:  "曾祥银",
		ExpAt: time.Now().Add(time.Hour * 2).Unix(),
	}
	t.Log(auth.Encode("key"))
}
func TestVerify(t *testing.T) {
	raw := "eyJpZCI6MSwibnVtIjoi6LSm5Y-3IiwibmFtZSI6IuabvuelpemTtiIsInJvbGUiOjEwMDAsImV4cCI6MTU4MzI0NjkxMH0=.3hVx7outOhWuz-P_QygaeYvCKD4="
	jwtAuth, err := Verify(raw, "key")
	t.Log(jwtAuth, err)
}
func BenchmarkEncode(b *testing.B) {
	auth := JwtAuth{
		Id:    1,
		Role:  1000,
		Num:   "账号",
		Name:  "曾祥银",
		ExpAt: time.Now().Add(time.Hour * 2).Unix(),
	}
	for i := 0; i < b.N; i++ {
		auth.Encode("key")
	}
}
func BenchmarkVerify(b *testing.B) {
	raw := "eyJpZCI6MSwibnVtIjoi6LSm5Y-3IiwibmFtZSI6IuabvuelpemTtiIsInJvbGUiOjEwMDAsImV4cCI6MTU4MzI0NjkxMH0=.3hVx7outOhWuz-P_QygaeYvCKD4="
	for i := 0; i < b.N; i++ {
		Verify(raw, "key")
	}

}
