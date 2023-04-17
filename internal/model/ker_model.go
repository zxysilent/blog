package model

// ------------------------------------------------------ 信息 ------------------------------------------------------
// StateGo go information
type StateGo struct {
	ARCH         string `json:"arch"`
	OS           string `json:"os"`
	Version      string `json:"version"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

// ------------------------------------------------------ 日志 ------------------------------------------------------
// Log 日志
type Log struct {
	Id      int    `xorm:"INT PK AUTOINCR" json:"id"`           //
	AdminId int    `xorm:"INT DEFAULT 0" json:"admin_id"`       //用户id
	Module  string `xorm:"VARCHAR(255)" json:"module"`          //模块
	Action  string `xorm:"VARCHAR(255)" json:"action"`          //做了什么
	Admin   *Admin `xorm:"-" swaggerignore:"true" json:"admin"` //附加
	Data    string `xorm:"TEXT" json:"data"`                    //数据
	Created int64  `xorm:"BIGINT" json:"created"`               //
}

func (*Log) TableName() string {
	return sysPrefix + "log"
}

const sysPrefix = ""
