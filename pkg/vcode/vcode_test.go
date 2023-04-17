package vcode

import (
	"testing"
)

func TestVcode(t *testing.T) {
	v := New(4, "xx")
	t.Logf("%+v\n", v)
	t.Log(Check("xx", v.Real, "xx"))
}
func TestCheck(t *testing.T) {
	v := Check("xx", "NSCrCt2dCpuyRSO9bqKUfg.1681221854.jXfVwm6J3vbSlTPCOPL51g", "7403")
	t.Logf("%+v\n", v)
}
