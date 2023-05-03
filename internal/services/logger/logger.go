package logger

import (
	"log"
	"os"
	"time"
)

var (
	LogService *log.Logger
	LogFile    *os.File
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t := time.Now()

	logpath := pwd + "/log/" + t.Format("2006-01-02") + ".log"
	if _, err := os.Stat(logpath); err == nil {
		LogFile, err = os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		LogFile, err = os.Create(logpath)
	}

	if err != nil {
		panic(err)
	}
	LogService = log.New(LogFile, "", log.LstdFlags|log.Lshortfile)
}

func Println(logMessage ...any) {
	LogService.Println(logMessage...)
	log.Println(logMessage...)
}
