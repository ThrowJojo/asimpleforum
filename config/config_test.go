package config

import (
	"testing"
	"fmt"
)

func TestGetDBConfig(t *testing.T) {
	if dbConfig, err := LoadConfigData(); err != nil {
		t.Error("Unexpected error getting DB Config: ", err)
	} else {
		fmt.Println(dbConfig)
	}
}

func TestGetConnectionString(t *testing.T) {
	if connectionString, err := GetConnectionString(); err != nil {
		t.Error("Unexpected error getting connection string: ", err)
	} else {
		fmt.Println("Connection string: ", connectionString)
	}
}