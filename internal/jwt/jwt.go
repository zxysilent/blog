package jwt

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// JwtAuth 凭证载体
type JwtAuth struct {
	Id    int    `json:"id"`
	Num   string `json:"num"`
	Name  string `json:"name"`
	Role  int    `json:"role"`
	ExpAt int64  `json:"exp"`
}

// Encode 生成 token
func (jc *JwtAuth) Encode(key string) string {
	data, _ := json.Marshal(jc)
	bStr := base64.URLEncoding.EncodeToString(data)
	h := sha1.New()
	h.Write([]byte(bStr + key))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return bStr + "." + sign
}

// Verify 验证
func Verify(raw string, key string) (*JwtAuth, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return nil, errors.New("非法的token: " + raw)
	}
	h := sha1.New()
	h.Write([]byte(parts[0] + key))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if sign != parts[1] {
		return nil, errors.New("token 非法")
	}
	datab, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("base64解码失败,error: " + err.Error())
	}
	ja := &JwtAuth{}
	err = json.Unmarshal(datab, ja)
	if err != nil {
		return nil, errors.New("json解码失败,error: " + err.Error())
	}
	if time.Now().Unix() > ja.ExpAt {
		return nil, errors.New("toekn 已经过期")
	}
	return ja, nil
}
