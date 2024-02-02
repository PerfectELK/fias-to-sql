package main

import (
	"fmt"
	"github.com/PerfectELK/go-import-fias/internal/app"
	"github.com/PerfectELK/go-import-fias/internal/services/logger"
)

func main() {
	defer logger.LogFile.Close()
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
