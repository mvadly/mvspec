package generator

import (
	"encoding/json"
	"fmt"
	"os"
)

type OpenAPISpec struct {
	OpenAPI    string                          `json:"openapi"`
	Info       Info                            `json:"info"`
	Servers    []Server                        `json:"servers,omitempty"`
	Paths      map[string]map[string]Operation `json:"paths"`
	Components Components                      `json:"components,omitempty"`
}

type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

type Schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Ref        string            `json:"$ref,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Example    interface{}       `json:"example,omitempty"`
	Enum       []interface{}     `json:"enum,omitempty"`
}

type SecurityScheme struct {
	Type         string            `json:"type"`
	Scheme       string            `json:"scheme,omitempty"`
	BearerFormat string            `json:"bearerFormat,omitempty"`
	Flow         string            `json:"flow,omitempty"`
	TokenUrl     string            `json:"tokenUrl,omitempty"`
	Scopes       map[string]string `json:"scopes,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationId string                `json:"operationId,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Consumes   []string            `json:"consumes,omitempty"`
	Produces  []string            `json:"produces,omitempty"`
	Responses  map[string]Response  `json:"responses,omitempty"`
	Security   []map[string][]string `json:"security,omitempty"`
	Deprecated bool                  `json:"deprecated,omitempty"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Schema      Schema `json:"schema,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema  Schema      `json:"schema,omitempty"`
	Example interface{} `json:"example,omitempty"`
}

type Response struct {
	Description   string               `json:"description"`
	Headers       map[string]Parameter `json:"headers,omitempty"`
	Content       map[string]MediaType `json:"content,omitempty"`
	Example       interface{}          `json:"example,omitempty"`
	RequestExample interface{}         `json:"requestExample,omitempty"`
}

func Write(filename string, spec *OpenAPISpec) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}
