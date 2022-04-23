package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

/*
Go App Profile Setting using Viper pkg
*/
func SetAppProfile(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("환경설정 파일을 찾을 수 없습니다.")
		}
		return nil, err
	}
	return v, nil
}

/*
Viper Unmarshal is all setup in the decode Config struct
*/
func ParseViper(v *viper.Viper) (*Config, error) {
	var c *Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return c, nil
}