package config

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOss struct {
	EndPoint        string
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
}

func NewAliyunOSS(name string) (*oss.Client, string, error) {
	ossConf := AliyunOss{}
	err := conf.ReadSection(name, &ossConf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read config for %s section: %s", name, err)
	}
	client, err := oss.New(ossConf.EndPoint, ossConf.AccessKeyId, ossConf.AccessKeySecret)
	if err != nil {
		return nil, "", err
	}
	return client, ossConf.BucketName, nil
}
