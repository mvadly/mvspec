package util

import (
	"testing"
)

func TestCleanRef(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"response.User", "User"},
		{"models.Response", "Response"},
		{"util.Error", "Error"},
		{"response.models.User", "User"},
		{"plain", "plain"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CleanRef(tt.input)
			if result != tt.expected {
				t.Errorf("CleanRef(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s       string
		substr  string
		want    bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "world", true},
		{"Hello World", "xyz", false},
		{"", "", true},
		{"", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := ContainsIgnoreCase(tt.s, tt.substr)
			if result != tt.want {
				t.Errorf("ContainsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.want)
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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExtractPathFromRouter(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractPathFromRouter(%q) = %q, want %q", tt.input, result, tt.expected)
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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExtractMethodFromHandler(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractMethodFromHandler(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseJSONExample(t *testing.T) {
	tests := []struct {
		input    string
		wantNil  bool
	}{
		{`{"key":"value"}`, false},
		{`123`, false},
		{`"test"`, false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseJSONExample(tt.input)
			isNil := result == nil
			if isNil != tt.wantNil {
				t.Errorf("ParseJSONExample(%q) = %v, wantNil %v", tt.input, result, tt.wantNil)
			}
		})
	}
}