package main

import (
	"fmt"
	"gitlabFileScanner/internal/app"
	"os"
)

func main() {
	err := app.Start()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
