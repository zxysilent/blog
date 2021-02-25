package upload

import (
	"blog/conf"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"path/filepath"
	"time"
)

type AliyunOSS struct{}

func (*AliyunOSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	bucket, err := newBucket()
	if err != nil {
		return "", "", errors.New("new aliyun oss bucket failed, err:" + err.Error())
	}

	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("file open failed, err:" + openError.Error())
	}

	//上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename

	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		return "", "", errors.New("file put aliyun oss failed, err:" + err.Error())
	}

	return conf.App.Oss.AliyunOssBucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (*AliyunOSS) DeleteFile(key string) error {
	bucket, err := newBucket()
	if err != nil {
		return errors.New("new aliyun oss bucket failed, err:" + err.Error())
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		return errors.New("delete aliyun oss file failed, err:" + err.Error())
	}

	return nil
}

func newBucket() (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(conf.App.Oss.AliyunOssEndpoint, conf.App.Oss.AliyunOssAccessKeyId, conf.App.Oss.AliyunOssAccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(conf.App.Oss.AliyunOssBucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
