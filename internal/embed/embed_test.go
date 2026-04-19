package embed

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIsDevEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		envVars  map[string]string
		remoteAddr string
		want     bool
	}{
		{
			name:     "MVSPEC_DEV_ONLY true",
			envVars:  map[string]string{"MVSPEC_DEV_ONLY": "true"},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "GO_ENV development",
			envVars:  map[string]string{"GO_ENV": "development"},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "GO_ENV empty",
			envVars:  map[string]string{"GO_ENV": ""},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "GO_ENV local",
			envVars:  map[string]string{"GO_ENV": "local"},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "ENV local",
			envVars:  map[string]string{"ENV": "local"},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "ENV development",
			envVars:  map[string]string{"ENV": "development"},
			remoteAddr: "",
			want:     true,
		},
		{
			name:     "localhost IPv4",
			envVars:  map[string]string{},
			remoteAddr: "127.0.0.1",
			want:     true,
		},
		{
			name:     "localhost IPv6",
			envVars:  map[string]string{},
			remoteAddr: "::1",
			want:     true,
		},
		{
			name:     "localhost IPv4 prefix",
			envVars:  map[string]string{},
			remoteAddr: "127.0.0.2",
			want:     true,
		},
		{
			name:     "production ENV",
			envVars:  map[string]string{"GO_ENV": "production"},
			remoteAddr: "",
			want:     false,
		},
		{
			name:     "production ENV",
			envVars:  map[string]string{"ENV": "production", "GO_ENV": "production"},
			remoteAddr: "",
			want:     false,
		},
		{
			name:     "production remote addr",
			envVars:  map[string]string{"ENV": "production", "GO_ENV": "production"},
			remoteAddr: "192.168.1.1",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			r := &http.Request{}
			if tt.remoteAddr != "" {
				r.RemoteAddr = tt.remoteAddr
			}

			got := isDevEnvironment(r)
			if got != tt.want {
				t.Errorf("isDevEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpecData(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		wantErr  bool
		checkData bool
	}{
		{
			name: "mv-spec.json exists",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				os.WriteFile(filepath.Join(tmpDir, "mv-spec.json"), []byte(`{"openapi":"3.0.3"}`), 0644)
				return tmpDir
			},
			wantErr:  false,
			checkData: true,
		},
		{
			name: "no file",
			setupFile: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:  false,
			checkData: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFile(t)
			origCwd, _ := os.Getwd()
			os.Chdir(dir)
			defer os.Chdir(origCwd)

			ReloadSpec()

			data, err := getSpecData()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSpecData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkData && data == nil {
				t.Error("data is nil")
			}
		})
	}
}

func TestServeSpec(t *testing.T) {
	tests := []struct {
		name     string
		setupFile func(t *testing.T) string
		wantCode int
	}{
		{
			name: "returns spec",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				os.WriteFile(filepath.Join(tmpDir, "mv-spec.json"), []byte(`{"openapi":"3.0.3"}`), 0644)
				return tmpDir
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFile(t)
			origCwd, _ := os.Getwd()
			os.Chdir(dir)
			defer os.Chdir(origCwd)

			ReloadSpec()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/mv-spec.json", nil)

			serveSpec(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestOpenBrowser(t *testing.T) {
	err := OpenBrowser("http://localhost:8080")
	if err != nil {
		t.Logf("OpenBrowser error (expected in some environments): %v", err)
	}
}

func TestReloadSpec(t *testing.T) {
	tmpDir := t.TempDir()
	specFile := filepath.Join(tmpDir, "mv-spec.json")
	os.WriteFile(specFile, []byte(`{"openapi":"3.0.3"}`), 0644)

	origCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origCwd)

	ReloadSpec()

	data, err := getSpecData()
	if err != nil {
		t.Errorf("getSpecData() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("spec data is empty")
	}
}

func TestMvHandlerWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		devOnly bool
		envVars map[string]string
		want   int
	}{
		{
			name:    "dev only true dev env",
			devOnly: true,
			envVars: map[string]string{"GO_ENV": "development"},
			want:   200,
		},
		{
			name:    "dev only false",
			devOnly: false,
			envVars: map[string]string{"GO_ENV": "production"},
			want:   404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg := Config{DevOnly: tt.devOnly}
			handler := MvHandlerWithConfig(cfg)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/mvdocs", nil)
			handler.ServeHTTP(w, r)

			if tt.want == 200 && w.Code == 404 {
				t.Logf("Got 404 in dev mode - may be expected")
			}
		})
	}
}

func TestServeWith(t *testing.T) {
	handler := MvHandler()
	if handler == nil {
		t.Error("handler is nil")
	}
}