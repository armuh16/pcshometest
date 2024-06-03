package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	str2duration "github.com/xhit/go-str2duration/v2"
)

var c *Config

const (
	Development = "Development"
	Production  = "Production"
)

type Config struct {
	Env      string `yaml:"env"`
	AppURL   string `yaml:"appURL"`
	Postgres struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslMode"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbName"`
	} `yaml:"postgres"`
	Auth struct {
		ExpireAccessToken          string `yaml:"expireAccessToken"`
		ExpireRefreshToken         string `yaml:"expireRefreshToken"`
		Secret                     string `yaml:"secret"`
		ExpireAccessTokenDuration  time.Duration
		ExpireRefreshTokenDuration time.Duration
	} `yaml:"auth"`
}

func Get() *Config {
	return c
}

func SetConfig() {
	var err error

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err.Error())
	}

	c.Auth.ExpireAccessTokenDuration, err = str2duration.ParseDuration(c.Auth.ExpireAccessToken)
	if err != nil {
		panic(fmt.Sprintf("config auth access expired duration string not valid: %s", err.Error()))
	}

	c.Auth.ExpireRefreshTokenDuration, err = str2duration.ParseDuration(c.Auth.ExpireRefreshToken)
	if err != nil {
		panic(fmt.Sprintf("config auth refresh token duration string not valid: %s", err.Error()))
	}

	viper.WatchConfig()
}
