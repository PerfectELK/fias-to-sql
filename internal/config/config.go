package config

import (
	"github.com/joho/godotenv"
	"os"
)

var configMap map[string]string
var configMapRedeclaredWithApp map[string]bool

func InitConfig() error {
	err := godotenv.Load(".env")

	configMap = make(map[string]string)
	configMapRedeclaredWithApp = make(map[string]bool)

	appConfig(configMap)
	dbConfig(configMap)
	fiasConfig(configMap)

	return err
}

func GetConfig(key string) string {
	val, ok := configMap[key]
	if !ok {
		return ""
	}
	_, isRedeclared := configMapRedeclaredWithApp[key]
	if isRedeclared {
		return val
	}
	envVal := os.Getenv(key)
	if envVal != "" {
		return envVal
	}
	return val
}

func SetConfig(key string, value string) {
	configMapRedeclaredWithApp[key] = true
	configMap[key] = value
}
