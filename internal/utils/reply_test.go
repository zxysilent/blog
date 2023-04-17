package utils

import "testing"

func TestReply(t *testing.T) {
	Succ("succ")
	Fail("fail")
	Page("page", []int{}, 0)
	ErrDeny("deny")
	ErrIpt("ipt")
	ErrOpt("opt")
	ErrSvr("svr")
	ErrToken("token")
	Ext("ext")
}
func TestReplyData(t *testing.T) {
	Succ("succ", struct{ A, b int }{})
	Fail("fail", struct{ A, b int }{})
	Page("page", []int{}, 0)
	ErrDeny("deny", struct{ A, b int }{})
	ErrIpt("ipt", struct{ A, b int }{})
	ErrOpt("opt", struct{ A, b int }{})
	ErrSvr("svr", struct{ A, b int }{})
	ErrToken("token")
	Ext("ext", struct{ A, b int }{})
}
