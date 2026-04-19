# mvspec

API Specification Generator - Generate OpenAPI 3.x specs from your Go/JS code with swaggo-compatible annotations.

## Installation

```bash
go install github.com/mvadly/mvspec/cmd/mvspec@latest
```

Or build from source:

```bash
git clone https://github.com/mvadly/mvspec
cd mvspec
go build -o mvspec ./cmd/mvspec
```

## Usage

### Generate API spec

```bash
# Auto-detect language and generate
mvspec

# Force Go mode
mvspec --lang go

# Force JavaScript/TypeScript mode
mvspec --lang js
```

### Create config template

```bash
mvspec init
```

This creates `mvspec.yaml` in current directory.

### Other commands

```bash
mvspec version    # Print version
mvspec fmt         # Format annotations in code
mvspec validate   # Validate without generating
mvspec embed      # Generate embedded docs handler and UI
mvspec remove     # Remove generated files (mv-spec.json, mv-docs/)
```

## Configuration

Create `mvspec.yaml` in your project root:

```yaml
title: My API
version: 1.0
description: API description
output: mv-spec.json
exclude:
  - ./internal
  - ./vendor
parseTypes: true
servers:
  - url: http://localhost:8080
    description: Local
  - url: https://api.example.com
    description: Production
env:
  - name: API_KEY
    value: ""
    description: "API Key for authentication"
  - name: TOKEN
    value: ""
    description: "Bearer token for API access"
```

### Servers (Multiple Base URLs)

Define multiple server URLs in your config. The API docs UI will show a dropdown to switch between them.

- **1 server**: Used automatically as the base URL
- **2+ servers**: Dropdown appears in the request bar, users can switch between environments

```yaml
servers:
  - url: http://localhost:8080
    description: Local Development
  - url: https://dev.api.example.com
    description: Development
  - url: https://qa.api.example.com
    description: QA
  - url: https://staging.api.example.com
    description: Staging
  - url: https://api.example.com
    description: Production
```

### Environment Variables

Pre-define environment variables that will be available in the API docs UI. Users can edit values and use them in requests with `{{VARIABLE_NAME}}` syntax.

```yaml
env:
  - name: API_KEY
    value: ""
    description: "API Key for authentication"
  - name: TOKEN
    value: ""
    description: "Bearer token"
  - name: BASE_URL
    value: "https://api.example.com"
    description: "Base URL for API calls"
```

The UI allows users to:
- View and edit environment variables
- Reset to default values from config
- Use variables in URLs, headers, and body with `{{VARIABLE_NAME}}` syntax

## Supported Annotations

Uses swaggo-compatible format:

```go
// @Summary     Get user by ID
// @Description Retrieves a user from the system
// @Tags        users
// @Accept      json
// @Produce    json
// @Param      id path int true "User ID"
// @Param      include query string false "Fields to include"
// @Success    200 {object} UserResponse "Success"
// @Failure    400 {object} Error "Bad request"
// @Failure    404 {object} Error "Not found"
// @Router     /users/{id} [get]
```

### Request Body

`@Param` with `body` generates an OpenAPI 3.x `requestBody`:

```go
// @Param request body CreateUserRequest true "User data"
```

This produces:
```json
"requestBody": {
  "description": "User data",
  "required": true,
  "content": {
    "application/json": {
      "schema": { "$ref": "#/components/schemas/CreateUserRequest" }
    }
  }
}
```

### Form Data (multipart/form-data)

`@Param` with `formData` generates `multipart/form-data` content-type:

```go
// @Param body formData UploadRequest true "Upload form"
```

This produces:
```json
"requestBody": {
  "description": "Upload form",
  "required": true,
  "content": {
    "multipart/form-data": {
      "schema": { "$ref": "#/components/schemas/UploadRequest" }
    }
  }
}
```

For individual form fields:

```go
// @Param name formData string true "Name"
// @Param file formData file true "File to upload"
```

This produces formData parameters in the operation.

### Response Examples

Add inline JSON examples to `@Success` and `@Failure` annotations:

```go
// @Success    200 {object} UserResponse "Success" {"id":1,"name":"John","email":"john@example.com"}
// @Failure    400 {object} Error "Bad request" {"code":"99","message":"Invalid input"}
```

### Request + Response Examples

