package repo

import (
	"blog/internal/model"
	"context"
	"sync"
	"time"

	"github.com/zxysilent/logs"
	"xorm.io/xorm"
)

// 操作日志写入队列
//
// 日志是旁路数据：写入失败不能影响业务，写入耗时也不能拖慢请求。
// 因此中间件只做一次非阻塞入队，真正的落库在单独协程里完成。
// 队列满时丢弃并记一条警告——宁可漏记，也不能阻塞正常请求。
const logQueueSize = 256

var (
	logQueueOnce sync.Once
	logQueue     chan *model.Log
	// logPending 统计已入队但尚未落库的条目，供 LogFlush 等待队列排空
	logPending sync.WaitGroup
)

// logWorker 启动后台写日志协程，进程内只启动一次
func logWorker() {
	logQueueOnce.Do(func() {
		logQueue = make(chan *model.Log, logQueueSize)
		go func() {
			for mod := range logQueue {
				if err := logInsert(mod); err != nil {
					logs.With().Caller(false).Err(err).Str("perm", mod.Perm).Warn("repo.log.insert")
				}
				logPending.Done()
			}
		}()
	})
}

// LogPush 非阻塞入队
// 队列满时直接丢弃，保证调用方（中间件）永远不会被日志拖住
func LogPush(mod *model.Log) {
	if mod == nil {
		return
	}
	logWorker()
	logPending.Add(1)
	select {
	case logQueue <- mod:
	default:
		logPending.Done()
		logs.With().Caller(false).Str("perm", mod.Perm).Warn("repo.log.queue.full, dropped")
	}
}

// LogFlush 等待已入队的日志全部落库
// 仅用于测试与进程优雅退出；超时后直接返回，不阻断调用方
func LogFlush(timeout time.Duration) {
	logWorker()
	done := make(chan struct{})
	go func() {
		logPending.Wait()
		close(done)
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
	}
}

// LogAdd 同步写入一条日志
// 需要确定写入顺序时使用（如测试），业务路径请用 LogPush
func LogAdd(ctx context.Context, mod *model.Log) error {
	if mod == nil {
		return nil
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	_, err := sess.InsertOne(mod)
	return err
}

// logInsert 后台协程写库，不持有请求上下文
func logInsert(mod *model.Log) error {
	sess := db.NewSession()
	defer sess.Close()
	_, err := sess.InsertOne(mod)
	return err
}

// logCondition 组装日志查询条件
func logCondition(sess *xorm.Session, filter *model.LogFilter) {
	if filter.UserId != nil && *filter.UserId != "" {
		sess.And("user_id = ?", *filter.UserId)
	}
	if filter.Module != "" {
		sess.And("module = ?", filter.Module)
	}
	if filter.Perm != "" {
		sess.And("perm = ?", filter.Perm)
	}
	if filter.Method != "" {
		sess.And("method = ?", filter.Method)
	}
	if filter.Status != nil && *filter.Status > 0 {
		sess.And("status = ?", *filter.Status)
	}
	if filter.From > 0 {
		sess.And("created >= ?", filter.From)
	}
	if filter.To > 0 {
		sess.And("created <= ?", filter.To)
	}
}

// LogList 操作日志分页列表
func LogList(ctx context.Context, args *model.LogListArgs, cols ...string) ([]model.Log, int, error) {
	mods := make([]model.Log, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	logCondition(sess, &args.LogFilter)
	if args.Mult != "" {
		sess.And("(action like ? or path like ? or module like ?)",
			args.SafeMult(), args.SafeMult(), args.SafeMult())
	}
	sess.Desc("created").Desc("id")
	if args.IsZero() {
		err := sess.Find(&mods)
		logFillUser(ctx, mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	logFillUser(ctx, mods)
	return mods, int(count), err
}

// logFillUser 批量填充操作人，避免 N+1
func logFillUser(ctx context.Context, mods []model.Log) {
	if len(mods) < 1 {
		return
	}
	ids := make([]string, 0, len(mods))
	for idx := range mods {
		if mods[idx].UserId != "" {
			ids = append(ids, mods[idx].UserId)
		}
	}
	if len(ids) < 1 {
		return
	}
	ump := UserMapIds(ctx, ids)
	for idx := range mods {
		if u, ok := ump[mods[idx].UserId]; ok {
			mods[idx].User = u
		}
	}
}

// LogGet 操作日志详情
func LogGet(ctx context.Context, id string) (*model.Log, error) {
	mod, err := GetAny[model.Log](ctx, id)
	if err != nil {
		return nil, err
	}
	if mod.UserId != "" {
		if u, err := UserGet(ctx, mod.UserId); err == nil {
			mod.User = u
		}
	}
	return &mod, nil
}

// LogDrop 删除单条日志
// 日志表是唯一允许清理的历史数据，不参与任何业务逻辑
func LogDrop(ctx context.Context, id string) error {
	return DelAny[model.Log](ctx, id)
}

// LogClear 按条件清空日志，返回删除条数
// before 为 0 表示清空全部；否则只删早于该时刻的记录
func LogClear(ctx context.Context, before int64) (int64, error) {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if before > 0 {
		sess.Where("created < ?", before)
	} else {
		// xorm 删除必须带条件，这里用一个恒真条件
		sess.Where("1 = 1")
	}
	return sess.Delete(&model.Log{})
}
