package validator

import (
	"errors"
	"strconv"
	"strings"
)

func Idcard(idcard string) error {
	if len(idcard) != 18 {
		return errors.New("请输入18位身份证号")
	}
	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	valid := []byte{'1', '0', 'x', '9', '8', '7', '6', '5', '4', '3', '2'}
	idcard = strings.ToLower(idcard)
	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(idcard[i : i+1])
		// log.Println(num)
		sum += num * weight[i]
	}
	if valid[sum%11] != idcard[17] {
		return errors.New("非法的身份证号")
	}
	return nil
}
