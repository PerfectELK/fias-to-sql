package config

import "os"

func appConfig(m map[string]string) {
	p, _ := os.Getwd()
	m["APP_DEBUG"] = ""
	m["APP_THREAD_NUMBER"] = ""
	m["APP_ROOT"] = p
	os.Setenv("APP_ROOT", p)
}
