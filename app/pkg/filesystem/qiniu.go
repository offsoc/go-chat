package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	_ "github.com/qiniu/go-sdk/v7/storage"
	"go-chat/config"
	"path/filepath"
	"strings"
	"time"
)

// QiniuFilesystem
// @link 对接文档 https://developer.qiniu.com/kodo/1238/go#upload-flow
type QiniuFilesystem struct {
	conf *config.Config
	mac  *qbox.Mac
}

func NewQiniuFilesystem(conf *config.Config) *QiniuFilesystem {
	return &QiniuFilesystem{
		conf: conf,
		mac:  qbox.NewMac(conf.Filesystem.Qiniu.AccessKey, conf.Filesystem.Qiniu.SecretKey),
	}
}

// Token 获取上传凭证
// todo token 需要加入缓存
func (s *QiniuFilesystem) Token() string {
	putPolicy := storage.PutPolicy{
		Scope:   s.conf.Filesystem.Qiniu.Bucket,
		Expires: 7200,
	}

	return putPolicy.UploadToken(s.mac)
}

func (s *QiniuFilesystem) Write(data []byte, filePath string) error {
	filePath = strings.TrimLeft(filePath, "/")

	cfg := storage.Config{
		Zone:     &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS: true,                 // 是否使用https域名
	}

	// 七牛标准的上传回复内容
	ret := storage.PutRet{}

	// 可选配置
	params := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)

	err := formUploader.Put(context.Background(), &ret, s.Token(), filePath, bytes.NewReader(data), int64(len(data)), &params)
	if err != nil {
		return err
	}

	fmt.Println(ret.Key, ret.Hash)

	return nil
}

func (s *QiniuFilesystem) WriteLocal(localFile string, filePath string) error {
	cfg := storage.Config{
		Zone:     &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS: true,                 // 是否使用https域名
	}

	// 七牛标准的上传回复内容
	ret := storage.PutRet{}

	// 可选配置
	params := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)

	filePath = strings.TrimLeft(filePath, "/")

	if err := formUploader.PutFile(context.Background(), &ret, s.Token(), filePath, localFile, &params); err != nil {
		return err
	}

	return nil
}

func (s *QiniuFilesystem) Copy(srcPath, filePath string) error {
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuadong,
	}

	bucket := s.conf.Filesystem.Qiniu.Bucket

	bucketManager := storage.NewBucketManager(s.mac, &cfg)

	return bucketManager.Copy(bucket, srcPath, bucket, filePath, false)
}

func (s *QiniuFilesystem) Delete(filePath string) error {
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuadong,
	}

	bucket := s.conf.Filesystem.Qiniu.Bucket

	bucketManager := storage.NewBucketManager(s.mac, &cfg)

	if err := bucketManager.Delete(bucket, filePath); err != nil {
		return err
	}

	return nil
}

func (s *QiniuFilesystem) DeleteDir(path string) error {
	return errors.New("七牛云无删除文件夹接口")
}

func (s *QiniuFilesystem) CreateDir(path string) error {
	return errors.New("七牛云无创建文件夹接口")
}

func (s *QiniuFilesystem) Stat(filePath string) (*FileStat, error) {
	cfg := storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuadong,
	}

	bucket := s.conf.Filesystem.Qiniu.Bucket

	bucketManager := storage.NewBucketManager(s.mac, &cfg)

	fileInfo, err := bucketManager.Stat(bucket, filePath)
	if err != nil {
		return nil, err
	}

	return &FileStat{
		Name:        filepath.Base(filePath),
		Size:        fileInfo.Fsize,
		Ext:         filepath.Ext(filePath),
		MimeType:    fileInfo.MimeType,
		LastModTime: storage.ParsePutTime(fileInfo.PutTime),
	}, nil
}

func (s *QiniuFilesystem) PublicUrl(filePath string) string {
	return storage.MakePublicURL(s.conf.Filesystem.Qiniu.Domain, filePath)
}

func (s *QiniuFilesystem) PrivateUrl(filePath string, timeout int) string {
	deadline := time.Now().Add(time.Second * time.Duration(timeout)).Unix()

	return storage.MakePrivateURL(s.mac, s.conf.Filesystem.Qiniu.Domain, filePath, deadline)
}

// List 获取目录下的所有文件
func (s *QiniuFilesystem) List(dir string) {

}
