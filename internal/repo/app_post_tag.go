package repo

import "blog/internal/model"

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tagId int) int {
	var count int
	db.SQL(`SELECT count(post_id) as count FROM post_tag WHERE tag_id=?`, tagId).Get(&count)
	return count
}

// TagPostPage 通过标签查询文章分页
func TagPostPage(filter *model.PostTagFilter) ([]model.PostTag, error) {
	mods := make([]model.PostTag, 0, filter.Ps)
	err := db.SQL(`SELECT id,post_id,tag_id FROM post_tag WHERE tag_id=? ORDER BY post_id DESC LIMIT ?,? `, filter.TagId, (filter.Pi-1)*filter.Ps, filter.Ps).Find(&mods)
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].PostId, ids) {
				ids = append(ids, mods[idx].PostId)
			}
		}
		mapSet := PostMapIds(ids)
		for idx := range mods {
			mods[idx].Post = mapSet[mods[idx].PostId]
		}
	}
	return mods, err
}

// TagPostAdds 添加文章标签[]
func TagPostAdds(mods *[]model.PostTag) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.Insert(mods); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagPostDrop 删除标签时候
// 删除对应的标签_文章
func TagPostDrop(tagId int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("tag_id = ?", tagId)
	if _, err := sess.Delete(&model.PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}

// PostTagDrops 修改文章时候
// 删除对应的标签_文章
func PostTagDrops(postId int, tagIds []int) error {
	if len(tagIds) < 1 {
		return nil
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("post_id = ?", postId).In("tag_id", tagIds)
	if _, err := sess.Delete(&model.PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}

// PostTagDrop 删除文章时候
// 删除对应的标签_文章
func PostTagDrop(postId int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("post_id = ?", postId)
	if _, err := sess.Delete(&model.PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}
