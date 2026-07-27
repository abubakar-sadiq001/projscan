package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// FileInfo stores summary information about a single file discovered during scanning.
type FileInfo struct {
	// Name is the filename as it appears on disk.
	Name string
	// Lines is the number of lines found in the file.
	Lines int
	// IsCompFile indicates whether the file looks like a React-style component.
	IsCompFile bool
}

// FileStats aggregates cumulative statistics for one file extension.
type FileStats struct {
	// Count is the number of files seen for this extension.
	Count int
	// TotalSize is the combined size in bytes of all files with this extension.
	TotalSize int64
	// TotalLines is the combined line count for all files with this extension.
	TotalLines int
	// FileInfo contains the per-file details used for reporting.
	FileInfo []FileInfo
	// TotalCompFiles counts how many files in this extension were detected as component-like.
	TotalCompFiles int
}

// Stats stores scanning results grouped by file extension.
var Stats = make(map[string]*FileStats)

// processedCount is currently unused, but reserved for future progress-tracking features.
var processedCount int

// showSpinner prints a simple animated spinner while scanning is in progress.
func showSpinner(stop <-chan bool) {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	dots := []string{"...", ".", ".."}

	i := 0
	for {
		select {
		case <-stop:
			// Clear the spinner line before returning.
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r%s Scanning project%s ", spinnerChars[i%len(spinnerChars)], dots[i%len(dots)])
			i++
		}
	}
}

// ReportPanic recovers from a panic and prints the underlying error if one occurred.
func ReportPanic() {
	p := recover()
	if p == nil {
		return
	}
	err, ok := p.(error)
	if ok {
		fmt.Println(err)
	} else {
		panic(p)
	}
}

// ScanDirectory recursively walks a directory tree and records stats for matching files.
func ScanDirectory(path string, extensionArgsSlice []string) {
	// Start a background spinner so users get feedback during the scan.
	stopSpinner := make(chan bool)
	go showSpinner(stopSpinner)
	defer close(stopSpinner)

	// These extensions are considered possible component-style files.
	suffixSlice := []string{".js", ".jsx", ".ts", ".tsx"}

	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// Walk every entry in the current directory.
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(path, fileName)

		if file.IsDir() {
			// Recurse into subdirectories so the scan covers the whole project.
			ScanDirectory(filePath, extensionArgsSlice)
		} else {
			fileExtension := filepath.Ext(fileName)

			// If the user supplied extension filters, skip anything that does not match.
			if len(extensionArgsSlice) > 0 && !slices.Contains(extensionArgsSlice, fileExtension) {
				continue
			}

			// Skip files that do not have an extension.
			if fileExtension == "" {
				continue
			}

			// Create a stats bucket for this extension the first time it appears.
			if _, exists := Stats[fileExtension]; !exists {
				Stats[fileExtension] = &FileStats{
					Count:          0,
					TotalSize:      0,
					TotalLines:     0,
					FileInfo:       []FileInfo{},
					TotalCompFiles: 0,
				}
			}

			// Count this file under the matching extension category.
			Stats[fileExtension].Count++

			// Read file metadata so we can track total size.
			fileInf, err := file.Info()
			if err != nil {
				log.Printf("Error opening file %s: %v", filePath, err)
				continue
			}
			Stats[fileExtension].TotalSize += fileInf.Size()

			// Open the file and count its lines.
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			const maxCapacity = 10 * 1024 * 1024 // 10MB buffer
			buf := make([]byte, maxCapacity)
			scanner.Buffer(buf, maxCapacity)

			lineCount := 0

			for scanner.Scan() {
				lineCount++
			}

			if err := scanner.Err(); err != nil {
				log.Printf("Warning: Error reading file %s: %v (skipping line count)\n", filepath.Base(filePath), err)
			}

			// Add the current file's line count to the extension total.
			Stats[fileExtension].TotalLines += lineCount

			// Determine whether this file should be considered a component-like file.
			isCompFile := false
			if fileExtension != "" {
				for _, suffix := range suffixSlice {
					if fileExtension == suffix {
						// Component files are typically named with an uppercase first letter.
						if len(fileName) > 0 && fileName[0] >= 'A' && fileName[0] <= 'Z' {
							isCompFile = true
							break
						}
					}
				}
			}

			if isCompFile {
				Stats[fileExtension].TotalCompFiles++
			}

			// Store the file's details so the report can print them later.
			Stats[fileExtension].FileInfo = append(Stats[fileExtension].FileInfo, FileInfo{
				Name:       fileName,
				Lines:      lineCount,
				IsCompFile: isCompFile,
			})
		}
	}
}

// PrintStats prints a human-readable report of the scan results.
func PrintStats() {
	// Count component files across all extensions for the summary header.
	totalComponentFiles := 0
	var extensions []string
	for ext := range Stats {
		extensions = append(extensions, ext)
	}
	sort.Strings(extensions)

	for _, ext := range extensions {
		stat := Stats[ext]
		totalComponentFiles += stat.TotalCompFiles
	}

	if totalComponentFiles > 0 {
		fmt.Printf("\n\n")
		fmt.Printf("React Components Found: %d\n\n\n", totalComponentFiles)
	}

	printLargeComponents(Stats, extensions)

	// Print the per-file breakdown for each extension.
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("FILES PER EXTENSION (with line counts)")
	fmt.Println(strings.Repeat("=", 70))

	for _, ext := range extensions {
		stat := Stats[ext]
		fmt.Printf("\n%s (%d files):\n", ext, stat.Count)
		for i, info := range stat.FileInfo {
			fmt.Printf("  %d. %-30s %8d lines\n", i+1, info.Name, info.Lines)
		}
	}

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

	for _, ext := range extensions {
		stat := Stats[ext]

		// Recalculate the total lines for this extension from the stored file details.
		extTotalLines := 0
		for _, info := range stat.FileInfo {
			extTotalLines += info.Lines
		}

		fmt.Printf("%-15s | %-8d | %-12s | %-15d\n",
			ext,
			stat.Count,
			formatFileSize(stat.TotalSize),
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
		formatFileSize(totalSize),
		totalLines,
	)
}

// formatFileSize converts byte counts into a more human-readable format such as KB or MB.
func formatFileSize(bytes int64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d bytes", bytes)
	}

	// Determine the appropriate unit for the size.
	units := []string{"bytes", "kb", "MB", "GB", "TB"}
	exp := 0
	size := float64(bytes)

	for size >= unit && exp < len(units)-1 {
		size /= unit
		exp++
	}

	// If the size is still in bytes, show it without decimals.
	if exp == 0 {
		return fmt.Sprintf("%d %s", int64(size), units[exp])
	}

	// For larger units, show two decimal places for readability.
	return fmt.Sprintf("%.2f %s", size, units[exp])
}

// printLargeComponents highlights component files that are unusually large.
func printLargeComponents(Stats map[string]*FileStats, extensions []string) {
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
