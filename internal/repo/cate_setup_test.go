package repo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"blog/conf"
	"blog/internal/model"

	"github.com/zxysilent/logs"
)

// init 在任何测试（以及它们启动的 goroutine）之前导走日志输出
// 放到 setup 里会与 repo 自身的日志协程并发读写 logs 的全局配置，触发数据竞争
func init() {
	logs.SetOutput(io.Discard)
}

// setupCateDB 初始化一套隔离的内存库环境
//   - 每个测试用独立命名的内存库，测试之间互不影响
//   - 走真实的 conf.Init + Init 流程，因此表结构、派生表重置、
//     幂等种子、栏目初始化都会被覆盖到
//   - 结束后关闭引擎，内存库随之销毁，不留任何落盘文件
func setupCateDB(t *testing.T) context.Context {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_", "\\", "_").Replace(t.Name())
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
	path := filepath.Join(t.TempDir(), "conf.toml")
	body := fmt.Sprintf(`title = "test"
mode = "test"
port = 0
host = "127.0.0.1"
token_exp = 1
token_secret = "test-secret"

[db]
type = "sqlite"
sqlite = %q
`, dsn)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write conf: %v", err)
	}
	conf.Init(context.Background(), path)
	Init(context.Background())
	engine := db
	t.Cleanup(func() { _ = engine.Close() })
	return context.Background()
}

// mustCateAdd 新建栏目并断言成功
func mustCateAdd(t *testing.T, ctx context.Context, mod model.Cate) *model.Cate {
	t.Helper()
	mod.BeforeInsert()
	if err := CateAdd(ctx, &mod); err != nil {
		t.Fatalf("CateAdd(%s): %v", mod.Name, err)
	}
	got, ok := CateCached(mod.Id)
	if !ok {
		t.Fatalf("CateAdd(%s): 未写入缓存", mod.Name)
	}
	copy := *got
	return &copy
}
