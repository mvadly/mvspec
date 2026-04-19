package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name      string
		spec     *OpenAPISpec
		wantErr  bool
		checkFile bool
		checkFunc func(*testing.T, *OpenAPISpec)
	}{
		{
			name: "valid spec",
			spec: &OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: Info{
					Title:   "Test API",
					Version: "1.0",
				},
				Paths: map[string]map[string]Operation{},
			},
			wantErr:   false,
			checkFile: true,
			checkFunc: func(t *testing.T, spec *OpenAPISpec) {
				if spec.OpenAPI != "3.0.3" {
					t.Errorf("OpenAPI = %q, want %q", spec.OpenAPI, "3.0.3")
				}
			},
		},
		{
			name: "with components",
			spec: &OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: Info{
					Title:   "Test API",
					Version: "1.0",
				},
				Paths: map[string]map[string]Operation{},
				Components: Components{
					Schemas: map[string]Schema{
						"User": {
							Type: "object",
							Properties: map[string]Schema{
								"name": {Type: "string"},
							},
						},
					},
				},
			},
			wantErr:   false,
			checkFile: true,
		},
		{
			name: "empty spec",
			spec: &OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: Info{
					Title:   "",
					Version: "",
				},
				Paths: map[string]map[string]Operation{},
			},
			wantErr:   false,
			checkFile: true,
		},
		{
			name: "with servers",
			spec: &OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: Info{Title: "Test", Version: "1.0"},
				Paths: map[string]map[string]Operation{},
				Servers: []Server{
					{URL: "http://localhost:8080", Description: "Local"},
				},
			},
			wantErr:   false,
			checkFile: true,
		},
		{
			name: "with security schemes",
			spec: &OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: Info{Title: "Test", Version: "1.0"},
				Paths: map[string]map[string]Operation{},
				Components: Components{
					SecuritySchemes: map[string]SecurityScheme{
						"bearerAuth": {
							Type:         "http",
							Scheme:       "bearer",
							BearerFormat: "JWT",
						},
					},
				},
			},
			wantErr:   false,
			checkFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "spec.json")

			err := Write(filePath, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFile && err == nil {
				data, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("ReadFile() error = %v", err)
					return
				}
				if len(data) == 0 {
					t.Error("File is empty")
				}
			}

			if tt.checkFunc != nil && tt.spec != nil {
				tt.checkFunc(t, tt.spec)
			}
		})
	}
}