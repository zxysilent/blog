package kernel

import (
	"blog/internal/repo"
	"blog/pkg/targz"

	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zxysilent/logs"
)

const backupPath = ".backs/"

func initCron() {
	// 标准格式不支持秒
	task := cron.New(cron.WithSeconds())
	// 每天三点备份
	// |0 0 3 1/1 * ?|
	// 每3天三点备份
	// |0 0 3 1/3 * ?|
	// 测试
	// |0/10 * * * * ?|每10秒
	// |0 0/1 * * * ?|每1分钟
	// |0 0/5 * * * ?|每5分钟
	task.AddFunc("0 0 3 1/1 * ?", func() {
		backupDb()
	})
	// 每月最后一天23:40
	task.AddFunc("0 40 23 L * ?", func() {
		backupStatic()
	})
	go task.Run()
}

func backupDb() (string, error) {
	now := time.Now()
	dir := now.Format("20060102_150405d")
	archive := backupPath + "/db/" + dir
	os.MkdirAll(archive, 0755)
	name := now.Format("/200601021504.sql")
	logs.Info("[back sql] ", repo.GetDb().DumpAllToFile(archive+name))
	// 保留最近 90 天
	// os.Remove(backupPath + "/db/" + now.AddDate(0, 0, -90).Format("20060102.tar.gz"))
	logs.Info("[targz] ", targz.Targz(archive+".tar.gz", archive))
	os.RemoveAll(archive)
	return archive + ".tar.gz", nil
}

func backupStatic() (string, error) {
	now := time.Now()
	name := now.Format("20060102_150405s.tar.gz")
	archive := backupPath + "static/"
	os.MkdirAll(archive, 0755)
	err := targz.Targz(archive+name, "static")
	if err != nil {
		return "", err
	}
	return archive + name, nil
}
