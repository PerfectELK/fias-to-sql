package config

import (
	"github.com/joho/godotenv"
	"os"
)

var configMap map[string]string
var configMapRedeclaredWithApp map[string]bool

func InitConfig(isLoadEnv ...bool) error {
	IsLoadEnv := true
	if len(isLoadEnv) > 0 {
		IsLoadEnv = isLoadEnv[0]
	}

	if IsLoadEnv {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	configMap = make(map[string]string)
	configMapRedeclaredWithApp = make(map[string]bool)

	appConfig(configMap)
	dbConfig(configMap)
	fiasConfig(configMap)

	return nil
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
