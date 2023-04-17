package targz

import "testing"

func TestTargz(t *testing.T) {
	err := Targz("out.tar.gz", "targz.go")
	if err != nil {
		t.Error(err)
	}
}
func TestUargz(t *testing.T) {
	err := UnTargz("out/", "out.tar.gz")
	if err != nil {
		t.Error(err)
	}
}
