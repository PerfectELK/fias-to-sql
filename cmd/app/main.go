package main

import (
	"fias_to_sql/internal/app"
	"fias_to_sql/internal/services/logger"
	"fmt"
)

func main() {
	defer logger.LogFile.Close()
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
