package config

import (
	"github.com/spf13/viper"
	"log"
)

type Service struct {
	Address string `yaml:"address"`
}

type Config struct {
	Services struct {
		Api      Service `yaml:"api"`
		Account  Service `yaml:"account"`
		Matching Service `yaml:"matching"`
		Swipes   Service `yaml:"swipes"`
	} `yaml:"services"`
	Database struct {
		Account struct {
			Dsn string `yaml:"dsn"`
		} `yaml:"account"`
		Swipes struct {
			Dsn string `yaml:"dsn"`
		} `yaml:"swipes"`
		Redis struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Password string `yaml:"password"`
			Username string `yaml:"username"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"database"`
	Auth struct {
		Jwt string `yaml:"jwt"`
	} `yaml:"auth"`
	Public struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"database"`
	} `yaml:"public"`
	S3 struct {
		Bucket   string `yaml:"bucket"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"s3"`
}

func LoadConfig(path, mode string) *Config {
	viper.SetConfigName("config." + mode)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("bad path to config: '%s', mode: %s", path, mode)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("YAML parsing error")
	}
	return &config
}

func (conf *Config) GetRedisAddr() string {
	return conf.Database.Redis.Host + ":" + conf.Database.Redis.Port
}
