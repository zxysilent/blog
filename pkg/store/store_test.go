package store

import (
	"fmt"
	"log"
	"testing"
)

// go clean -testcache

func TestBlot(t *testing.T) {
	store, err := New("test")
	log.Println(err)
	store.Set("k", "")
	s := S{
		Name: "111",
		Size: 1111222,
	}
	store.Setx("json", s)
	log.Println(store.Get("k"))
	log.Println(store.Get("json1"))

	log.Println(store.With("111").Setx("k", s))
	var s1 S
	log.Println(store.With("111").Getx("k", &s1), s1)
	fmt.Println(s1)
}

type S struct {
	Name string
	Size int64
}
