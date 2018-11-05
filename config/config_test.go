package config

import (
	"testing"

	"github.com/article/config"
)

func TestGetConfig(t *testing.T) {
	config := config.GetConfig()

	if config.Mongo.Ip != "127.0.0.1:27017" {
		t.Errorf("Enabled Mongodb deafult Ip does not match with expected value")
	}

	if config.Mongo.DbName != "article" {
		t.Errorf("Enabled Mongodb deafult dbname does not match with expected value")
	}
	if config.Server.Port != ":8001" {
		t.Errorf("Enabled Server Default port is not match with expected value")
	}
}
