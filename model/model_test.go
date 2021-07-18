package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGen(t *testing.T) {
	defer os.RemoveAll("logs")
	//name: [{ required: true, message: "请填写", trigger: "blur", max: 128 }],
	generate(User{})

}
func generate(ipt interface{}) {
	// 包装代码-方便下面调整
	opt := check(ipt)
	fieldNum := opt.NumField()
	for i := 0; i < fieldNum; i++ {
		tags := opt.Field(i).Tag
		json := tags.Get("json")
		if json != "" {
			fmt.Printf("%s\n", json)
		}
	}
}

// check确保可以获取结构体属性
func check(ipt interface{}) reflect.Type {
	opt := reflect.TypeOf(ipt)
	if opt.Kind() == reflect.Ptr {
		opt = opt.Elem()
	}
	if opt.Kind() != reflect.Struct {
		panic("type error, not Struct ：" + opt.Kind().String())
	}
	return opt
}
