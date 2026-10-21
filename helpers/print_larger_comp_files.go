package helpers

import (
	"fmt"
	"strings"
)

func (f *FileStats) printLargeComponentFiles(Stats map[string]*FileStats, extensions []string) {
	var largeComponents []struct {
		name  string
		lines int
	}

	// Collect every component file that exceeds the warning threshold.
	for _, ext := range extensions {
		for _, info := range Stats[ext].FileInfo {
			if info.Lines > 300 && info.IsCompFile {
				largeComponents = append(largeComponents, struct {
					name  string
					lines int
				}{info.Name, info.Lines})
			}
		}
	}

	count := len(largeComponents)
	if count == 0 {
		return
	}

	// Print the header for the large-component section.
	fmt.Println("Large Components")
	fmt.Println(strings.Repeat("-", 16))

	// Print each large component with its line count.
	for _, comp := range largeComponents {
		fmt.Printf("%-8s %8d lines\n", comp.name, comp.lines)
	}

	// Print a warning message about the large components.
	textValue := "components"
	if count <= 1 {
		textValue = "component"
	}

	fmt.Printf("\n")
	fmt.Printf(`⚠️  Warning:
       %d %s exceed 300 lines.
       Consider splitting them into smaller components.`, count, textValue)
}
