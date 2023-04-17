package repo

import (
	"blog/internal/model"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// ------------------------------------------------------ model.Basic 基础配置 ------------------------------------------------------

var memBasic struct {
	mem  *model.Basic
	lock sync.RWMutex
}

func CfgBasic() model.Basic {
	memBasic.lock.Lock()
	defer memBasic.lock.Unlock()
	cache := *memBasic.mem
	return cache
}

func initBasic() error {
	memBasic.lock.Lock()
	defer memBasic.lock.Unlock()
	mod, has := DictGet(model.DictKeyBasic)
	if !has {
		return errors.New("Not found")
	}
	memBasic.mem = mod.ToBasic()
	logs.Info("basic cache")
	return nil
}

var memBlog struct {
	mem  *model.Blog
	lock sync.RWMutex
}

func CfgBlog() model.Blog {
	memBlog.lock.Lock()
	defer memBlog.lock.Unlock()
	cache := *memBlog.mem
	return cache
}

func initBlog() error {
	memBlog.lock.Lock()
	defer memBlog.lock.Unlock()
	mod, has := DictGet(model.DictKeyBlog)
	if !has {
		return errors.New("Not found")
	}
	memBlog.mem = mod.ToBlog()
	logs.Info("blog cache")
	return nil
}

// ------------------------------------------------------ 日志 ------------------------------------------------------

// LogPage 记录分页
func LogPage(pi int, ps int, cols ...string) ([]model.Log, error) {
	mods := make([]model.Log, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Created").Limit(ps, (pi-1)*ps).Find(&mods)
	for idx := range mods {
		mods[idx].Admin, _ = AdminGet(mods[idx].AdminId)
	}
	return mods, err
}

// LogCount 记录分页总数
func LogCount() int {
	mod := &model.Log{}
	sess := db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// LogAdd 添加记录
func LogAdd(mod *model.Log) error {
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

func LOG(any interface{}, modeule string, action string, data ...interface{}) {
	adminId := 0
	switch ctx := any.(type) {
	case int:
		adminId = ctx
	case echo.Context:
		adminId, _ = ctx.Get("uid").(int)
	default:
	}
	mod := model.Log{
		AdminId: adminId,
		Module:  modeule,
		Action:  action,
	}
	if len(data) > 0 {
		switch ext := data[0].(type) {
		case int:
			mod.Data = strconv.Itoa(ext)
		case string:
			mod.Data = ext
		default:
			b, _ := json.Marshal(data[0])
			mod.Data = string(b)
		}
	}
	mod.Created = time.Now().Unix()
	LogAdd(&mod)
}
