package targz

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Targz 压缩文件或目录为 .tar.gz
// dst 未以 .tar.gz 结尾时自动补全
func Targz(dst string, src string) error {
	if !strings.HasSuffix(dst, ".tar.gz") {
		dst += ".tar.gz"
	}
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("源路径不可用: %w", err)
	}
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	// 关闭顺序必须是 tw → gw → d，否则尾部数据写不进去
	gw, err := gzip.NewWriterLevel(d, 5)
	if err != nil {
		d.Close()
		return err
	}
	tw := tar.NewWriter(gw)
	walkErr := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		// 记录相对路径，避免解压时带出源机器上的绝对目录
		name, err := filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(name)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
	// 无论遍历是否出错都要把尾部写完，再按顺序关闭
	if err := tw.Close(); err != nil && walkErr == nil {
		walkErr = err
	}
	if err := gw.Close(); err != nil && walkErr == nil {
		walkErr = err
	}
	if err := d.Close(); err != nil && walkErr == nil {
		walkErr = err
	}
	if walkErr != nil {
		// 压缩失败时不留半成品
		os.Remove(dst)
	}
	return walkErr
}

// UnTargz 解压 .tar.gz 到 dst 目录
// 只允许写入 dst 之下的路径，防止压缩包内含 ../ 造成越界写入
func UnTargz(dst string, src string) error {
	if dst == "" {
		return errors.New("目标目录不能为空")
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	root, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	// 注意：这里必须是 Open，Create 会把待解压的文件截断为空
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		target, err := safeJoin(root, hdr.Name)
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		case tar.TypeReg:
		default:
			// 符号链接等类型直接跳过，避免不可预期的落盘行为
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode)&0o777)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tr)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

// safeJoin 把归档内的相对路径拼到 root 下，拒绝越界路径
func safeJoin(root, name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("非法的归档路径: %s", name)
	}
	target := filepath.Join(root, clean)
	rel, err := filepath.Rel(root, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("非法的归档路径: %s", name)
	}
	return target, nil
}
