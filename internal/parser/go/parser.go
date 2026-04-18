package goparser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/mvadly/mvspec/internal/config"
	"github.com/mvadly/mvspec/internal/generator"
)

type Parser struct {
	cfg    *config.Config
	fset   token.FileSet
	files  map[string]*ast.File
	routes []Route
	types  map[string]*TypeInfo
	operas []Operation
}

type Route struct {
	Path    string
	Method  string
	Handler string
}

type TypeInfo struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name    string
	Type    string
	JSONTag string
}

type Annotation struct {
	Summary     string
	Description string
	Tags        []string
	Accept      string
	Produce     string
	Params      []Param
	Success     []Response
	Failure     []Response
	Router      string
	Security    []string
}

type Param struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Description string
}

type Response struct {
	Code           int
	Type           string
	Description    string
	Example        string
	RequestExample string
}

type Operation struct {
	Path       string
	Method     string
	Handler    string
	Annotation *Annotation
}

func Generate(cfg *config.Config) error {
	p := &Parser{
		cfg:    cfg,
		files:  make(map[string]*ast.File),
		routes: []Route{},
		types:  make(map[string]*TypeInfo),
	}

	if err := p.scanFiles(); err != nil {
		return fmt.Errorf("scan files: %w", err)
	}

	if err := p.parseRoutes(); err != nil {
		return fmt.Errorf("parse routes: %w", err)
	}

	if err := p.parseAnnotations(); err != nil {
		return fmt.Errorf("parse annotations: %w", err)
	}

	if cfg.ParseTypes {
		if err := p.parseTypes(); err != nil {
			return fmt.Errorf("parse types: %w", err)
		}
	}

	spec := p.generateSpec()

	return generator.Write(cfg.Output, spec)
}

func (p *Parser) scanFiles() error {
	excludeMap := make(map[string]bool)
	for _, e := range p.cfg.Exclude {
		excludeMap[e] = true
	}

	var walkDir func(dir string) error
	walkDir = func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		for _, entry := range entries {
			name := entry.Name()
			path := dir + "/" + name

			if excludeMap[name] || strings.HasPrefix(name, ".") {
				continue
			}

			if entry.IsDir() {
				walkDir(path)
				continue
			}

			if !strings.HasSuffix(name, ".go") {
				continue
			}
			if strings.HasSuffix(name, "_test.go") {
				continue
			}

			f, err := parser.ParseFile(&p.fset, path, nil, parser.ParseComments)
			if err != nil {
				continue
			}
			p.files[path] = f
		}
		return nil
	}

	return walkDir(".")
}

func (p *Parser) parseRoutes() error {
	for _, f := range p.files {
		for _, decl := range f.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			var routeGrp string

			ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				switch sel.Sel.Name {
				case "Group":
					if len(call.Args) > 0 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							routeGrp = strings.Trim(lit.Value, `"`)
						}
					}
				case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD":
					if len(call.Args) >= 2 {
						pathArg := ""
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							pathArg = strings.Trim(lit.Value, `"`)
						}
						var handlerName string
						if ident, ok := call.Args[1].(*ast.Ident); ok {
							handlerName = ident.Name
						} else if sel2, ok := call.Args[1].(*ast.SelectorExpr); ok {
							handlerName = sel2.Sel.Name
						}

						fullPath := routeGrp + pathArg
						route := Route{
							Path:    fullPath,
							Method:  strings.ToUpper(sel.Sel.Name),
							Handler: handlerName,
						}
						p.routes = append(p.routes, route)
					}
				}
				return true
			})
		}
	}

	return nil
}

func (p *Parser) parseAnnotations() error {
	annoPattern := regexp.MustCompile(`@(\w+)\s+(.*)`)

	for _, f := range p.files {
		for _, decl := range f.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			handlerName := funcDecl.Name.Name
			if funcDecl.Doc == nil {
				continue
			}

			var anno Annotation

			for _, comment := range funcDecl.Doc.List {
				line := comment.Text[1:]
				line = strings.ReplaceAll(line, "\t", " ")

				parts := annoPattern.FindStringSubmatch(line)
				if len(parts) < 3 {
					continue
				}

				value := strings.TrimSpace(parts[2])
				switch parts[1] {
				case "Summary":
					anno.Summary = value
				case "Description":
					anno.Description = value
				case "Tags":
					anno.Tags = splitTags(value)
				case "Accept":
					anno.Accept = value
				case "Produce":
					anno.Produce = value
				case "Param":
					if param := parseParam(value); param != nil {
						anno.Params = append(anno.Params, *param)
					}
				case "Success":
					if resp := parseResponse(value, true); resp != nil {
						anno.Success = append(anno.Success, *resp)
					}
				case "Failure":
					if resp := parseResponse(value, false); resp != nil {
						anno.Failure = append(anno.Failure, *resp)
					}
				case "Router":
					anno.Router = value
				case "Security":
					anno.Security = append(anno.Security, value)
				}
			}

			if anno.Summary != "" || anno.Router != "" {
				p.operas = append(p.operas, Operation{
					Handler:    handlerName,
					Annotation: &anno,
				})
			}
		}
	}

	return nil
}

