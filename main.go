/*
Copyright © 2026 Abubakar Sadiq Abubakarsadiqi532@gmail.com
*/
package main

import (
	"os"

	"github.com/abubakar-sadiq001/projscan/cmd"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" || os.Args[1] == "-v" {
		PrintVersion()
		return
	}

	cmd.Execute()
}
