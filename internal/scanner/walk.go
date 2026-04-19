package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type WalkFunc func(string, string)

func WalkDir(dir string, excludeMap map[string]bool, fn WalkFunc) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)

		if excludeMap[name] || strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			WalkDir(path, excludeMap, fn)
			continue
		}

		fn(path, name)
	}
}