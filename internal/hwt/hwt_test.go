package hwt

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	auth := Auth{
		Id:     1,
		RoleId: 1,
		Num:    "账号",
		ExpAt:  time.Now().Add(time.Hour * 24).Unix(),
	}
	t.Log(auth.Encode("key"))
}
func TestVerify(t *testing.T) {
	raw := "eyJpIjoxLCJuIjoi6LSm5Y-3IiwicmkiOjEsImUiOjE2MTkwNTIxNDF9.BgqifIFamVp_A_IEQ9lSIw"
	auth := Auth{}
	err := auth.Decode(raw, "key")
	t.Log(auth, err)
}
func BenchmarkEncode(b *testing.B) {
	auth := Auth{
		Id:     1,
		Num:    "账号",
		RoleId: 1,
		ExpAt:  time.Now().Add(time.Hour * 2).Unix(),
	}
	for i := 0; i < b.N; i++ {
		auth.Encode("key")
	}
}
func BenchmarkVerify(b *testing.B) {
	raw := "eyJpIjoxLCJuIjoi6LSm5Y-3IiwicmkiOjEsImUiOjE2MTkwNTIxNDF9.BgqifIFamVp_A_IEQ9lSIw"
	auth := Auth{}
	for i := 0; i < b.N; i++ {
		auth.Decode(raw, "key")
	}
}
