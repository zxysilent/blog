package repo

import (
	"blog/internal/model"
	"time"
)

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tagId string) int {
	var count int
	db.SQL(`SELECT count(post_id) as count FROM post_tag WHERE tag_id=?`, tagId).Get(&count)
	return count
}

// // TagPostPage 通过标签查询文章分页
// func TagPostPage(filter *model.PostTagFilter) ([]model.PostTag, error) {
// 	mods := make([]model.PostTag, 0, filter.Ps)
// 	err := db.SQL(`SELECT id,post_id,tag_id FROM post_tag WHERE tag_id=? ORDER BY post_id DESC LIMIT ?,? `, filter.TagId, (filter.Pi-1)*filter.Ps, filter.Ps).Find(&mods)
// 	if len(mods) > 0 {
// 		ids := make([]string, 0, len(mods))
// 		for idx := range mods {
// 			if !slices.Contains(ids, mods[idx].PostId) {
// 				ids = append(ids, mods[idx].PostId)
// 			}
// 		}
// 		mapSet := PostMapIds(ids)
// 		for idx := range mods {
// 			mods[idx].Post = mapSet[mods[idx].PostId]
// 		}
// 	}
// 	return mods, err
// }

// TagPostDrop 删除标签时删除对应的标签_文章
func TagPostDrop(tagId string) error {
	sess := db.NewSession()
	defer sess.Close()
	if _, err := sess.Where("tag_id = ?", tagId).Delete(&model.PostTag{}); err != nil {
		return err
	}
	return nil
}

// PostTags 查询文章关联的标签
func PostTags(postId string) []model.Tag {
	mods := make([]model.Tag, 0, 8)
	db.SQL("SELECT tag.* FROM post_tag ,tag WHERE tag.id=tag_id AND post_id=?", postId).Find(&mods)
	return mods
}

// PostTagSync 全量替换文章标签
func PostTagSync(postId string, tagIds []string) error {
	sess := db.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Where("post_id = ?", postId).Delete(&model.PostTag{}); err != nil {
		sess.Rollback()
		return err
	}
	if len(tagIds) > 0 {
		milli := time.Now().UnixMilli()
		mods := make([]model.PostTag, 0, len(tagIds))
		for _, tid := range tagIds {
			mods = append(mods, model.PostTag{PostId: postId, TagId: tid, Updated: milli, Created: milli})
		}
		for idx := range mods {
			mods[idx].BeforeInsert()
		}
		if _, err := sess.Insert(&mods); err != nil {
			sess.Rollback()
			return err
		}
	}
	return sess.Commit()
}
