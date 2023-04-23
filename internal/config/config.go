package config

import (
	"github.com/joho/godotenv"
	"os"
)

var configMap map[string]string

func InitConfig() error {
	err := godotenv.Load(".env")

	configMap = make(map[string]string)
	dbConfig(configMap)
	fiasConfig(configMap)

	return err
}

func GetConfig(key string) string {
	val, ok := configMap[key]
	if !ok {
		return ""
	}
	envVal := os.Getenv(key)
	if envVal != "" {
		return envVal
	}
	return val
}
