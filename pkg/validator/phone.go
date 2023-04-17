package validator

import (
	"errors"
	"strings"
)

func Phone(phone string) error {
	if len(phone) != 11 {
		return errors.New("请输入11位电话号码")
	}

	phone = strings.ToLower(phone)
	if phone[0] != '1' {
		return errors.New("非法的电话号码")
	}
	for i := 1; i < len(phone); i++ {
		if phone[i] > '9' || phone[i] < '0' {
			return errors.New("非法的电话号码")
		}
	}
	return nil
}
