package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// TagExist 标签是否存在
func TagExist(ctx context.Context, mod *model.Tag) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// TagGet 标签单条数据
func TagGet(ctx context.Context, mod *model.Tag) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if has, _ := sess.Get(mod); !has {
		return ErrNotExist
	}
	return nil
}

// TagList 标签列表数据
func TagList(ctx context.Context, filter *model.TagListArgs, cols ...string) ([]model.Tag, int, error) {
	mods := make([]model.Tag, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	sess.Asc("created")
	if filter.Mult != "" {
		sess.And("(name like ? or intro like ?)", filter.SafeMult(), filter.SafeMult())
	}
	if filter.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	sess.Limit(filter.GetLimit(), filter.GetOffset())
	count, err := sess.FindAndCount(&mods)
	return mods, int(count), err
}

// TagCount 标签数据总数
func TagCount(ctx context.Context, filter *model.TagListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if filter.Mult != "" {
		sess.And("(name like ? or intro like ?)", filter.SafeMult(), filter.SafeMult())
	}
	count, _ := sess.Count(&model.Tag{})
	return int(count)
}

// TagAdd 标签添加数据
func TagAdd(ctx context.Context, mod *model.Tag) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagEdit 标签修改数据
func TagEdit(ctx context.Context, mod *model.Tag, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagDrop 标签删除数据
func TagDrop(ctx context.Context, mod *model.Tag) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if _, err := sess.Delete(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagIds 标签id集合查询
func TagIds(ctx context.Context, ids []string) []model.Tag {
	mods := make([]model.Tag, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// TagMapIds 标签id集合查询
func TagMapIds(ctx context.Context, ids []string) map[string]*model.Tag {
	mods := make([]model.Tag, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.Tag, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}

// TagStatus 所有标签统计 当前标签下有文章才显示
func TagStatus() ([]model.TagState, error) {
	mods := make([]model.TagState, 0, 8)
	err := db.SQL("SELECT `name`,intro,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Find(&mods)
	return mods, err
}
