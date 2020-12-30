package hmt

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// HmtAuth 凭证载体
type HmtAuth struct {
	Id    int    `json:"id"`
	Num   string `json:"num"`
	Name  string `json:"name"`
	Role  int    `json:"role"`
	ExpAt int64  `json:"exp"`
}

// Encode 生成 token
func (jc *HmtAuth) Encode(key string) string {
	data, _ := json.Marshal(jc)
	bStr := base64.RawURLEncoding.EncodeToString(data)
	hm := hmac.New(md5.New, []byte(key))
	hm.Write([]byte(bStr))
	sign := base64.RawURLEncoding.EncodeToString(hm.Sum(nil))
	return bStr + "." + sign
}

// Decode 生成 token
func (jc *HmtAuth) Decode(raw string, key string) error {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return errors.New("非法的token: " + raw)
	}
	hm := hmac.New(md5.New, []byte(key))
	hm.Write([]byte(parts[0]))
	sign := base64.RawURLEncoding.EncodeToString(hm.Sum(nil))
	if sign != parts[1] {
		return errors.New("token 非法")
	}
	datab, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return errors.New("base64解码失败,error: " + err.Error())
	}
	err = json.Unmarshal(datab, jc)
	if err != nil {
		return errors.New("json解码失败,error: " + err.Error())
	}
	if time.Now().Unix() > jc.ExpAt {
		return errors.New("toekn 已经过期")
	}
	return nil
}
