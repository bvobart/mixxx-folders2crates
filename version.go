package main

import "fmt"

var (
	version = "dev-snapshot"
	commit  = "unknown"
	date    = "latest"
)

func printVersion() {
	fmt.Println("mixxx-folders2crates version:", version)
	fmt.Println("commit:", commit)
	fmt.Println("date:", date)
}
