# mvspec

API Specification Generator - Generate OpenAPI 3.x specs from your Go/JS code with swaggo-compatible annotations.

## Installation

```bash
go install github.com/mvadly/mvspec@latest
```

Or build from source:

```bash
git clone https://github.com/mvadly/mvspec
cd mvspec
go build -o mvspec ./cmd
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
mvspec fmt         # Format annotations in code
mvspec validate   # Validate without generating
mvspec embed      # Generate embedded docs handler and UI
```

## Configuration

Create `mvspec.yaml` in your project root:

```yaml
title: My API
version: 1.0
description: API description
host: api.example.com
basePath: /v1
output: mv-spec.json
exclude:
  - ./internal
  - ./vendor
parseTypes: true
```

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

### Global annotations (in main.go or entry file)

```go
// @title           API Title
// @version         1.0
// @description     API description
// @host           api.example.com
// @basePath       /v1
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

### Integration

```go
// In your main.go or router setup
import "github.com/mvadly/mvspec/mv-docs"

// Add route
r.GET("/mvdocs", gin.WrapF(mvdocs.MvHandler()))
```

The UI includes:
- Collections panel (organized API endpoints)
- Request builder (method, URL, headers, body)
- Response viewer (status, time, body with syntax highlighting)
- History (recent requests)
- Environment variables support
- Green glassmorphism theme (#10B981)

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

## License

MIT License - see [LICENSE](LICENSE)