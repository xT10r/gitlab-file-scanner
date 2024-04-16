package main

import (
	"fmt"
	"gitlabFileScanner/cmd"
	"os"
)

func main() {
	err := cmd.Start()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
