package scanner

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

	WalkDir(dir, excludeMap, func(path, name string) {
		s.Files[path] = &FileInfo{
			Path: path,
			Name: name,
		}
	})
	return nil
}