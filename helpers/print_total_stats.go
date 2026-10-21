package helpers

import (
	"fmt"
	"strings"
)

func (f *FileStats) PrintTotalFilesStat() {

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("FILE STATISTICS BY EXTENSION")
	fmt.Println(strings.Repeat("=", 70))

	// Print the summary table header.
	fmt.Printf("%-15s | %-8s | %-12s | %-15s\n",
		"Extension", "Count", "Size", "Total Lines")
	fmt.Println(strings.Repeat("-", 70))

	totalFiles := 0
	var totalSize int64
	totalLines := 0

	for _, ext := range f.Extensions {
		stat := Stats[ext]

		// Recalculate the total lines for this extension from the stored file details.
		extTotalLines := 0
		for _, info := range stat.FileInfo {
			extTotalLines += info.Lines
		}

		fmt.Printf("%-15s | %-8d | %-12s | %-15d\n",
			ext,
			stat.Count,
			f.formatFileSize(stat.TotalSize),
			extTotalLines,
		)

		totalFiles += stat.Count
		totalSize += stat.TotalSize
		totalLines += extTotalLines
		totalComponentFiles += stat.TotalCompFiles
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%-15s | %-8d | %-12s | %-15d\n",
		"TOTAL",
		totalFiles,
		f.formatFileSize(totalSize),
		totalLines,
	)
}
