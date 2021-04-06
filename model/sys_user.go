package model

import "time"

// User 用户
type User struct {
	Id     int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	Name   string    `xorm:"VARCHAR(255) comment('姓名')" json:"name"`
	Num    string    `xorm:"VARCHAR(255) UNIQUE comment('账号')" json:"num"`
	Passwd string    `xorm:"VARCHAR(255) comment('密码')" json:"passwd"`
	RoleId int       `xorm:"INT(11) DEFAULT 0 comment('角色')" json:"role_id"`
	Email  string    `xorm:"VARCHAR(255) UNIQUE comment('邮箱')" json:"email"`
	Phone  string    `xorm:"VARCHAR(255) UNIQUE comment('电话')" json:"phone"`
	Lock   bool      `xorm:"TINYINT(4) DEFAULT 0 comment('锁定')" json:"lock"`
	Ecount int       `xorm:"INT(11) DEFAULT 0 comment('错误次数')" json:"ecount"`
	Ltime  time.Time `xorm:"DATETIME comment('上次登录时间')" json:"ltime"`
	Ctime  time.Time `xorm:"DATETIME comment('创建时间')" json:"ctime"`
	Role   Role      `xorm:"-" json:"role"`  //用户角色
	Auths  []Api     `xorm:"-" json:"auths"` //权限集合
	Menus  []Menu    `xorm:"-" json:"menus"` //菜单导航
}

func (User) TableName() string {
	return "sys_user"
}

//UserByNum 通过账号获取用户信息
func UserByNum(num string) (*User, bool) {
	mod := &User{}
	has, _ := Db.Where("num=?", num).Get(mod)
	return mod, has
}

//UserExist 判断是否存在当前账号
func UserExist(num string) bool {
	has, _ := Db.Exist(&User{
		Num: num,
	})
	return has
}

// UserGet 单条用户信息
func UserGet(id int) (*User, bool) {
	mod := &User{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// UserAll 所有用户信息
func UserAll() ([]User, error) {
	mods := make([]User, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// UserPage 用户分页信息
func UserPage(pi int, ps int, cols ...string) ([]User, error) {
	mods := make([]User, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// UserCount 用户分页信息总数
func UserCount() int {
	mod := &User{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// UserAdd 添加用户信息
func UserAdd(mod *User) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// UserEdit 编辑用户信息
func UserEdit(mod *User, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// UserIds 返回用户信息-ids
func UserIds(ids []int) map[int]*User {
	mods := make([]User, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*User, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// UserDrop 删除单条用户信息
func UserDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&User{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
