package token

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Token[T Share | Auth | string] struct {
	ExpAt int64 `json:"e"` //s
	Claim T     `json:"c"`
}

type Share struct {
	PostId string `json:"p"`
}

// Auth 凭证载体
// 只携带身份信息，权限判定统一由服务端按角色查询 grant.Routes 完成
type Auth struct {
	UserId string `json:"u"` //用户id
	RoleId string `json:"r"` //角色id
}

// Encode 生成 token
func (in Token[T]) Encode(key string) string {
	data, _ := json.Marshal(in)
	bStr := base64.RawURLEncoding.EncodeToString(data)
	hm := hmac.New(md5.New, []byte(key))
	hm.Write([]byte(bStr))
	sign := base64.RawURLEncoding.EncodeToString(hm.Sum(nil))
	return bStr + "." + sign
}

// Decode 验证
func (in *Token[T]) Decode(raw string, key string) error {
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
	err = json.Unmarshal(datab, in)
	if err != nil {
		return errors.New("json解码失败,error: " + err.Error())
	}
	//ExpAt <=0 不过期
	if in.ExpAt > 0 && time.Now().Unix() > in.ExpAt {
		return errors.New("toekn 已经过期")
	}
	return nil
}
