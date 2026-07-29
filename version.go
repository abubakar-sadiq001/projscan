package main

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func PrintVersion() {
	fmt.Printf("ProjScan %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built: %s\n", date)
}