func (p *Parser) parseTypes() error {
	for _, f := range p.files {
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				info := &TypeInfo{Name: typeSpec.Name.Name}

				for _, field := range structType.Fields.List {
					if len(field.Names) == 0 {
						continue
					}

					fieldName := field.Names[0].Name
					typeName := typeToString(field.Type)
					jsonTag := strings.ToLower(fieldName)

					info.Fields = append(info.Fields, Field{
						Name:    fieldName,
						Type:    typeName,
						JSONTag: jsonTag,
					})
				}

				p.types[info.Name] = info
			}
		}
	}

	return nil
}

func (p *Parser) generateSpec() *generator.OpenAPISpec {
	spec := &generator.OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: generator.Info{
			Title:       p.cfg.Title,
			Version:     p.cfg.Version,
			Description: p.cfg.Description,
		},
		Paths: make(map[string]map[string]generator.Operation),
		Components: generator.Components{
			Schemas: make(map[string]generator.Schema),
		},
	}

	// Use servers from config if defined, otherwise fallback to Host+BasePath
	if len(p.cfg.Servers) > 0 {
		servers := make([]generator.Server, len(p.cfg.Servers))
		for i, s := range p.cfg.Servers {
			servers[i] = generator.Server{
				URL:         s.URL,
				Description: s.Description,
			}
		}
		spec.Servers = servers
	} else if p.cfg.Host != "" {
		spec.Servers = []generator.Server{
			{URL: p.cfg.Host + p.cfg.BasePath},
		}
	}

	for _, route := range p.routes {
		fullPath := route.Path
		if !strings.HasPrefix(fullPath, "/") {
			fullPath = "/" + fullPath
		}

		if spec.Paths[fullPath] == nil {
			spec.Paths[fullPath] = make(map[string]generator.Operation)
		}

		var op generator.Operation
		for _, o := range p.operas {
			matched := false
			var anno *Annotation

			if o.Handler == route.Handler {
				matched = true
				anno = o.Annotation
			}

			if !matched && o.Annotation.Router != "" {
				routerPath := extractPathFromRouter(o.Annotation.Router)
				if routerPath == fullPath {
					matched = true
					anno = o.Annotation
				}
			}

			if !matched {
				routeMethod := extractMethodFromHandler(route.Handler)
				if routeMethod != "" && containsIgnoreCase(o.Handler, routeMethod) {
					matched = true
					anno = o.Annotation
				}
			}

			if !matched && o.Annotation.Router != "" {
				routeMethod := extractMethodFromHandler(route.Handler)
				if routeMethod != "" && containsIgnoreCase(o.Handler, routeMethod) {
					routerPath := extractPathFromRouter(o.Annotation.Router)
					if routerPath == fullPath {
						matched = true
						anno = o.Annotation
					}
				}
			}

			if !matched && o.Annotation.Router != "" {
				routeMethod := extractMethodFromHandler(route.Handler)
				if routeMethod != "" && containsIgnoreCase(o.Handler, routeMethod) {
					routerPath := extractPathFromRouter(o.Annotation.Router)
					if routerPath == fullPath {
						matched = true
						anno = o.Annotation
					}
				}
			}

			if matched && anno != nil {
				op.Summary = anno.Summary
				op.Description = anno.Description
				op.Tags = anno.Tags

				for _, param := range anno.Params {
					if param.In == "body" || param.In == "formData" {
						schemaRef := param.Type
						if !strings.HasPrefix(schemaRef, "#/") {
							if strings.HasPrefix(schemaRef, "{") {
								schemaRef = strings.Trim(schemaRef, "{}")
							}
							schemaRef = "#/components/schemas/" + schemaRef
						}
						schemaRef = strings.ReplaceAll(schemaRef, "response.", "")
						schemaRef = strings.ReplaceAll(schemaRef, "models.", "")
						schemaRef = strings.ReplaceAll(schemaRef, "util.", "")
						op.RequestBody = &generator.RequestBody{
							Description: param.Description,
							Required:    param.Required,
							Content: map[string]generator.MediaType{
								"application/json": {
									Schema: generator.Schema{Ref: schemaRef},
								},
							},
						}
						continue
					}
					pm := generator.Parameter{
						Name:        param.Name,
						In:          param.In,
						Required:    param.Required,
						Description: param.Description,
						Schema:      generator.Schema{Type: param.Type},
					}
					op.Parameters = append(op.Parameters, pm)
				}

				for _, succ := range anno.Success {
					respRef := succ.Type
					respRef = strings.ReplaceAll(respRef, "response.", "")
					respRef = strings.ReplaceAll(respRef, "models.", "")
					respRef = strings.ReplaceAll(respRef, "util.", "")
					resp := generator.Response{
						Description: succ.Description,
						Content: map[string]generator.MediaType{
							"application/json": {
								Schema: generator.Schema{Ref: respRef},
							},
						},
						Example:        parseJSONExample(succ.Example),
						RequestExample: parseJSONExample(succ.RequestExample),
					}
					if op.Responses == nil {
						op.Responses = make(map[string]generator.Response)
					}
					op.Responses[fmt.Sprintf("%d", succ.Code)] = resp
				}

				for _, fail := range anno.Failure {
					respRef := fail.Type
					respRef = strings.ReplaceAll(respRef, "response.", "")
					respRef = strings.ReplaceAll(respRef, "models.", "")
					respRef = strings.ReplaceAll(respRef, "util.", "")
					resp := generator.Response{
						Description: fail.Description,
						Content: map[string]generator.MediaType{
							"application/json": {
								Schema: generator.Schema{Ref: respRef},
							},
						},
						Example:        parseJSONExample(fail.Example),
						RequestExample: parseJSONExample(fail.RequestExample),
					}
					if op.Responses == nil {
						op.Responses = make(map[string]generator.Response)
					}
					op.Responses[fmt.Sprintf("%d", fail.Code)] = resp
				}
				break
			}
		}

		if op.Summary == "" {
			if route.Handler != "" {
				op.Summary = route.Handler
			} else {
				routeMethod := extractMethodFromHandler(route.Handler)
				for _, o := range p.operas {
					if o.Annotation.Router != "" && routeMethod != "" && containsIgnoreCase(o.Handler, routeMethod) {
						routerParts := strings.Fields(o.Annotation.Router)
						if len(routerParts) >= 1 {
							routerPath := routerParts[0]
							if !strings.HasPrefix(routerPath, "/") {
								routerPath = "/" + routerPath
							}
							if routerPath == fullPath {
								anno := o.Annotation
								op.Summary = anno.Summary
								op.Description = anno.Description
								op.Tags = anno.Tags

								for _, param := range anno.Params {
									if param.In == "body" || param.In == "formData" {
										schemaRef := param.Type
										if !strings.HasPrefix(schemaRef, "#/") {
											if strings.HasPrefix(schemaRef, "{") {
												schemaRef = strings.Trim(schemaRef, "{}")
											}
											schemaRef = "#/components/schemas/" + schemaRef
										}
										schemaRef = strings.ReplaceAll(schemaRef, "response.", "")
										schemaRef = strings.ReplaceAll(schemaRef, "models.", "")
										schemaRef = strings.ReplaceAll(schemaRef, "util.", "")
										op.RequestBody = &generator.RequestBody{
											Description: param.Description,
											Required:    param.Required,
											Content: map[string]generator.MediaType{
												"application/json": {
													Schema: generator.Schema{Ref: schemaRef},
												},
											},
										}
										continue
									}
									pm := generator.Parameter{
										Name:        param.Name,
										In:          param.In,
										Required:    param.Required,
										Description: param.Description,
										Schema:      generator.Schema{Type: param.Type},
									}
									op.Parameters = append(op.Parameters, pm)
								}

								for _, succ := range anno.Success {
									respRef := succ.Type
									respRef = strings.ReplaceAll(respRef, "response.", "")
									respRef = strings.ReplaceAll(respRef, "models.", "")
									respRef = strings.ReplaceAll(respRef, "util.", "")
									resp := generator.Response{
										Description: succ.Description,
										Content: map[string]generator.MediaType{
											"application/json": {
												Schema: generator.Schema{Ref: respRef},
											},
										},
										Example:        parseJSONExample(succ.Example),
										RequestExample: parseJSONExample(succ.RequestExample),
									}
									if op.Responses == nil {
										op.Responses = make(map[string]generator.Response)
									}
									op.Responses[fmt.Sprintf("%d", succ.Code)] = resp
								}

								for _, fail := range anno.Failure {
									respRef := fail.Type
									respRef = strings.ReplaceAll(respRef, "response.", "")
									respRef = strings.ReplaceAll(respRef, "models.", "")
									respRef = strings.ReplaceAll(respRef, "util.", "")
									resp := generator.Response{
										Description: fail.Description,
										Content: map[string]generator.MediaType{
											"application/json": {
												Schema: generator.Schema{Ref: respRef},
											},
										},
										Example:        parseJSONExample(fail.Example),
										RequestExample: parseJSONExample(fail.RequestExample),
									}
									if op.Responses == nil {
										op.Responses = make(map[string]generator.Response)
									}
									op.Responses[fmt.Sprintf("%d", fail.Code)] = resp
								}
								break
							}
						}
					}
				}
			}
		}

		spec.Paths[fullPath][strings.ToLower(route.Method)] = op
	}

	for name, t := range p.types {
		schema := generator.Schema{
			Type:       "object",
			Properties: make(map[string]generator.Schema),
		}

		for _, f := range t.Fields {
			prop := generator.Schema{
				Type: inferJSONType(f.Type),
			}
			schema.Properties[f.JSONTag] = prop
		}

		spec.Components.Schemas[name] = schema
	}

	return spec
}

