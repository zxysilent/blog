package upload

import (
	"blog/conf"
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
	"time"
)

type Qiniu struct{}

func (*Qiniu) UploadFile(file *multipart.FileHeader) (string, string, error) {
	putPolicy := storage.PutPolicy{Scope: conf.App.Oss.QiniuBucket}
	mac := qbox.NewMac(conf.App.Oss.QiniuAccessKey, conf.App.Oss.QiniuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfig()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "logo"}}

	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("file open failed, err:" + openError.Error())
	}
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", "", errors.New("file put qiniu failed, err:" + putErr.Error())
	}
	return conf.App.Oss.QiniuImgPath + "/" + ret.Key, ret.Key, nil
}

func (*Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(conf.App.Oss.QiniuAccessKey, conf.App.Oss.QiniuSecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(conf.App.Oss.QiniuBucket, key); err != nil {
		return errors.New("delete qiniu file failed, err:" + err.Error())
	}
	return nil
}

func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      conf.App.Oss.QiniuUseHTTPS,
		UseCdnDomains: conf.App.Oss.QiniuUseCdnDomains,
	}
	switch conf.App.Oss.QiniuZone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &cfg
}
