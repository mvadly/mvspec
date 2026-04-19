package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		wantErr    bool
		checkFunc  func(*testing.T, *Config)
	}{
		{
			name:       "valid config",
			configFile:  "title: Test API\nversion: 1.0\noutput: spec.json",
			wantErr:    false,
			checkFunc: func(t *testing.T, cfg *Config) {
				if cfg.Title != "Test API" {
					t.Errorf("Title = %q, want %q", cfg.Title, "Test API")
				}
				if cfg.Version != "1.0" {
					t.Errorf("Version = %q, want %q", cfg.Version, "1.0")
				}
				if cfg.Output != "spec.json" {
					t.Errorf("Output = %q, want %q", cfg.Output, "spec.json")
				}
			},
		},
		{
			name:       "file not found",
			configFile:  "",
			wantErr:    true,
			checkFunc:  nil,
		},
		{
			name:       "invalid yaml",
			configFile:  "title: [invalid yaml",
			wantErr:    true,
			checkFunc:  nil,
		},
		{
			name:       "default output",
			configFile:  "title: Test",
			wantErr:    false,
			checkFunc: func(t *testing.T, cfg *Config) {
				if cfg.Output != "mv-spec.json" {
					t.Errorf("Output = %q, want %q", cfg.Output, "mv-spec.json")
				}
			},
		},
		{
			name:       "with servers",
			configFile:  "title: Test\nservers:\n  - url: http://localhost:8080\n    description: Local",
			wantErr:    false,
			checkFunc: func(t *testing.T, cfg *Config) {
				if len(cfg.Servers) != 1 {
					t.Errorf("Servers length = %d, want 1", len(cfg.Servers))
				}
				if cfg.Servers[0].URL != "http://localhost:8080" {
					t.Errorf("Server URL = %q, want %q", cfg.Servers[0].URL, "http://localhost:8080")
				}
			},
		},
		{
			name:       "with exclude",
			configFile:  "title: Test\nexclude:\n  - vendor\n  - .git",
			wantErr:    false,
			checkFunc: func(t *testing.T, cfg *Config) {
				if len(cfg.Exclude) != 2 {
					t.Errorf("Exclude length = %d, want 2", len(cfg.Exclude))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.configFile == "" {
				_, err := Load("/nonexistent/path/config.yaml")
				if (err != nil) != tt.wantErr {
					t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.configFile), 0644); err != nil {
				t.Fatalf("WriteFile error: %v", err)
			}

			cfg, err := Load(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFunc != nil && cfg != nil {
				tt.checkFunc(t, cfg)
			}
		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantErr   bool
		checkFile bool
	}{
		{
			name: "valid save",
			config: &Config{
				Title:   "Test API",
				Version: "1.0",
			},
			wantErr:   false,
			checkFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			err := Save(configPath, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFile && err == nil {
				if _, err := os.ReadFile(configPath); err != nil {
					t.Errorf("ReadFile() error = %v", err)
				}
			}
		})
	}
}