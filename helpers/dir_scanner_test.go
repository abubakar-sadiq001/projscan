package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 bytes"},
		{512, "512 bytes"},
		{1024, "1.00 kb"},
		{1500, "1.46 kb"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}
	for _, tt := range tests {
		got := formatFileSize(tt.bytes)
		if got != tt.want {
			t.Errorf("formatFileSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestReportPanic(t *testing.T) {
	// No panic – should do nothing.
	func() {
		defer ReportPanic()
		// Nothing
	}()

	// Panic with error – should print the error and continue.
	var out bytes.Buffer
	// Redirect stdout to capture output.
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	func() {
		defer ReportPanic()
		panic(fmt.Errorf("test error"))
	}()

	w.Close()
	os.Stdout = oldStdout
	io.Copy(&out, r)
	got := out.String()
	if !strings.Contains(got, "test error") {
		t.Errorf("ReportPanic did not print expected error; got %q", got)
	}

	// Panic with non-error – should re-panic.
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with non-error value, but did not panic")
		} else if r != "not an error" {
			t.Errorf("expected panic value 'not an error', got %v", r)
		}
	}()
	func() {
		defer ReportPanic()
		panic("not an error")
	}()
}

func TestScanDirectory(t *testing.T) {
	// Create a temporary directory with files.
	tmpDir := t.TempDir()

	// Create some files with various extensions.
	files := []struct {
		name    string
		content string
	}{
		{"foo.js", "line1\nline2\n"},
		{"Bar.tsx", "line1\nline2\nline3\n"},
		{"baz.ts", "one\n"},
		{"util.go", "package main\n"},
		{"Component.js", "line1\nline2\nline3\nline4\nline5\n"},
		{"sub/Another.tsx", "line1\nline2\n"},
		{"sub/ignore.txt", "hello\n"},
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f.name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Reset global Stats before test.
	Stats = make(map[string]*FileStats)
	processedCount = 0

	// Provide extension filter: only .js, .tsx, .ts (no .go or .txt).
	filter := []string{".js", ".tsx", ".ts"}
	ScanDirectory(tmpDir, filter)

	// Verify Stats.
	expectedExtensions := map[string]struct {
		count      int
		totalLines int
		compFiles  int
		fileNames  []string
	}{
		".js": {
			count:      2,
			totalLines: 2 + 5, // foo.js (2), Component.js (5)
			compFiles:  1,     // Component.js is component
			fileNames:  []string{"foo.js", "Component.js"},
		},
		".tsx": {
			count:      2,
			totalLines: 3 + 2, // Bar.tsx (3), sub/Another.tsx (2)
			compFiles:  2,     // both start with uppercase
			fileNames:  []string{"Bar.tsx", "Another.tsx"},
		},
		".ts": {
			count:      1,
			totalLines: 1, // baz.ts
			compFiles:  0, // starts with lowercase
			fileNames:  []string{"baz.ts"},
		},
	}

	if len(Stats) != len(expectedExtensions) {
		t.Errorf("expected %d extensions, got %d", len(expectedExtensions), len(Stats))
	}

	for ext, exp := range expectedExtensions {
		stat, ok := Stats[ext]
		if !ok {
			t.Errorf("extension %q not found in Stats", ext)
			continue
		}
		if stat.Count != exp.count {
			t.Errorf("%s: Count = %d, want %d", ext, stat.Count, exp.count)
		}
		if stat.TotalLines != exp.totalLines {
			t.Errorf("%s: TotalLines = %d, want %d", ext, stat.TotalLines, exp.totalLines)
		}
		if stat.TotalCompFiles != exp.compFiles {
			t.Errorf("%s: TotalCompFiles = %d, want %d", ext, stat.TotalCompFiles, exp.compFiles)
		}
		// Check file names (order may vary, so we sort)
		var names []string
		for _, fi := range stat.FileInfo {
			names = append(names, fi.Name)
		}
		// Sort both slices for comparison.
		sortStrings(names)
		sortStrings(exp.fileNames)
		if !stringSlicesEqual(names, exp.fileNames) {
			t.Errorf("%s: FileInfo names = %v, want %v", ext, names, exp.fileNames)
		}
	}
}

func TestPrintStats(t *testing.T) {
	// Populate Stats with controlled data.
	Stats = map[string]*FileStats{
		".js": {
			Count:          2,
			TotalSize:      1234,
			TotalLines:     100,
			TotalCompFiles: 1,
			FileInfo: []FileInfo{
				{Name: "foo.js", Lines: 40, IsCompFile: false},
				{Name: "Bar.js", Lines: 60, IsCompFile: true},
			},
		},
		".tsx": {
			Count:          1,
			TotalSize:      567,
			TotalLines:     50,
			TotalCompFiles: 1,
			FileInfo: []FileInfo{
				{Name: "Component.tsx", Lines: 50, IsCompFile: true},
			},
		},
	}

	// Also include a large component to trigger the warning.
	Stats[".tsx"].FileInfo = append(Stats[".tsx"].FileInfo, FileInfo{
		Name:       "Huge.tsx",
		Lines:      400,
		IsCompFile: true,
	})

	// Capture stdout.
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintStats()

	w.Close()
	os.Stdout = oldStdout
	var out bytes.Buffer
	io.Copy(&out, r)
	output := out.String()

	// Check for expected sections.
	checks := []string{
		"React Components Found: 2",
		"FILES PER EXTENSION",
		".js (2 files):",
		"1. foo.js                                40 lines",
		"2. Bar.js                                60 lines",
		".tsx (2 files):",
		"1. Component.tsx                         50 lines",
		"2. Huge.tsx                             400 lines",
		"FILE STATISTICS BY EXTENSION",
		"Extension       | Count    | Size         | Total Lines",
		".js             | 2        | 1.21 kb      | 100",
		".tsx            | 2        | 567 bytes    | 450", // 50+400
		"TOTAL           | 4        | 1.76 kb      | 550",
		"Large Components",
		"Huge.tsx  400 lines",
		"⚠️  Warning:",
		"1 component exceed 300 lines.",
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Errorf("PrintStats output missing expected substring: %q\nFull output:\n%s", want, output)
		}
	}
}

// Helper functions for slice comparison.
func sortStrings(s []string) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
