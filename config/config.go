package config

import "os"

// MongoConfig ...
type MongoConfig struct {
	Ip     string `json:"ip"`
	DbName string `json:"dbName"`
}

// ServerConfig ....
type ServerConfig struct {
	Port string `json:"port"`
}

// Config ...
type Config struct {
	Mongo  *MongoConfig  `json:"mongo"`
	Server *ServerConfig `json:"server"`
}

// GetConfig ...
func GetConfig() *Config {
	return &Config{
		Mongo: &MongoConfig{
			Ip:     envOrDefaultString("article:mongo:ip", "127.0.0.1:27017"),
			DbName: envOrDefaultString("article:mongo:dbName", "article")},
		Server: &ServerConfig{Port: envOrDefaultString("article:server:port", ":8001")},
	}
}

func envOrDefaultString(envVar string, defaultValue string) string {

	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}

	return value
}
