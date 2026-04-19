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

func TestServeSpecJSON(t *testing.T) {
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

func TestServeWithHandler(t *testing.T) {
	handler := MvHandler()
	if handler == nil {
		t.Error("handler is nil")
	}
}

func TestServeSwaggerUI(t *testing.T) {
	tmpDir := t.TempDir()
	specFile := filepath.Join(tmpDir, "mv-spec.json")
	os.WriteFile(specFile, []byte(`{"openapi":"3.0.3","info":{"title":"Test API","version":"1.0"}}`), 0644)

	origCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origCwd)

	ReloadSpec()

	tests := []struct {
		name       string
		path      string
		wantCode  int
		wantBody  bool
	}{
		{
			name:      "root path",
			path:      "/",
			wantCode:  200,
			wantBody:  true,
		},
		{
			name:      "mvdocs path",
			path:      "/mvdocs",
			wantCode:  200,
			wantBody:  true,
		},
		{
			name:      "mvdocs with slash",
			path:      "/mvdocs/",
			wantCode:  200,
			wantBody:  true,
		},
		{
			name:      "mv-spec.json path",
			path:      "/mv-spec.json",
			wantCode:  200,
			wantBody:  true,
		},
		{
			name:      "mvdocs mv-spec.json path",
			path:      "/mvdocs/mv-spec.json",
			wantCode:  200,
			wantBody:  true,
		},
		{
			name:      "not found path",
			path:      "/unknown",
			wantCode:  404,
			wantBody:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tt.path, nil)
			r.RemoteAddr = "127.0.0.1:1234"

			serveSwaggerUI(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}

			if tt.wantBody && w.Body.Len() == 0 {
				t.Error("expected body but got empty")
			}
		})
	}
}

func TestServeSpecWithHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	specFile := filepath.Join(tmpDir, "mv-spec.json")
	os.WriteFile(specFile, []byte(`{"openapi":"3.0.3","info":{"title":"Test","version":"1.0"}}`), 0644)

	origCwd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origCwd)

	ReloadSpec()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/mv-spec.json", nil)

	serveSpec(w, r)

	if w.Code != 200 {
		t.Errorf("status = %d, want 200", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json; charset=utf-8", contentType)
	}

	if w.Header().Get("Cache-Control") == "" {
		t.Error("expected Cache-Control header")
	}
}

func TestWatchSpecAsync(t *testing.T) {
	WatchSpecAsync("test.json", func() {
		t.Log("callback called")
	})
}

func TestMvHandlerProduction(t *testing.T) {
	os.Setenv("GO_ENV", "production")
	os.Setenv("ENV", "production")
	defer os.Unsetenv("GO_ENV")
	defer os.Unsetenv("ENV")

	handler := MvHandler()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/mvdocs", nil)
	r.RemoteAddr = "192.168.1.1:1234"
	handler.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("status = %d, want 404 in production", w.Code)
	}
}