Add both request and response examples for each status code using `request:{...}`:

```go
// @Success    200 {object} Response "Login successful" {"responseCode":"00","responseMessage":"Success"} request:{"username":"john","password":"123"}
// @Failure    400 {object} Response "Invalid credentials" {"responseCode":"01","responseMessage":"Wrong password"} request:{"username":"john","password":"wrong"}
// @Failure    500 {object} Response "Server error" {"responseCode":"99","responseMessage":"Error"} request:{"username":"john","password":"123"}
```

**Format:** `status_code {object} Type "description" {response_json} request:{request_json}`

### Global annotations (in main.go or entry file)

```go
// @title           API Title
// @version         1.0
// @description     API description
```

## Supported Frameworks

### Go
- gin-gonic/gin
- labstack/echo
- gofiber/fiber
- go-chi/chi
- gorilla/mux

### JavaScript/TypeScript (coming soon)
- Express
- Fastify
- NestJS

## Output

Generates `mv-spec.json` (OpenAPI 3.0.3) that can be imported into:
- Postman
- Swagger UI
- Redoc
- Any OpenAPI tool

## Embeddable API Docs UI

Generate a self-hosted API documentation with a custom Postman-like UI:

```bash
mvspec embed
```

This creates:
- `mv-docs/docs.go` - Go handler for serving the docs
- `mv-docs/index.html` - Custom API testing UI with green glass theme
- `mv-docs/styles.css` - Styling
- `mv-docs/app.js` - UI logic

### UI Features

- **Postman-like layout**: Side-by-side request and response panels by default
- **Layout toggle**: Switch between side-by-side and stacked vertical layout
- **Sidebar toggle**: Show/hide the sidebar with collections and history
- **Responsive design**: Works on all screen sizes
- Collections panel (organized API endpoints by tags)
- Request builder (method, URL, headers, body, auth)
- Response viewer (status, time, body, headers with syntax highlighting)
- History (recent requests)
- Environment variables support with reset to defaults
- Green glassmorphism theme (#10B981)

### Integration

```go
// In your main.go or router setup
import "your-module/mv-docs"

// Add route
r.GET("/mvdocs/*path", gin.WrapH(mvdocs.MvHandler()))
```

Access at `http://localhost:8080/mvdocs`

### Environment Variables

- `MVSPEC_DEV_ONLY=true` - Enable only in development mode
- `GO_ENV=development|local` - Auto-detected as dev mode

## Example

```go
// @Summary     Get user by ID
// @Description Retrieves a user from the system by ID
// @Tags        users
// @Accept      json
// @Produce    json
// @Param       id path int true "User ID"
// @Success     200 {object} UserResponse
// @Failure     404 {object} ErrorResponse
// @Router      /users/{id} [get]
func GetUserById(c *gin.Context) {
    // handler code here
}

// @Summary     Create user
// @Description Creates a new user
// @Tags        users
// @Accept      json
// @Produce    json
// @Param       request body CreateUserRequest true "User data"
// @Success     201 {object} UserResponse
// @Failure     400 {object} ErrorResponse
// @Router      /users [post]
func CreateUser(c *gin.Context) {
    // handler code here
}
```

Run:

```bash
mvspec
```

Output `mv-spec.json`:

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "My API",
    "version": "1.0"
  },
  "paths": {
    "/users/{id}": {
      "get": {
        "summary": "Get user by ID",
        "description": "Retrieves a user from the system by ID",
        "tags": ["users"],
        "parameters": [...],
        "responses": {...}
      }
    }
  }
}
```

## Recent Changes

### v1.x Updates

- **Form Data Support**: Parser now correctly handles `@Param body formData` and `@Param name formData` annotations to generate `multipart/form-data` content-type
- **Request/Response Examples**: Inline JSON examples in annotations are parsed and displayed in the UI
- **Content-Type Auto-Detection**: UI automatically detects content-type from endpoint spec and shows appropriate body editor (JSON, form-urlencoded, form-data)
- **Form Data File Upload**: UI shows file input for multipart/form-data endpoints

### UI Features

- Postman-like interface with glassmorphism theme
- Environment variables support
- Multiple server URLs support
- Request/Response examples display
- Form editor for form-urlencoded and multipart/form-data

## License

MIT License - see [LICENSE](LICENSE)