package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Targz 压缩文件 .tar.gz
func Targz(dst string, src string) error {
	if !strings.HasSuffix(dst, ".tar.gz") {
		dst += ".tar.gz"
	}
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	gw, _ := gzip.NewWriterLevel(d, 5)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			file.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

// 解压文件
func UnTargz(dst string, src string) error {
	if dst[len(dst)-1] != '/' {
		dst += "/"
	}
	file, err := os.Create(src)
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
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dst + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}
func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(name[0:strings.LastIndex(name, "/")], 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
