# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.4] - 2026-05-31

### Added
- **Multiple Server URLs**: Config and UI dropdown for switching between environments
- **Environment Variables**: Pre-defined variables in config, UI modal with `{{variable}}` substitution
- **Reset to Defaults**: Button in env modal to restore config defaults
- **Sidebar Toggle**: Show/hide sidebar with localStorage persistence
- **Response Headers Tab**: View response headers separately
- **Request Headers in Result**: Show sent request headers alongside response
- **Form-Data Support**: File upload with multipart/form-data, individual form fields
- **Content-Type Auto-Detection**: UI detects and shows appropriate body editor (json, form, form-data)
- **File Type Detection**: `*multipart.FileHeader` correctly detected as file (format: binary)
- **Form Struct Tags**: Parser uses `form` tags instead of field names for schema properties
- **Combined Examples**: Request + Response examples displayed together with Try button
- **API Header**: Title and description displayed in UI
- **Responsive Design**: Mobile breakpoints and improved layout handling

### Fixed
- Form-data parsing in annotations and UI display
- Form tag and file type detection for form-data schemas
- Separate file input for each file field in multipart form-data
- Response/request example parsing in annotations (request: prefix support)
- JavaScript duplicate variable declaration

### Changed
- `mvspec init` generates new config format with servers and env variables
- Major UI restructuring for better layout control
- Updated README with new features and examples

## [v1.0.0] - 2026-04-18

### Added
- **Postman-like UI**: Side-by-side request and response panels by default
- **Layout Toggle**: SVG icons to switch between horizontal (side-by-side) and vertical (stacked) layout
- **Request/Response Examples**: Both examples displayed in Examples tab with Try button
- **Fixed Request Bar**: Request bar stays at top, only panels toggle with layout

### Fixed
- Toggle icons now show only one at a time based on layout state
- Syntax error in examples HTML generation

### Changed
- Updated README with new UI features
- Restructured embed.go HTML/CSS for better layout control

---

## [v0.0.0] - Previous versions

Earlier versions focused on core OpenAPI spec generation from Go/JS annotations.

### Features (prior versions)
- OpenAPI 3.x spec generation
- swaggo-compatible annotations
- Request body support
- Response examples
- Multiple framework support (gin, echo, fiber, chi, mux)
- Embedded docs handler generation