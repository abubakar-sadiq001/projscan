package helpers

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var expectedExtensions = map[string]struct {
	count      int
	totalSize  int64
	totalLines int
	fileInfo   []struct {
		name       string
		lines      int
		isCompFile bool
	}
	totalCompFiles int
}{
	".js": {
		count:          2,
		totalLines:     7, // foo.js (2), Component.js (5)
		totalCompFiles: 1, // Component.js is component
		fileInfo: []struct {
			name       string
			lines      int
			isCompFile bool
		}{
			{name: "foo.js", lines: 2, isCompFile: false},
			{name: "Component.js", lines: 5, isCompFile: true},
		},
	},
	".tsx": {
		count:          2,
		totalLines:     5, // Bar.tsx (3), sub/Another.tsx (2)
		totalCompFiles: 2, // both start with uppercase
		fileInfo: []struct {
			name       string
			lines      int
			isCompFile bool
		}{
			{name: "Bar.tsx", lines: 3, isCompFile: true},
			{name: "Another.tsx", lines: 2, isCompFile: true},
		},
	},
	".ts": {
		count:          1,
		totalLines:     1, // baz.ts
		totalCompFiles: 0, // starts with lowercase
		fileInfo: []struct {
			name       string
			lines      int
			isCompFile bool
		}{
			{name: "baz.ts", lines: 1, isCompFile: false},
		},
	},
}

var files = []struct{ name, content string }{
	{"foo.js", "line1\nline2\n"},
	{"Bar.tsx", "line1\nline2\nline3\n"},
	{"baz.ts", "one\n"},
	{"util.go", "package main\n"},
	{"Component.js", "line1\nline2\nline3\nline4\nline5\n"},
	{"sub/Another.tsx", "line1\nline2\n"},
	{"sub/ignore.txt", "hello\n"},
}

func Test_Scan_Directory(t *testing.T) {
	Stats = make(map[string]*FileStats)
	// Mu := sync.Mutex{}
	var wg sync.WaitGroup

	tmpDir := t.TempDir()

	for _, f := range files {
		p := filepath.Join(tmpDir, f.name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(f.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	filter := []string{}

	wg.Add(1)
	ScanDirectory(tmpDir, filter, &wg)
	wg.Wait()

	if len(Stats) != 5 {
		t.Fatalf("expected 5 extensions, got %d: %v", len(Stats), Stats)
	}

}

func Test_Scan_Directory_With_Filter(t *testing.T) {
	Stats = make(map[string]*FileStats)
	Mu = sync.Mutex{}
	var wg sync.WaitGroup

	tmpDir := t.TempDir()

	for _, f := range files {
		p := filepath.Join(tmpDir, f.name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(f.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	filter := []string{".js", ".tsx", ".ts"}

	wg.Add(1)
	ScanDirectory(tmpDir, filter, &wg)
	wg.Wait()

	if len(Stats) != 3 {
		t.Fatalf("expected 3 extensions, got %d: %v", len(Stats), Stats)
	}

	js := Stats[".js"]
	if js == nil {
		t.Fatal("missing .js")
	}
	if js.Count != 2 {
		t.Errorf(".js count = %d, want 2", js.Count)
	}
	if js.TotalLines != 7 {
		t.Errorf(".js lines = %d, want 7", js.TotalLines)
	}
	if js.TotalCompFiles != 1 {
		t.Errorf(".js comp = %d, want 1", js.TotalCompFiles)
	}

	tsx := Stats[".tsx"]
	if tsx == nil {
		t.Fatal("missing .tsx")
	}
	if tsx.Count != 2 {
		t.Errorf(".tsx count = %d, want 2", tsx.Count)
	}
	if tsx.TotalLines != 5 {
		t.Errorf(".tsx lines = %d, want 5", tsx.TotalLines)
	}
	if tsx.TotalCompFiles != 2 {
		t.Errorf(".tsx comp = %d, want 2", tsx.TotalCompFiles)
	}

	ts := Stats[".ts"]
	if ts == nil {
		t.Fatal("missing .ts")
	}
	if ts.Count != 1 {
		t.Errorf(".ts count = %d, want 1", ts.Count)
	}
	if ts.TotalLines != 1 {
		t.Errorf(".ts lines = %d, want 1", ts.TotalLines)
	}
	if ts.TotalCompFiles != 0 {
		t.Errorf(".ts comp = %d, want 0", ts.TotalCompFiles)
	}
}
