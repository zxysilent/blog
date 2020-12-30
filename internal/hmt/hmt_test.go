package hmt

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	auth := HmtAuth{
		Id:    1,
		Role:  1000,
		Num:   "账号",
		Name:  "曾祥银",
		ExpAt: time.Now().Add(time.Hour * 2).Unix(),
	}
	t.Log(auth.Encode("key"))
}
func TestVerify(t *testing.T) {
	raw := "eyJpZCI6MSwibnVtIjoi6LSm5Y-3IiwibmFtZSI6IuabvuelpemTtiIsInJvbGUiOjEwMDAsImV4cCI6MTU5NzczMDgzM30.RgHOln4bWCX-O1OFpjImuvG9aoE"
	auth := HmtAuth{}
	err := auth.Decode(raw, "key")
	t.Log(auth, err)
}
func BenchmarkEncode(b *testing.B) {
	auth := HmtAuth{
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
	raw := "eyJpZCI6MSwibnVtIjoi6LSm5Y-3IiwibmFtZSI6IuabvuelpemTtiIsInJvbGUiOjEwMDAsImV4cCI6MTU5NzczMDgzM30.RgHOln4bWCX-O1OFpjImuvG9aoE"
	auth := HmtAuth{}
	for i := 0; i < b.N; i++ {
		auth.Decode(raw, "key")
	}
}
