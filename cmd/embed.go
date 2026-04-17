package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var outputPath string

func runEmbed() error {
	flag.StringVar(&outputPath, "o", "./mv-docs", "Output directory")
	flag.Parse()

	dir := outputPath
	if dir == "" {
		dir = "./mv-docs"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if _, err := os.Stat("mv-spec.json"); err != nil {
		return fmt.Errorf("run mvspec first")
	}

	content := getDocsContent()
	path := filepath.Join(dir, "docs.go")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", path)
	fmt.Printf("\nUsage: r.GET(\"/mvdocs\", gin.WrapF(mvdocs.MvHandler()))\n")
	fmt.Printf("Access: http://localhost:8080/mvdocs\n")

	return nil
}

func getDocsContent() string {
	return "// Package mvdocs provides embedded API documentation handler.\n" +
		"// GENERATED FILE - DO NOT EDIT\n" +
		"// Run 'mvspec embed' to regenerate\n" +
		"\n" +
		"package mvdocs\n" +
		"\n" +
		"import (\n" +
		"\t\"net/http\"\n" +
		"\t\"os\"\n" +
		"\t\"strings\"\n" +
		"\t\"sync\"\n" +
		")\n" +
		"\n" +
		"var specOnce sync.Once\n" +
		"var specData []byte\n" +
		"\n" +
		"// MvHandler returns HTTP handler for API documentation.\n" +
		"func MvHandler() http.Handler {\n" +
		"\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n" +
		"\t\tif !isDev(r) {\n" +
		"\t\t\thttp.NotFound(w, r)\n" +
		"\t\t\treturn\n" +
		"\t\t}\n" +
		"\t\tserve(w, r)\n" +
		"\t})\n" +
		"}\n" +
		"\n" +
		"func isDev(r *http.Request) bool {\n" +
		"\tif os.Getenv(\"MVSPEC_DEV_ONLY\") == \"true\" {\n" +
		"\t\treturn true\n" +
		"\t}\n" +
		"\tenv := os.Getenv(\"GO_ENV\")\n" +
		"\tif env == \"\" || env == \"development\" || env == \"local\" {\n" +
		"\t\treturn true\n" +
		"\t}\n" +
		"\treturn strings.HasPrefix(r.RemoteAddr, \"127.\")\n" +
		"}\n" +
		"\n" +
		"func serve(w http.ResponseWriter, r *http.Request) {\n" +
		"\tpath := strings.TrimSuffix(r.URL.Path, \"/\")\n" +
		"\n" +
		"\tfilePath := strings.TrimPrefix(path, \"/mvdocs/\")\n" +
		"\n" +
		"\tif filePath == \"\" || filePath == \"index.html\" || path == \"/mvdocs\" {\n" +
		"\t\tw.Header().Set(\"Content-Type\", \"text/html\")\n" +
		"\t\tserveIndexHTML(w)\n" +
		"\t\treturn\n" +
		"\t}\n" +
		"\n" +
		"\tif filePath == \"mv-spec.json\" {\n" +
		"\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n" +
		"\t\tserveSpec(w)\n" +
		"\t\treturn\n" +
		"\t}\n" +
		"\n" +
		"\t// Try to read from mv-docs folder\n" +
		"\tstaticExt := map[string]string{\n" +
		"\t\t\"styles.css\": \"text/css\",\n" +
		"\t\t\"app.js\": \"application/javascript\",\n" +
		"\t\t\"index.html\": \"text/html\",\n" +
		"\t}\n" +
		"\n" +
		"\tfor filename, contentType := range staticExt {\n" +
		"\t\tif filePath == filename {\n" +
		"\t\t\tdata, err := os.ReadFile(\"mv-docs/\" + filename)\n" +
		"\t\t\tif err == nil {\n" +
		"\t\t\t\tw.Header().Set(\"Content-Type\", contentType)\n" +
		"\t\t\t\tw.Write(data)\n" +
		"\t\t\t\treturn\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\n" +
		"\t// Default to index\n" +
		"\tw.Header().Set(\"Content-Type\", \"text/html\")\n" +
		"\tserveIndexHTML(w)\n" +
		"}\n" +
		"\n" +
		"func serveIndexHTML(w http.ResponseWriter) {\n" +
		"\tdata, _ := os.ReadFile(\"mv-docs/index.html\")\n" +
		"\tif len(data) == 0 {\n" +
		"\t\tdata = []byte(\"<html><body><h1>MVAPI Docs</h1><p>Run mvspec embed first</p></body></html>\")\n" +
		"\t}\n" +
		"\tw.Write(data)\n" +
		"}\n" +
		"\n" +
		"func serveSpec(w http.ResponseWriter) {\n" +
		"\tdata, _ := os.ReadFile(\"mv-spec.json\")\n" +
		"\tif len(data) == 0 {\n" +
		"\t\tdata = []byte(\"{\\\"openapi\\\":\\\"3.0.0\\\",\\\"info\\\":{\\\"title\\\":\\\"API\\\"},\\\"paths\\\":{}}\")\n" +
		"\t}\n" +
		"\tw.Write(data)\n" +
		"}\n" +
		"\n" +
		"func ReloadSpec() {\n" +
		"\tspecOnce = sync.Once{}\n" +
		"}\n"
}
