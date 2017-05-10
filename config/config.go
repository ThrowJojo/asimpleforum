package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"github.com/spf13/viper"
)

type ConfigData struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Secret string `yaml:"secret"`
}

// Loads config.yaml file with viper
func LoadConfigWithViper() (*ConfigData, error) {

	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	configData := &ConfigData{User: viper.GetString("user"), Password: viper.GetString("password"), Database: viper.GetString("database"), Secret: viper.GetString("secret")}
	return configData, nil

}

// Formats a connection with data loaded from config.yaml
func GetConnectionString() (string, error) {
	if dbConfig, err := LoadConfigWithViper(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Database), nil
	}
}