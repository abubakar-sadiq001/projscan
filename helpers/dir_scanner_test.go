package helpers

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	filter := []string{}

	wg.Add(1)
	ScanDirectory(tmpDir, filter, &wg)
	wg.Wait()

	assert.Len(t, Stats, 5)
	assert.NotNil(t, Stats)
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

	js := Stats[".js"]
	tsx := Stats[".tsx"]
	ts := Stats[".ts"]

	assert.Len(t, Stats, 3)

	assert.NotNil(t, js)
	assert.Equal(t, js.Count, 2)
	assert.Equal(t, js.TotalLines, 7)
	assert.Equal(t, js.TotalCompFiles, 1)

	assert.NotNil(t, tsx)
	assert.Equal(t, tsx.Count, 2)
	assert.Equal(t, tsx.TotalLines, 5)
	assert.Equal(t, tsx.TotalCompFiles, 2)

	assert.NotNil(t, ts)
	assert.Equal(t, ts.Count, 1)
	assert.Equal(t, ts.TotalLines, 1)
	assert.Equal(t, ts.TotalCompFiles, 0)
}
