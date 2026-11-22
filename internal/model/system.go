package model

// StateGo 服务器运行时信息
// 由 runtime 包实时取值，不落库
type StateGo struct {
	ARCH         string `json:"arch"`          //架构
	OS           string `json:"os"`            //操作系统
	Version      string `json:"version"`       //Go 版本
	NumCPU       int    `json:"num_cpu"`       //CPU 核数
	NumGoroutine int    `json:"num_goroutine"` //协程数
}
