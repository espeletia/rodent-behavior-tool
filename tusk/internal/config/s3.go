package config

import "fmt"

type S3Config struct {
	Credentials       string
	AccessKey         string
	SecretAccessKey   string
	URL               string
	Region            string
	UploadsPathPrefix string
	Bucket            string
	Prod              bool
}

func loadS3Config() S3Config {
	s3Config := &S3Config{}
	v := configViper("s3")
	err := v.BindEnv("URL", "S3_URL")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.BindEnv("Credentials", "S3_CREDENTIALS")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.BindEnv("AccessKey", "S3_ACCESS_KEY")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.BindEnv("SecretAccessKey", "S3_SECRET_ACCESS_KEY")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.BindEnv("Bucket", "BUCKET")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.BindEnv("Prod", "PROD_ENV")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = v.Unmarshal(s3Config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return *s3Config
}
