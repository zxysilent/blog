package repo

import (
	"blog/internal/model"
	"context"
	"sync"
)

// DictList 字典分页
func DictList(ctx context.Context, filter *model.DictListArgs, cols ...string) ([]model.Dict, int, error) {
	mods := make([]model.Dict, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if filter.Name != "" {
		sess.And("name like concat('%',?,'%')", fixlike(filter.Name))
	}
	sess.Asc("sort").Desc("created")
	if filter.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	sess.Limit(filter.GetLimit(), filter.GetOffset())
	count, err := sess.FindAndCount(&mods)
	return mods, int(count), err
}

// DictEdit 编辑字典
func DictEdit(ctx context.Context, mod *model.Dict, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Omit("key")
	if _, err := sess.ID(mod.Key).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	SyncCfg(ctx)
	return nil
}

var _dict = _MemDict{mu: &sync.RWMutex{}, cs: model.NewDictSet()}

func GetCfg() *model.DictSet {
	_dict.mu.RLock()
	defer _dict.mu.RUnlock()
	return _dict.cs
}

func GetCfgx(key string) *model.Dict {
	_dict.mu.RLock()
	defer _dict.mu.RUnlock()
	return _dict.cs.Get(key)
}

func SyncCfg(ctx context.Context) {
	_dict.mu.Lock()
	defer _dict.mu.Unlock()
	mods, _, _ := DictList(ctx, &model.DictListArgs{})
	for _, mod := range mods {
		switch mod.Key {
		case model.DictKeyBlog:
			_dict.cs.Blog = mod.ToBlog()
		case model.DictKeyBasic:
			_dict.cs.Basic = mod.ToBasic()
		default:
			_dict.cs.Set(mod.Key, &mod)
		}
	}
}

type _MemDict struct {
	cs *model.DictSet
	mu *sync.RWMutex
}
