package repo

import (
	"blog/internal/model"
	"context"
)

// DictExist 字典主键是否已存在
func DictExist(ctx context.Context, id string) bool {
	if id == "" {
		return false
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.ID(id).NoAutoCondition().Exist(&model.Dict{})
	return exist
}

// DictGet 按主键查询字典（含选项树）
func DictGet(ctx context.Context, id string) (*model.Dict, error) {
	mod, err := GetAny[model.Dict](ctx, id)
	if err != nil {
		return nil, err
	}
	mods := []model.Dict{mod}
	dictEntryFill(ctx, mods)
	return &mods[0], nil
}

// DictList 字典分页列表（含选项树）
func DictList(ctx context.Context, args *model.DictListArgs) ([]model.Dict, int, error) {
	mods := make([]model.Dict, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if args.Mult != "" {
		sess.And("(id like ? or name like ? or remark like ?)", args.SafeMult(), args.SafeMult(), args.SafeMult())
	}
	if args.Status != nil {
		sess.And("status = ?", *args.Status)
	}
	sess.Asc("id")
	var count int64
	var err error
	if args.IsZero() {
		err = sess.Find(&mods)
		count = int64(len(mods))
	} else {
		count, err = sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	}
	if err != nil {
		return mods, 0, err
	}
	dictEntryFill(ctx, mods)
	return mods, int(count), nil
}

// DictAll 全部字典（含选项树）
func DictAll(ctx context.Context) ([]model.Dict, error) {
	mods := make([]model.Dict, 0, 16)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Asc("id").Find(&mods); err != nil {
		return nil, err
	}
	dictEntryFill(ctx, mods)
	return mods, nil
}

// DictAdd 新增字典及其选项（单事务）
func DictAdd(ctx context.Context, mod *model.Dict, rows []model.DictEntry) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	if len(rows) > 0 {
		for idx := range rows {
			rows[idx].DictId = mod.Id
		}
		if _, err := sess.Insert(&rows); err != nil {
			sess.Rollback()
			return err
		}
	}
	return sess.Commit()
}

// DictEdit 修改字典，rows 非 nil 时全量替换选项（单事务）
func DictEdit(ctx context.Context, mod *model.Dict, rows []model.DictEntry, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	if rows != nil {
		if _, err := sess.Where("dict_id = ?", mod.Id).Delete(&model.DictEntry{}); err != nil {
			sess.Rollback()
			return err
		}
		if len(rows) > 0 {
			if _, err := sess.Insert(&rows); err != nil {
				sess.Rollback()
				return err
			}
		}
	}
	return sess.Commit()
}

// DictDrop 删除字典并级联删除其选项（单事务）
func DictDrop(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Where("dict_id = ?", id).Delete(&model.DictEntry{}); err != nil {
		sess.Rollback()
		return err
	}
	if _, err := sess.ID(id).NoAutoCondition().Delete(&model.Dict{}); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// DictEntryList 查询指定字典的选项（扁平，按排序）
func DictEntryList(ctx context.Context, dictId string) ([]model.DictEntry, error) {
	rows := make([]model.DictEntry, 0, 16)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.Where("dict_id = ?", dictId).Asc("sort").Asc("field").Find(&rows)
	return rows, err
}

// dictEntryFill 批量填充选项树，避免 N+1
func dictEntryFill(ctx context.Context, mods []model.Dict) {
	if len(mods) < 1 {
		return
	}
	ids := make([]string, 0, len(mods))
	for idx := range mods {
		ids = append(ids, mods[idx].Id)
	}
	rows := make([]model.DictEntry, 0, len(mods)*8)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.In("dict_id", ids).Asc("sort").Asc("field").Find(&rows); err != nil {
		for idx := range mods {
			mods[idx].Entries = []model.DictEntry{}
		}
		return
	}
	group := make(map[string][]model.DictEntry, len(mods))
	for _, r := range rows {
		group[r.DictId] = append(group[r.DictId], r)
	}
	for idx := range mods {
		mods[idx].BuildTree(group[mods[idx].Id])
	}
}
