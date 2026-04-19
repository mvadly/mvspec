package goparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/mvadly/mvspec/internal/config"
	"github.com/mvadly/mvspec/internal/generator"
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
		{`123`, 123.0},
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
		})
	}
}

func TestGenerateIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(goFile, []byte(`
package main

// @Summary Get users
// @Description Get all users
// @Tags User
// @Accept json
// @Produce json
// @Param body formData UserForm true "User form"
// @Success 200 {object} UserResponse
// @Router /users [get]
func GetUsers() type

type UserForm struct {
	Name string `+"`form:\"name\"`"+`
	Age  int    `+"`form:\"age\"`"+`
}

type UserResponse struct {
	Code string
	Data string
}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg := &config.Config{
		Title:       "Test API",
		Version:     "1.0.0",
		Output:      "test_spec.json",
		ParseTypes:  true,
		Exclude:    []string{"vendor"},
		Servers:    []config.ServerConfig{{URL: "http://localhost:8080"}},
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if _, err := os.Stat(cfg.Output); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	os.Remove(cfg.Output)
}

func TestGenerateNoParseTypes(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(goFile, []byte(`
package main

// @Summary Get users
// @Router /users [get]
func GetUsers() {}
`), 0644)

	cfg := &config.Config{
		Title:      "Test API",
		Version:    "1.0.0",
		Output:    "test_spec2.json",
		ParseTypes: false,
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	os.Remove(cfg.Output)
}

func TestGenerateSpec(t *testing.T) {
	p := &Parser{
		cfg:    &config.Config{Title: "Test", Version: "1.0"},
		fset:   *token.NewFileSet(),
		files:  make(map[string]*ast.File),
		routes: []Route{{Path: "/test", Method: "GET", Handler: "GetTest"}},
		types: make(map[string]*TypeInfo),
		operas: []Operation{
			{
				Handler: "GetTest",
				Annotation: &Annotation{
					Summary:     "Get test",
					Description: "Get test endpoint",
					Tags:        []string{"Test"},
				},
			},
		},
	}

	spec := p.generateSpec()
	if spec == nil {
		t.Fatal("generateSpec() returned nil")
	}
	if spec.OpenAPI != "3.0.3" {
		t.Errorf("OpenAPI = %q, want 3.0.3", spec.OpenAPI)
	}
	if spec.Info.Title != "Test" {
		t.Errorf("Info.Title = %q, want Test", spec.Info.Title)
	}
	if spec.Paths == nil {
		t.Error("Paths is nil")
	}
}

func TestWriteGenerator(t *testing.T) {
	spec := &generator.OpenAPISpec{
		OpenAPI: "3.0.3",
		Info:   generator.Info{Title: "Test", Version: "1.0"},
		Paths:  map[string]map[string]generator.Operation{},
	}

	tmpFile := filepath.Join(t.TempDir(), "spec.json")
	err := generator.Write(tmpFile, spec)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	os.Remove(tmpFile)
}

func TestTypeToString(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"string type", "string"},
		{"int type", "int"},
		{"pointer to string", "*string"},
		{"slice of string", "[]string"},
		{"map type", "map[string]int"},
		{"interface type", "interface{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", "package main\ntype T "+tt.code, 0)
			if err != nil {
				t.Skipf("Skipping: cannot parse code")
			}

			var typeSpec *ast.TypeSpec
			for _, d := range f.Decls {
				if gd, ok := d.(*ast.GenDecl); ok && len(gd.Specs) > 0 {
					if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
						typeSpec = ts
						break
					}
				}
			}
			if typeSpec == nil {
				t.Skipf("No type spec found")
			}

			result := typeToString(typeSpec.Type)
			if len(result) == 0 {
				t.Error("expected non-empty result")
			}
		})
	}
}

func TestGenerateFullSpec(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(goFile, []byte(`
package main

// @Summary Get users
// @Description Get all users
// @Tags User
// @Accept json
// @Produce json
// @Param body formData LoginForm true "Login form"
// @Success 200 {object} Response
// @Router /users [get]
func GetUsers() type

type LoginForm struct {
	Username string `+"`form:\"username\"`"+`
	Password string `+"`form:\"password\"`"+`
}

type Response struct {
	Code string
	Data string
}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg := &config.Config{
		Title:       "Full Test API",
		Version:     "2.0.0",
		Output:      "full_test.json",
		ParseTypes:  true,
		Exclude:    []string{},
		Servers:    []config.ServerConfig{{URL: "http://localhost:8080"}},
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if _, err := os.Stat(cfg.Output); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	os.Remove(cfg.Output)
}

func TestGenerateWithHost(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(goFile, []byte(`package main`), 0644)

	cfg := &config.Config{
		Title:       "Host Test API",
		Version:     "1.0.0",
		Output:      "host_test.json",
		Host:        "localhost:8080",
		BasePath:    "/api",
		ParseTypes:  false,
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	os.Remove(cfg.Output)
}

func TestGenerateWithSecurity(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(goFile, []byte(`
package main

// @Summary Get users
// @Security ApiKeyAuth
// @Router /users [get]
func GetUsers() {}

// @Summary Create user
// @Security OAuth2Implicit
// @Router /users [post]
func CreateUser() {}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg := &config.Config{
		Title:       "Security Test API",
		Version:     "1.0.0",
		Output:      "security_test.json",
		ParseTypes:  false,
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	os.Remove(cfg.Output)
}

func TestParseRoutesWithGinRouter(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "router.go")
	err := os.WriteFile(goFile, []byte(`
package main

import "github.com/gin-gonic/gin"

func GetUsers() {}
func CreateUser() {}
func UpdateUser() {}
func DeleteUser() {}
func GetItems() {}

func main() {
	r := gin.Default()
	r.GET("/users", GetUsers)
	r.GET("/items", GetItems)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	api := r.Group("/api")
	api.GET("/users", GetUsers)
}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	p := &Parser{
		cfg:    &config.Config{Exclude: []string{}},
		fset:   *token.NewFileSet(),
		files:  make(map[string]*ast.File),
		routes: []Route{},
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
	if err != nil {
		t.Skipf("Skipping: cannot parse file: %v", err)
	}
	p.files[goFile] = f

	err = p.parseRoutes()
	if err != nil {
		t.Errorf("parseRoutes() error = %v", err)
	}

	t.Logf("Found %d routes", len(p.routes))
}

func TestParseAnnotationsIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "handlers.go")
	err := os.WriteFile(goFile, []byte(`
package main

// @Summary Get all users
// @Description Returns users list
// @Tags User,Admin
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Security ApiKeyAuth
// @Router /users [get]
func GetUsers() {}

// @Summary Create user
// @Tags User
// @Param body body CreateUserRequest true "User"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser() {}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	p := &Parser{
		cfg:    &config.Config{Exclude: []string{}},
		fset:   *token.NewFileSet(),
		files:  make(map[string]*ast.File),
		operas: []Operation{},
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
	if err != nil {
		t.Skipf("Skipping: cannot parse file: %v", err)
	}
	p.files[goFile] = f

	err = p.parseAnnotations()
	if err != nil {
		t.Errorf("parseAnnotations() error = %v", err)
	}

	t.Logf("Found %d operations", len(p.operas))
}

func TestParseTypesIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "types.go")
	err := os.WriteFile(goFile, []byte(`
package main

import "mime/multipart"

type User struct {
	ID     int    `+"`json:\"id\"`"+`
	Name   string `+"`json:\"name\"`"+`
	Email  string `+"`json:\"email\"`"+`
	Avatar *multipart.FileHeader `+"`form:\"avatar\"`"+`
	Tags   []string `+"`json:\"tags\"`"+`
}

type CreateUserRequest struct {
	Name string `+"`form:\"name\"`"+`
}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	p := &Parser{
		cfg:    &config.Config{Exclude: []string{}},
		fset:   *token.NewFileSet(),
		files:  make(map[string]*ast.File),
		types: make(map[string]*TypeInfo),
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
	if err != nil {
		t.Skipf("Skipping: cannot parse file: %v", err)
	}
	p.files[goFile] = f

	err = p.parseTypes()
	if err != nil {
		t.Errorf("parseTypes() error = %v", err)
	}

	t.Logf("Found %d types", len(p.types))
}

func TestGenerateFullAPI(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(goFile, []byte(`
package main

import "mime/multipart"

// @Summary Get users
// @Tags User
// @Success 200 {object} UserList
// @Router /users [get]
func GetUsers() {}

// @Summary Create user
// @Tags User
// @Param body formData CreateUserRequest true "User"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser() {}

type User struct {
	ID   int    `+"`json:\"id\"`"+`
	Name string `+"`json:\"name\"`"+`
}

type CreateUserRequest struct {
	Name string `+"`form:\"name\"`"+`
}

type UserList struct {
	Data []User `+"`json:\"data\"`"+`
}
`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg := &config.Config{
		Title:       "Test API",
		Version:     "1.0.0",
		Output:     "test.json",
		ParseTypes: true,
		Servers:    []config.ServerConfig{{URL: "http://localhost:8080"}},
	}

	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	os.Remove(cfg.Output)
}