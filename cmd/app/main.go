package main

import (
	"fias_to_sql/internal/app"
	"fmt"
)

func main() {
	err := app.App()
	if err != nil {
		fmt.Println(err)
	}
}
