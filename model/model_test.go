package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestTag(t *testing.T) {
	defer os.RemoveAll("logs")
	//name: [{ required: true, message: "请填写", trigger: "blur", max: 128 }],
	generate(Menu{})

}
func generate(ptr interface{}) {
	t := reflect.TypeOf(ptr)
	// if t.Kind() == reflect.Ptr {
	// 	t = t.Elem()
	// }
	if t.Kind() != reflect.Struct {
		panic("Check type error not Struct")
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		tags := t.Field(i).Tag
		json := tags.Get("json")
		if json != "" {
			fmt.Printf("%s\n", json)
		}
	}
}
