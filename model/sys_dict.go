package model

// Dict 配置
type Dict struct {
	Key   string `xorm:"PK UNIQUE VARCHAR(64) comment('唯一Key')" json:"key"` //key
	Value string `xorm:"VARCHAR(255) comment('值')" json:"value"`            //值
	Inner bool   `xorm:"TINYINT(4) comment('内部禁止删除')" json:"inner"`         //内部
	Intro string `xorm:"VARCHAR(255) comment('说明')" json:"intro"`           //说明
}

func (Dict) TableName() string {
	return "sys_dict"
}

// DictGet 单条字典
func DictGet(key string) (*Dict, bool) {
	mod := &Dict{}
	has, _ := db.ID(key).Get(mod)
	return mod, has
}

// DictPage 字典分页
func DictPage(pi int, ps int, cols ...string) ([]Dict, error) {
	mods := make([]Dict, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// DictCount 字典分页总数
func DictCount() int {
	mod := &Dict{}
	sess := db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// DictAdd 添加字典
func DictAdd(mod *Dict) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// DictEdit 编辑字典
func DictEdit(mod *Dict, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Key).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// DictDrop 删除单条字典
func DictDrop(key string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(key).Delete(&Dict{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
