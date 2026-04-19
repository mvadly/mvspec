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
		{"int32", "integer"},
		{"int16", "integer"},
		{"int8", "integer"},
		{"uint", "integer"},
		{"uint64", "integer"},
		{"bool", "boolean"},
		{"float64", "number"},
		{"float32", "number"},
		{"file", "string"},
		{"*os.File", "string"},
		{"array", "object"},
		{"FileHeader", "string"},
		{"multipart.FileHeader", "string"},
		{"object", "object"},
		{"SomethingFileHeader", "string"},
		{"my_multipart_Data", "string"},
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

func TestParseResponse(t *testing.T) {
	tests := []struct {
		input    string
		success bool
		expected *Response
	}{
		{
			input:    `200 {object} ResponseModel "success response"`,
			success: true,
			expected: &Response{
				Code:        200,
				Type:        "ResponseModel",
				Description: "success response",
			},
		},
		{
			input:    `400 {object} ErrorModel "error"`,
			success: false,
			expected: &Response{
				Code:        400,
				Type:        "ErrorModel",
				Description: "error",
			},
		},
		{
			input:    `200 {object} ResponseModel "success" request:{"key":"value"}`,
			success: true,
			expected: &Response{
				Code:           200,
				Type:           "ResponseModel",
				Description:    "success",
				RequestExample: `{"key":"value"}`,
			},
		},
		{
			input:    `404 NotFound "resource not found"`,
			success: false,
			expected: &Response{
				Code:        404,
				Type:        "NotFound",
				Description: "resource not found",
			},
		},
		{
			input:    "",
			success: true,
			expected: nil,
		},
		{
			input:    "only code",
			success: true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseResponse(tt.input, tt.success)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseResponse() = %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("parseResponse() = nil, want %v", tt.expected)
			}
			if result.Code != tt.expected.Code {
				t.Errorf("Code = %d, want %d", result.Code, tt.expected.Code)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %q, want %q", result.Type, tt.expected.Type)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", result.Description, tt.expected.Description)
			}
			if result.RequestExample != tt.expected.RequestExample {
				t.Errorf("RequestExample = %q, want %q", result.RequestExample, tt.expected.RequestExample)
			}
		})
	}
}

func TestExtractPathFromRouter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/users", "/api/users"},
		{"api/users", "/api/users"},
		{"no-slash", "/no-slash"},
		{"", ""},
		{"/v1", "/v1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractPathFromRouter(tt.input)
			if result != tt.expected {
				t.Errorf("extractPathFromRouter(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractMethodFromHandler(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserController.GetUser", "GetUser"},
		{"*UserService.Create", "Create"},
		{"simple", "simple"},
		{"", ""},
		{"Controller.ActionName", "ActionName"},
		{"*Service.DoSomething", "DoSomething"},
		{"a.b.c", "c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractMethodFromHandler(tt.input)
			if result != tt.expected {
				t.Errorf("extractMethodFromHandler(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s       string
		substr  string
		expected bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "xyz", false},
		{"", "", true},
		{"Hello", "", true},
		{"", "hello", false},
		{"Test123", "test", true},
		{"TEST", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := containsIgnoreCase(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestParseJSONExample(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"key":"value"}`, map[string]interface{}{"key": "value"}},
		{`["a","b"]`, []interface{}{"a", "b"}},
		{`"string"`, "string"},
		{`123`, 123},
		{`true`, true},
		{`null`, nil},
		{"invalid{json", nil},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJSONExample(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseJSONExample(%q) = %v, want nil", tt.input, result)
				}
				return
			}
			if result == nil {
				t.Errorf("parseJSONExample(%q) = nil, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestSplitTags(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"tag1,tag2,tag3", []string{"tag1", "tag2", "tag3"}},
		{"tag1", []string{"tag1"}},
		{"", nil},
		{"  tag1 ,  tag2  ", []string{"tag1", "tag2"}},
		{"tag1,,tag2", []string{"tag1", "tag2"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := splitTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitTags(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range tt.expected {
				if i >= len(result) {
					break
				}
				if result[i] != tt.expected[i] {
					t.Errorf("splitTags(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}