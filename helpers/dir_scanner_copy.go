package helpers

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"slices"
// 	"sort"
// 	"strings"
// )

// type FileInfo struct {
// 	Name       string
// 	Lines      int
// 	IsCompFile bool
// }

// type FileStats struct {
// 	Count          int
// 	TotalSize      int64
// 	TotalLines     int
// 	FileInfo       []FileInfo
// 	TotalCompFiles int
// }

// var Stats = make(map[string]*FileStats)
// var processedCount int

// // spinner function
// func showSpinner(stop <-chan bool) {
// 	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
// 	dots := []string{"...", ".", ".."}

// 	i := 0
// 	for {
// 		select {
// 		case <-stop:
// 			fmt.Print("\r\033[K") // Clear the line
// 			return
// 		default:
// 			fmt.Printf("\r%s Scanning project%s ", spinnerChars[i%len(spinnerChars)], dots[i%len(dots)])
// 			i++
// 		}
// 	}
// }

// func ReportPanic() {
// 	p := recover()
// 	if p == nil {
// 		return
// 	}
// 	err, ok := p.(error)
// 	if ok {
// 		fmt.Println(err)
// 	} else {
// 		panic(p)
// 	}
// }

// // Scans through a directory or subdir and return all [.extension] files
// func ScanDirectory(path string, extensionArgsSlice []string) {
// 	// Start spinner in a goroutine
// 	stopSpinner := make(chan bool)
// 	go showSpinner(stopSpinner)
// 	defer close(stopSpinner) // Stop spinner

// 	suffixSlice := []string{".js", ".jsx", ".ts", ".tsx"}

// 	files, err := os.ReadDir(path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Loop over all files
// 	for _, file := range files {
// 		fileName := file.Name()                   // Return current file name
// 		filePath := filepath.Join(path, fileName) // Joins the path of the file with the filename

// 		if file.IsDir() {
// 			ScanDirectory(filePath, extensionArgsSlice)
// 		} else {

// 			fileExtension := filepath.Ext(fileName)

// 			/*
// 				This code does thesame thing with the following  code. it was commented because the following code looks more explicit and shorter

// 					if len(extensionArgsSlice) > 0 {
// 					matchFound := slices.Contains(extensionArgsSlice, fileExtension)

// 					if !matchFound {
// 						continue
// 					}
// 				}
// 			*/

// 			if len(extensionArgsSlice) > 0 && !slices.Contains(extensionArgsSlice, fileExtension) {
// 				continue // skip any file whose extension isn't in the provided list
// 			}

// 			// Checks for nil extension
// 			if fileExtension == "" {
// 				continue
// 			}

// 			// Initialize stats for nil extension
// 			if _, exists := Stats[fileExtension]; !exists {
// 				Stats[fileExtension] = &FileStats{
// 					Count:          0,
// 					TotalSize:      0,
// 					TotalLines:     0,
// 					FileInfo:       []FileInfo{},
// 					TotalCompFiles: 0,
// 				}
// 			}

// 			// Increment extension file count
// 			Stats[fileExtension].Count++
// 			// Get file info
// 			fileInf, err := file.Info()
// 			if err != nil {
// 				log.Printf("Error opening file %s: %v", filePath, err)
// 				continue
// 			}
// 			// Calculate the total size
// 			Stats[fileExtension].TotalSize += fileInf.Size()

// 			// Open and scan file
// 			f, err := os.Open(filePath)
// 			if err != nil {
// 				fmt.Println("Error:", err)
// 				return
// 			}
// 			defer f.Close()

// 			scanner := bufio.NewScanner(f)
// 			const maxCapacity = 10 * 1024 * 1024 // 10MB buffer
// 			buf := make([]byte, maxCapacity)
// 			scanner.Buffer(buf, maxCapacity)

// 			lineCount := 0

// 			for scanner.Scan() {
// 				lineCount++
// 			}

// 			if err := scanner.Err(); err != nil {
// 				log.Printf("Warning: Error reading file %s: %v (skipping line count)\n", filepath.Base(filePath), err)
// 			}

// 			Stats[fileExtension].TotalLines += lineCount

// 			// Check if it's a component file (starts with uppercase AND has valid extension)
// 			isCompFile := false
// 			if fileExtension != "" {
// 				// Check if extension is in the list
// 				for _, suffix := range suffixSlice {
// 					if fileExtension == suffix {
// 						// Check if filename starts with uppercase
// 						if len(fileName) > 0 && fileName[0] >= 'A' && fileName[0] <= 'Z' {
// 							isCompFile = true
// 							break
// 						}
// 					}
// 				}
// 			}

// 			if isCompFile {
// 				Stats[fileExtension].TotalCompFiles++
// 			}

// 			// Store filename and its line count together
// 			Stats[fileExtension].FileInfo = append(Stats[fileExtension].FileInfo, FileInfo{
// 				Name:       fileName,
// 				Lines:      lineCount,
// 				IsCompFile: isCompFile,
// 			})

