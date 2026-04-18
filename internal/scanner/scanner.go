package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type FileScanner struct {
	Exclude []string
	Files   map[string]*FileInfo
}

type FileInfo struct {
	Path string
	Name string
}

func New(exclude []string) *FileScanner {
	return &FileScanner{
		Exclude: exclude,
		Files:   make(map[string]*FileInfo),
	}
}

func (s *FileScanner) Scan(dir string) error {
	excludeMap := make(map[string]bool)
	for _, e := range s.Exclude {
		excludeMap[e] = true
	}

	var walkDir func(dir string) error
	walkDir = func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		for _, entry := range entries {
			name := entry.Name()
			path := filepath.Join(dir, name)

			if excludeMap[name] || strings.HasPrefix(name, ".") {
				continue
			}

			if entry.IsDir() {
				walkDir(path)
				continue
			}

			s.Files[path] = &FileInfo{
				Path: path,
				Name: name,
			}
		}
		return nil
	}

	return walkDir(dir)
}
