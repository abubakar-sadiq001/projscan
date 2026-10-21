package helpers

import (
	"fmt"
	"sort"
	"strings"
)

func (f *FileStats) PrintFilesWithLines() {
	// Count component files across all extensions for the summary header.

	for ext := range Stats {
		f.Extensions = append(f.Extensions, ext)
	}
	sort.Strings(f.Extensions)

	for _, ext := range f.Extensions {
		stat := Stats[ext]
		totalComponentFiles += stat.TotalCompFiles
	}

	if totalComponentFiles > 0 {
		fmt.Printf("\n\n")
		fmt.Printf("React Components Found: %d\n\n\n", totalComponentFiles)
	}

	// highlights component files that are unusually large.
	f.printLargeComponentFiles(Stats, f.Extensions)

	// Print the per-file breakdown for each extension.
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("FILES PER EXTENSION (with line counts)")
	fmt.Println(strings.Repeat("=", 70))

	for _, ext := range f.Extensions {
		stat := Stats[ext]
		fmt.Printf("\n%s (%d files):\n", ext, stat.Count)
		for i, info := range stat.FileInfo {
			fmt.Printf("  %d. %-30s %8d lines\n", i+1, info.Name, info.Lines)
		}
	}
}