// 		}
// 	}
// }

// // Print statistics for all extensions
// func PrintStats() {

// 	totalComponentFiles := 0
// 	var extensions []string
// 	for ext := range Stats {
// 		extensions = append(extensions, ext)
// 	}
// 	sort.Strings(extensions)

// 	for _, ext := range extensions {
// 		stat := Stats[ext]
// 		totalComponentFiles += stat.TotalCompFiles
// 	}
// 	if totalComponentFiles > 0 {
// 		fmt.Printf("\n\n")
// 		fmt.Printf("React Components Found: %d\n\n\n", totalComponentFiles)
// 	}

// 	printLargeComponents(Stats, extensions)

// 	// Print files with their line counts
// 	fmt.Println("\n" + strings.Repeat("=", 70))
// 	fmt.Println("FILES PER EXTENSION (with line counts)")
// 	fmt.Println(strings.Repeat("=", 70))

// 	for _, ext := range extensions {
// 		stat := Stats[ext]
// 		fmt.Printf("\n%s (%d files):\n", ext, stat.Count)
// 		for i, info := range stat.FileInfo {
// 			fmt.Printf("  %d. %-30s %8d lines\n", i+1, info.Name, info.Lines)
// 		}
// 	}

// 	fmt.Println(strings.Repeat("=", 70))
// 	fmt.Println("FILE STATISTICS BY EXTENSION")
// 	fmt.Println(strings.Repeat("=", 70))

// 	// Print header
// 	fmt.Printf("%-15s | %-8s | %-12s | %-15s\n",
// 		"Extension", "Count", "Size", "Total Lines")
// 	fmt.Println(strings.Repeat("-", 70))

// 	totalFiles := 0
// 	var totalSize int64
// 	totalLines := 0

// 	for _, ext := range extensions {
// 		stat := Stats[ext]

// 		// Calculate total lines for this extension
// 		extTotalLines := 0
// 		for _, info := range stat.FileInfo {
// 			extTotalLines += info.Lines
// 		}

// 		fmt.Printf("%-15s | %-8d | %-12s | %-15d\n",
// 			ext,
// 			stat.Count,
// 			formatFileSize(stat.TotalSize),
// 			extTotalLines,
// 		)

// 		totalFiles += stat.Count
// 		totalSize += stat.TotalSize
// 		totalLines += extTotalLines
// 		totalComponentFiles += stat.TotalCompFiles

// 	}

// 	fmt.Println(strings.Repeat("-", 70))
// 	fmt.Printf("%-15s | %-8d | %-12s | %-15d\n",
// 		"TOTAL",
// 		totalFiles,
// 		formatFileSize(totalSize),
// 		totalLines,
// 	)

// }

// // formatFileSize with proper unit display
// func formatFileSize(bytes int64) string {
// 	const unit = 1024

// 	if bytes < unit {
// 		return fmt.Sprintf("%d bytes", bytes)
// 	}

// 	// Determine the appropriate unit
// 	units := []string{"bytes", "kb", "MB", "GB", "TB"}
// 	exp := 0
// 	size := float64(bytes)

// 	for size >= unit && exp < len(units)-1 {
// 		size /= unit
// 		exp++
// 	}

// 	// If still in bytes, show no decimal places
// 	if exp == 0 {
// 		return fmt.Sprintf("%d %s", int64(size), units[exp])
// 	}

// 	// For KB and above, show 2 decimal places
// 	return fmt.Sprintf("%.2f %s", size, units[exp])
// }

// // I go add comment later
// func printLargeComponents(Stats map[string]*FileStats, extensions []string) {
// 	var largeComponents []struct {
// 		name  string
// 		lines int
// 	}

// 	// Single loop to collect all large components
// 	for _, ext := range extensions {
// 		for _, info := range Stats[ext].FileInfo {
// 			if info.Lines > 300 && info.IsCompFile {
// 				largeComponents = append(largeComponents, struct {
// 					name  string
// 					lines int
// 				}{info.Name, info.Lines})
// 			}
// 		}
// 	}

// 	count := len(largeComponents)
// 	if count == 0 {
// 		return
// 	}

// 	// Print header
// 	fmt.Println("Large Components")
// 	fmt.Println(strings.Repeat("-", 16))

// 	// Print all components
// 	for _, comp := range largeComponents {
// 		fmt.Printf("%-8s %8d lines\n", comp.name, comp.lines)
// 	}

// 	// Print warning
// 	textValue := "components"
// 	if count <= 1 {
// 		textValue = "component"
// 	}

// 	fmt.Printf("\n")
// 	fmt.Printf(`⚠️  Warning:
//        %d %s exceed 300 lines.
//        Consider splitting them into smaller components.`, count, textValue)
// }
