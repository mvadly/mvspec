package goparser

import (
	"testing"
)

func TestParseParamFormData(t *testing.T) {
	tests := []struct {
		input    string
		expected Param
	}{
		{
			input: "body formData models.SubmissionCluster true submission",
			expected: Param{
				Name:        "body",
				In:          "formData",
				Type:        "models.SubmissionCluster",
				Required:    true,
				Description: "submission",
			},
		},
		{
			input: "name formData string true Name field",
			expected: Param{
				Name:        "name",
				In:          "formData",
				Type:        "string",
				Required:    true,
				Description: "Name field",
			},
		},
		{
			input: "file formData file true File to upload",
			expected: Param{
				Name:        "file",
				In:          "formData",
				Type:        "file",
				Required:    true,
				Description: "File to upload",
			},
		},
		{
			input: "email formData string false Email",
			expected: Param{
				Name:        "email",
				In:          "formData",
				Type:        "string",
				Required:    false,
				Description: "Email",
			},
		},
		{
			input: "body formData ModelName true data",
			expected: Param{
				Name:        "body",
				In:          "formData",
				Type:        "ModelName",
				Required:    true,
				Description: "data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseParam(tt.input)
			if result == nil {
				t.Fatalf("parseParam returned nil for input: %s", tt.input)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.In != tt.expected.In {
				t.Errorf("In = %v, want %v", result.In, tt.expected.In)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
			if result.Required != tt.expected.Required {
				t.Errorf("Required = %v, want %v", result.Required, tt.expected.Required)
			}
		})
	}
}

func TestParseParamJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Param
	}{
		{
			input: "body body RequestBody true request body",
			expected: Param{
				Name:        "body",
				In:          "body",
				Type:        "RequestBody",
				Required:    true,
				Description: "request body",
			},
		},
		{
			input: "data body models.Data true JSON data",
			expected: Param{
				Name:        "data",
				In:          "body",
				Type:        "models.Data",
				Required:    true,
				Description: "JSON data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseParam(tt.input)
			if result == nil {
				t.Fatalf("parseParam returned nil for input: %s", tt.input)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.In != tt.expected.In {
				t.Errorf("In = %v, want %v", result.In, tt.expected.In)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
		})
	}
}

func TestParseParamQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected Param
	}{
		{
			input: "id query int true User ID",
			expected: Param{
				Name:        "id",
				In:          "query",
				Type:        "int",
				Required:    true,
				Description: "User ID",
			},
		},
		{
			input: "name query string false Filter by name",
			expected: Param{
				Name:        "name",
				In:          "query",
				Type:        "string",
				Required:    false,
				Description: "Filter by name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseParam(tt.input)
			if result == nil {
				t.Fatalf("parseParam returned nil for input: %s", tt.input)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.In != tt.expected.In {
				t.Errorf("In = %v, want %v", result.In, tt.expected.In)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
		})
	}
}

func TestParseParamHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected Param
	}{
		{
			input: "Authorization header string true Bearer token",
			expected: Param{
				Name:        "Authorization",
				In:          "header",
				Type:        "string",
				Required:    true,
				Description: "Bearer token",
			},
		},
		{
			input: "Content-Type header string false Content type",
			expected: Param{
				Name:        "Content-Type",
				In:          "header",
				Type:        "string",
				Required:    false,
				Description: "Content type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseParam(tt.input)
			if result == nil {
				t.Fatalf("parseParam returned nil for input: %s", tt.input)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.In != tt.expected.In {
				t.Errorf("In = %v, want %v", result.In, tt.expected.In)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
		})
	}
}

func TestInferJSONType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "string"},
		{"int", "integer"},
		{"int64", "integer"},
		{"bool", "boolean"},
		{"float64", "number"},
		{"file", "string"},
		{"*os.File", "string"},
		{"array", "object"},
		{"FileHeader", "string"},
		{"multipart.FileHeader", "string"},
		{"object", "object"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := inferJSONType(tt.input)
			if result != tt.expected {
				t.Errorf("inferJSONType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}