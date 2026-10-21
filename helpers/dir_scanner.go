package helpers

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var totalComponentFiles int
var Mu *sync.Mutex

// FileInfo stores summary information about a single file discovered during scanning.
type FileInfo struct {
	Name       string
	Lines      int
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
	Extensions     []string
}

func InitializeStats() *FileStats {
	return &FileStats{
		Count:          0,
		TotalSize:      0,
		TotalLines:     0,
		FileInfo:       []FileInfo{},
		TotalCompFiles: 0,
		Extensions:     []string{},
	}
}

// Stats stores scanning results grouped by file extension.
var Stats = make(map[string]*FileStats)

var sem = make(chan struct{}, 16)

// ScanDirectory recursively walks a directory tree and records stats for matching files.
func ScanDirectory(path string, extensionArgsSlice []string, wg *sync.WaitGroup) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() {
		<-sem
	}()

	suffixSlice := []string{".js", ".jsx", ".ts", ".tsx"}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Warning: cannot read %s: %v", path, err)
		return
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(path, fileName)

		if file.IsDir() {
			wg.Add(1)
			go ScanDirectory(filePath, extensionArgsSlice, wg)
			continue
		}

		fileExtension := filepath.Ext(fileName)
		if len(extensionArgsSlice) > 0 && !slices.Contains(extensionArgsSlice, fileExtension) {
			continue
		}
		if fileExtension == "" {
			continue
		}

		fileInf, err := file.Info()
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			continue
		}
		size := fileInf.Size()

		lineCount, err := CountLines(filePath)
		if err != nil {
			log.Printf("Warning: cannot read %s: %v", filePath, err)
			continue
		}

		isCompFile := false
		for _, suffix := range suffixSlice {
			if fileExtension == suffix && len(fileName) > 0 && fileName[0] >= 'A' && fileName[0] <= 'Z' {
				isCompFile = true
				break
			}
		}

		Mu.Lock()
		if _, exists := Stats[fileExtension]; !exists {
			Stats[fileExtension] = InitializeStats()
		}
		s := Stats[fileExtension]
		s.Count++
		s.TotalSize += size
		s.TotalLines += lineCount
		if isCompFile {
			s.TotalCompFiles++
		}
		s.FileInfo = append(s.FileInfo, FileInfo{
			Name:       fileName,
			Lines:      lineCount,
			IsCompFile: isCompFile,
		})
		Mu.Unlock()
	}
}
