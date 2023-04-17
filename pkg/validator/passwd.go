package validator

import (
	"errors"
	"strings"
)

func Passwd(passwd string) error {
	if len(passwd) < 6 {
		return errors.New("密码至少6位")
	}
	passwd = strings.ToLower(passwd)
	hasNum := false
	hasLetter := false
	for i := 0; i < len(passwd); i++ {
		if passwd[i] <= '9' && passwd[i] >= '0' {
			hasNum = true
		}
		if passwd[i] <= 'z' && passwd[i] >= 'a' {
			hasLetter = true
		}
	}
	if !hasNum {
		return errors.New("密码需要包含数字")
	}
	if !hasLetter {
		return errors.New("密码需要包含字母")
	}
	return nil
}
