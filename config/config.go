package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type ConfigData struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Secret string `yaml:"secret"`
}

// TODO: Change database name, secret
func LoadConfigData() (*ConfigData, error) {

	dbConfig := ConfigData{}

	contents, readErr := ioutil.ReadFile("../config.yaml")
	if readErr != nil {
		return nil, readErr
	}

	if yamlError := yaml.Unmarshal([]byte(contents), &dbConfig); yamlError != nil {
		return nil, yamlError
	}

	return &dbConfig, nil

}

func GetConnectionString() (string, error) {
	if dbConfig, err := LoadConfigData(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Database), nil
	}
}