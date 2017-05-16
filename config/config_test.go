package config

import (
	"testing"
	"fmt"
)

func TestGetConnectionString(t *testing.T) {
	if connectionString, err := GetConnectionString(true); err != nil {
		t.Error("Unexpected error getting connection string: ", err)
	} else {
		fmt.Println("Connection string: ", connectionString)
	}
}

func TestLoadConfigWithViper(t *testing.T) {
	LoadConfigWithViper()
}