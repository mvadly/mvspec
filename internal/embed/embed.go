package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	specOnce sync.Once
	specData []byte
	specPath string
)

type Config struct {
	DevOnly   bool
	Watch     bool
	ServePath string
}

var DefaultConfig = Config{
	DevOnly:   true,
	Watch:     true,
	ServePath: "/mvdocs",
}

const swaggerIndexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>API Documentation - {{.Title}}</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.0/swagger-ui.css">
  <style>
    html { box-sizing: border-box; }
    *, *:before, *:after { box-sizing: inherit; }
    #swagger-ui { max-width: 1460px; margin: 0 auto; padding: 20px; }
    .topbar { display: none !important; }
    @media screen and (max-width: 1200px) {
      #swagger-ui { padding: 10px; }
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.0/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.0/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const params = new URLSearchParams(window.location.search);
      const url = params.get('url') || window.configUrl || '/mv-spec.json';
      
      SwaggerUI({
        url: url,
        dom_id: '#swagger-ui',
        deepLinking: true,
        docExpansion: 'list',
        filter: true,
        showExtensions: true,
        showCommonExtensions: true,
        presets: [
          SwaggerUI.presets.apis,
          SwaggerUIStandalonePresets
        ],
        layout: 'StandaloneLayout'
      });
    };
  </script>
</body>
</html>`

var swaggerIndexTemplate = template.Must(template.New("swagger").Parse(swaggerIndexHTML))

type IndexData struct {
	Title   string
	Version string
}

func MvHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isDevEnvironment(r) {
			http.NotFound(w, r)
			return
		}
		serveSwaggerUI(w, r)
	})
}

func MvHandlerWithConfig(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.DevOnly && !isDevEnvironment(r) {
			http.NotFound(w, r)
			return
		}
		serveSwaggerUI(w, r)
	})
}

func isDevEnvironment(r *http.Request) bool {
	if os.Getenv("MVSPEC_DEV_ONLY") == "true" {
		return true
	}
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		return true
	}
	if env == "local" {
		return true
	}
	if os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development" {
		return true
	}
	return r.RemoteAddr == "127.0.0.1" || r.RemoteAddr == "::1" || strings.HasPrefix(r.RemoteAddr, "127.")
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimSuffix(path, "/")

	if path == "" || path == "/" || path == "/mvdocs" || path == "/mvdocs/" {
		serveIndex(w, r)
		return
	}

	if path == "/mv-spec.json" || path == "/mvdocs/mv-spec.json" {
		serveSpec(w, r)
		return
	}

	http.NotFound(w, r)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	title := "API Documentation"
	version := "1.0"

	if data, err := getSpecData(); err == nil {
		var spec map[string]interface{}
		if json.Unmarshal(data, &spec) == nil {
			if info, ok := spec["info"].(map[string]interface{}); ok {
				if t, ok := info["title"].(string); ok {
					title = t
				}
				if v, ok := info["version"].(string); ok {
					version = v
				}
			}
		}
	}

	data := IndexData{
		Title:   title,
		Version: version,
	}

	var buf bytes.Buffer
	if err := swaggerIndexTemplate.Execute(&buf, data); err != nil {
		fmt.Fprintf(w, "<html><body><h1>Error rendering template: %v</h1></body></html>", err)
		return
	}

	w.Write(buf.Bytes())
}

func serveSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	data, err := getSpecData()
	if err != nil {
		http.Error(w, "Spec not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}

func getSpecData() ([]byte, error) {
	specOnce.Do(func() {
		paths := []string{
			"mv-spec.json",
			"./mv-spec.json",
			"../mv-spec.json",
		}

		for _, p := range paths {
			if data, err := os.ReadFile(p); err == nil {
				specData = data
				specPath = p
				return
			}
		}

		specData = []byte(`{"openapi":"3.0.0","info":{"title":"API","version":"1.0"},"paths":{}}`)
	})

	if len(specData) == 0 {
		return nil, fmt.Errorf("spec not found")
	}
	return specData, nil
}

func ReloadSpec() {
	specOnce = sync.Once{}
	specData = nil
	specPath = ""
}

func WatchSpec(filepath string, fn func()) error {
	if filepath == "" {
		filepath = "mv-spec.json"
	}

	watcher, err := os.Open(filepath)
	if err != nil {
		return err
	}
	watcher.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ticker.C:
			if _, err := os.Stat(filepath); err == nil {
				ReloadSpec()
				if fn != nil {
					fn()
				}
			}
		}
	}
	return nil
}

func WatchSpecAsync(filepath string, fn func()) {
	go WatchSpec(filepath, fn)
}

func Serve(addr string) error {
	return ServeWith(addr, MvHandler())
}

func ServeWith(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return server.ListenAndServe()
}

func OpenBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = runCmd("cmd", "/c", "start", url)
	case "darwin":
		err = runCmd("open", url)
	case "linux":
		err = runCmd("xdg-open", url)
	default:
		err = runCmd("xdg-open", url)
	}
	return err
}

func runCmd(name string, arg ...string) error {
	cmd := &exec.Cmd{
		Path: name,
		Args: arg,
	}
	return cmd.Start()
}
