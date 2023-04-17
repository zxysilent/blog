package model

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJson(t *testing.T) {
	defer os.RemoveAll("logs")
	mod := Blog{}
	b, _ := json.Marshal(mod)
	t.Log(string(b))
}