func splitTags(s string) []string {
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func parseParam(s string) *Param {
	parts := strings.Fields(s)
	if len(parts) < 5 {
		return nil
	}

	required := parts[3] == "true"
	return &Param{
		Name:        parts[0],
		In:          parts[1],
		Type:        parts[2],
		Required:    required,
		Description: strings.Join(parts[4:], " "),
	}
}

func parseResponse(s string, success bool) *Response {
	// Trim leading/trailing spaces
	s = strings.TrimSpace(s)

	// Use regex to properly parse the annotation
	// Format: status_code {object} type "description" {response} request:{request}
	re := regexp.MustCompile(`^(\d+)\s+\{object\}\s+(\S+)\s+"([^"]+)"(?:\s*(\{[^}]+\}))?(?:\s*(request:\{[^}]*\}))?`)
	matches := re.FindStringSubmatch(s)

	if len(matches) < 4 {
		// Fallback to old parsing
		parts := strings.Fields(s)
		if len(parts) < 3 {
			return nil
		}
		code := 200
		if !success {
			fmt.Sscanf(parts[0], "%d", &code)
		}
		return &Response{
			Code:        code,
			Type:        parts[1],
			Description: strings.Join(parts[2:], " "),
		}
	}

	code := 200
	if !success {
		fmt.Sscanf(matches[1], "%d", &code)
	}

	responseExample := ""
	if len(matches) >= 5 && matches[4] != "" {
		responseExample = matches[4]
	}

	requestExample := ""
	if len(matches) >= 6 && matches[5] != "" {
		reqMatch := regexp.MustCompile(`request:(\{[^}]+\})`).FindStringSubmatch(matches[5])
		if len(reqMatch) >= 2 {
			requestExample = reqMatch[1]
		}
	}

	return &Response{
		Code:           code,
		Type:           matches[2],
		Description:    matches[3],
		Example:        responseExample,
		RequestExample: requestExample,
	}
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "map"
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "object"
	}
}

func inferJSONType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64":
		return "integer"
	case "float64", "float32":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "object"
	}
}

func extractPathFromRouter(router string) string {
	parts := strings.Fields(router)
	if len(parts) >= 1 {
		path := parts[0]
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return path
	}
	return ""
}

func extractMethodFromHandler(handler string) string {
	if handler == "" {
		return ""
	}
	parts := strings.Split(handler, ".")
	if len(parts) > 0 {
		method := parts[len(parts)-1]
		method = strings.TrimPrefix(method, "*")
		return method
	}
	return ""
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func parseJSONExample(jsonStr string) interface{} {
	if jsonStr == "" {
		return nil
	}
	var result interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil
	}
	return result
}
