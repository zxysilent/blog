package rndm

import (
	"sort"
	"strconv"
	"testing"
	"time"
)

// TestULID 基础格式校验
func TestULID(t *testing.T) {
	s := ULID()
	t.Log(s)
	if len(s) != ulidEncLen {
		t.Fatalf("长度应为 %d，实际 %d", ulidEncLen, len(s))
	}
	if !IsULID(s) {
		t.Fatalf("生成的 ULID 不合法: %s", s)
	}
	if s[0] > '7' {
		t.Fatalf("首字符必须 <= '7'，实际 %c", s[0])
	}
}

// TestULIDFrom 固定串派生：稳定、合法、不碰撞
func TestULIDFrom(t *testing.T) {
	a := ULIDFrom("permission:user:list")
	b := ULIDFrom("permission:user:list")
	t.Log(a)
	if a != b {
		t.Fatalf("同一 key 应得到相同结果: %s != %s", a, b)
	}
	if len(a) != ulidEncLen {
		t.Fatalf("长度应为 %d，实际 %d", ulidEncLen, len(a))
	}
	if !IsULID(a) {
		t.Fatalf("派生 ID 不合法: %s", a)
	}
	// 不同 key 不碰撞
	const n = 20000
	set := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		s := ULIDFrom("permission:key" + strconv.Itoa(i))
		if !IsULID(s) {
			t.Fatalf("派生 ID 不合法: %s", s)
		}
		if _, has := set[s]; has {
			t.Fatalf("派生 ID 碰撞: %s", s)
		}
		set[s] = struct{}{}
	}
}

// TestULIDLowercase 输出恒为小写
func TestULIDLowercase(t *testing.T) {
	check := func(s string) {
		for i := 0; i < len(s); i++ {
			if s[i] >= 'A' && s[i] <= 'Z' {
				t.Fatalf("ULID 应全小写，实际含大写: %s", s)
			}
		}
	}
	for i := 0; i < 1000; i++ {
		check(ULID())
		check(ULIDFrom("permission:key" + strconv.Itoa(i)))
	}
}

// TestULIDUnique 并发唯一性
func TestULIDUnique(t *testing.T) {
	const n = 100000
	set := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		s := ULID()
		if _, has := set[s]; has {
			t.Fatalf("出现重复 ULID: %s", s)
		}
		set[s] = struct{}{}
	}
	if len(set) != n {
		t.Fatalf("去重后数量 %d != %d", len(set), n)
	}
}

// TestULIDUniqueParallel 并发唯一性
func TestULIDUniqueParallel(t *testing.T) {
	const (
		worker = 8
		each   = 20000
	)
	ch := make(chan string, worker*each)
	done := make(chan struct{}, worker)
	for i := 0; i < worker; i++ {
		go func() {
			for j := 0; j < each; j++ {
				ch <- ULID()
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < worker; i++ {
		<-done
	}
	close(ch)
	set := make(map[string]struct{}, worker*each)
	for s := range ch {
		if !IsULID(s) {
			t.Fatalf("非法 ULID: %s", s)
		}
		if _, has := set[s]; has {
			t.Fatalf("并发下出现重复 ULID: %s", s)
		}
		set[s] = struct{}{}
	}
}

// TestULIDMonotonic 同毫秒内严格递增 + 字典序 == 时间序
func TestULIDMonotonic(t *testing.T) {
	prev := ULID()
	for i := 0; i < 50000; i++ {
		cur := ULID()
		if cur <= prev {
			t.Fatalf("非严格递增: %s <= %s", cur, prev)
		}
		prev = cur
	}
}

// TestULIDTimeOrder 跨毫秒的字典序与时间序一致
func TestULIDTimeOrder(t *testing.T) {
	var list []string
	for i := 0; i < 5; i++ {
		list = append(list, ULID())
		time.Sleep(2 * time.Millisecond)
	}
	sorted := append([]string(nil), list...)
	sort.Strings(sorted)
	for i := range list {
		if list[i] != sorted[i] {
			t.Fatalf("字典序与时间序不一致: %v vs %v", list, sorted)
		}
	}
}

// TestIsULID 非法输入
func TestIsULID(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"00000000000000000000000000", true},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAV", true},   //大写合法
		{"81ARZ3NDEKTSV4RRFFQ69G5FAV", false},  //首字符 8 越界
		{"01ARZ3NDEKTSV4RRFFQ69G5FA", false},   //长度 25
		{"01ARZ3NDEKTSV4RRFFQ69G5FAVV", false}, //长度 27
		{"01ARZ3NDEKTSV4RRFFQ69G5FAI", false},  //包含 I
		{"01ARZ3NDEKTSV4RRFFQ69G5FAL", false},  //包含 L
		{"01ARZ3NDEKTSV4RRFFQ69G5FAO", false},  //包含 O
		{"01ARZ3NDEKTSV4RRFFQ69G5FAU", false},  //包含 U
		{"01arz3ndektsv4rrffq69g5fav", true},   //小写也合法
	}
	for _, c := range cases {
		if got := IsULID(c.in); got != c.want {
			t.Errorf("IsULID(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// BenchmarkULID ULID 生成性能
func BenchmarkULID(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ULID()
	}
}
