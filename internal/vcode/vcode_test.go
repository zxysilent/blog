package vcode

import (
	"fmt"
	"os"
	"testing"

	"github.com/zxysilent/utils"
)

func TestVcode(t *testing.T) {
	s := "12345"
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}
	fout, _ := os.Create("out.png")
	NewImage(utils.RandDigitStr(6)).WriteTo(fout)
	fmt.Println(NewImage(utils.RandDigitStr(6)).Bs64())
	fout.Close()
}
