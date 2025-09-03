package config

import "go-echo-template/internal/shared/utils"

type ObjectConfig struct {
	S3 *S3Config
}

type S3Config struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

func newS3Config() *S3Config {
	return &S3Config{
		Region:    utils.MustGetStrEnv("S3_REGION"),
		Bucket:    utils.MustGetStrEnv("S3_BUCKET"),
		AccessKey: utils.MustGetStrEnv("S3_ACCESS_KEY"),
		SecretKey: utils.MustGetStrEnv("S3_SECRET_KEY"),
		Endpoint:  utils.MustGetStrEnv("S3_ENDPOINT"),
	}
}

func newObjectConfig() *ObjectConfig {
	return &ObjectConfig{
		S3: newS3Config(),
	}
}
