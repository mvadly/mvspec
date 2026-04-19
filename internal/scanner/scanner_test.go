package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	exclude := []string{"vendor", ".git"}
	s := New(exclude)

	if s.Exclude == nil {
		t.Error("Exclude is nil")
	}
	if s.Files == nil {
		t.Error("Files is nil")
	}
	if len(s.Exclude) != 2 {
		t.Errorf("Exclude length = %d, want 2", len(s.Exclude))
	}

	emptySlice := []string{}
	s2 := New(emptySlice)
	if s2.Exclude == nil {
		t.Error("Exclude is nil for empty slice")
	}
}

func TestFileScannerScan(t *testing.T) {
	tests := []struct {
		name          string
		setupDir      func(t *testing.T) string
		exclude      []string
		wantFileCount int
		wantErr      bool
		checkPaths   []string
	}{
		{
			name: "scans directory",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				os.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte("package main"), 0644)
				os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("package main"), 0644)
				return tmpDir
			},
			exclude:      []string{},
			wantFileCount: 2,
			wantErr:      false,
			checkPaths: []string{},
		},
		{
			name: "empty directory",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			exclude:      []string{},
			wantFileCount: 0,
			wantErr:      false,
			checkPaths:  nil,
		},
		{
			name: "with excludes exact match",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				os.WriteFile(filepath.Join(tmpDir, "file.go"), []byte("package main"), 0644)
				os.MkdirAll(filepath.Join(tmpDir, "vendor"), 0755)
				os.WriteFile(filepath.Join(tmpDir, "vendor/file.go"), []byte("package vendor"), 0644)
				return tmpDir
			},
			exclude:      []string{"vendor"},
			wantFileCount: 1,
			wantErr:      false,
			checkPaths:  []string{},
		},
		{
			name: "nonexistent directory",
			setupDir: func(t *testing.T) string {
				return "/nonexistent/path"
			},
			exclude:      []string{},
			wantFileCount: 0,
			wantErr:      false,
			checkPaths:  nil,
		},
		{
			name: "nested directories",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "sub")
				os.MkdirAll(subDir, 0755)
				os.WriteFile(filepath.Join(tmpDir, "root.go"), []byte("package main"), 0644)
				os.WriteFile(filepath.Join(subDir, "nested.go"), []byte("package main"), 0644)
				return tmpDir
			},
			exclude:      []string{},
			wantFileCount: 2,
			wantErr:      false,
			checkPaths:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			s := New(tt.exclude)

			err := s.Scan(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(s.Files) != tt.wantFileCount {
				t.Errorf("Files count = %d, want %d", len(s.Files), tt.wantFileCount)
			}

			if tt.checkPaths != nil && len(tt.checkPaths) > 0 {
				for _, path := range tt.checkPaths {
					if s.Files[path] == nil {
						t.Errorf("File %q not found in map", path)
					}
				}
			}
		})
	}
}