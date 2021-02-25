package upload

import (
	"blog/conf"
	"path/filepath"

	"errors"
	"github.com/zxysilent/utils"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type Local struct{}

func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", errors.New("open file failed, err:" + err.Error())
	}
	defer src.Close()
	basePath := conf.App.Oss.LocalPath + time.Now().Format(utils.FmtyyyyMMdd) + "/"
	//确保文件夹存在
	os.MkdirAll(basePath, 0777)
	fileName := utils.RandStr(16) + filepath.Ext(file.Filename)
	filePathName := basePath + fileName
	dst, err := os.Create(filePathName)
	if err != nil {
		return "", "", errors.New("create file failed, err:" + err.Error())
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", "", errors.New("copy file failed, err" + err.Error())
	}

	return "/" + filePathName, fileName, nil
}

func (*Local) DeleteFile(key string) error {
	p := conf.App.Oss.LocalPath + "/" + key
	if strings.Contains(p, conf.App.Oss.LocalPath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
