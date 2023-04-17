package token

import (
	"testing"
	"time"
)

func TestBaseEncode(t *testing.T) {
	mod := Token[string]{}
	mod.ExpAt = time.Now().Add(time.Hour * 24).Unix()
	t.Log(mod.Encode("key"))
}
func TestBaseVerify(t *testing.T) {
	raw := "eyJlYXQiOjE2ODE1NzE2NjN9.h7fVRA2qLw2336WjK-Ou4kDMYnE"
	mod := Token[string]{}
	err := mod.Decode(raw, "key")
	t.Log(mod, err)
}
func TestAuthEncode(t *testing.T) {
	mod := &Token[Auth]{
		Claim: Auth{
			UserId: 1,
			RoleId: 1000,
		},
	}
	mod.ExpAt = time.Now().Add(time.Hour * 24).Unix()
	t.Log(mod.Encode("key"))
}
func TestAuthVerify(t *testing.T) {
	raw := "eyJlYXQiOjE2ODE1NzI3OTV9.IhCdCZO1npx46cGM4ao9I5JzigY"
	mod := &Token[Auth]{}
	err := mod.Decode(raw, "key")
	t.Log(mod, err)
}
func BenchmarkEncode(b *testing.B) {
	mod := &Token[Auth]{
		Claim: Auth{
			UserId: 1,
			RoleId: 1000,
		},
	}
	mod.ExpAt = time.Now().Add(time.Hour * 24).Unix()
	for i := 0; i < b.N; i++ {
		mod.Encode("key")
	}
}
func BenchmarkVerify(b *testing.B) {
	raw := "eyJpZCI6MSwicmlkIjoxMDAwLCJlYXQiOjE2MjA5Nzc3NTl9.8H_3zTm7B6RJFbXtQ5CVZxPkoag"
	mod := &Token[Auth]{}
	for i := 0; i < b.N; i++ {
		mod.Decode(raw, "key")
	}
}
func TestShareEncode(t *testing.T) {
	mod := &Token[Share]{
		Claim: Share{PostId: 84},
	}
	mod.ExpAt = time.Now().Add(time.Hour * 24).Unix()
	t.Log(mod.Encode("key"))
}
func TestShareVerify(t *testing.T) {
	raw := "eyJlIjoxNjgxNjI0MTIyLCJjIjp7InAiOjg0fX0.u3W_s9rA8zfG5W0u5XZJKg"
	mod := &Token[Share]{}
	err := mod.Decode(raw, "key")
	t.Log(mod, err)
}
