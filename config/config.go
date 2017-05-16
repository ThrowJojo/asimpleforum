package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ConfigData struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	TestDatabase string `yaml:"test_database"`
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

	configData := &ConfigData{User: viper.GetString("user"),
		Password: viper.GetString("password"),
		Database: viper.GetString("database"),
		TestDatabase: viper.GetString("test_database"),
		Secret: viper.GetString("secret")}
	return configData, nil

}

// Formats a connection with data loaded from config.yaml
func GetConnectionString(test bool) (string, error) {
	if dbConfig, err := LoadConfigWithViper(); err != nil {
		return "", err
	} else {
		if test {
			return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.TestDatabase), nil
		} else {
			return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Database), nil
		}
	}
}