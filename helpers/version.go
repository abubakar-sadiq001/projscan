package helpers

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func PrintVersion() {
	fmt.Printf("ProjScan %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built: %s\n", Date)
}
