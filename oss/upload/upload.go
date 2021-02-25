package upload

import (
	"blog/conf"
	"mime/multipart"
)

type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewOss() OSS {
	switch conf.App.Oss.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	case "aliyun-oss":
		return &AliyunOSS{}
	default:
		return &Local{}
	}
}
