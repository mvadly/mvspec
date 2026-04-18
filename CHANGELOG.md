# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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