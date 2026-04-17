package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mvadly/mvspec/internal/config"
	goparser "github.com/mvadly/mvspec/internal/parser/go"
	jsparser "github.com/mvadly/mvspec/internal/parser/js"
)

var (
	lang       string
	output     string
	exclude    string
	parseTypes bool
	cfgFile    string
)

func main() {
	flag.StringVar(&lang, "lang", "auto", "Language: go, js, auto")
	flag.StringVar(&output, "output", "mv-spec.json", "Output file")
	flag.StringVar(&exclude, "exclude", "", "Directories to exclude (comma-separated)")
	flag.BoolVar(&parseTypes, "parseTypes", true, "Enable type inference")
	flag.StringVar(&cfgFile, "config", "mvspec.yaml", "Config file")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "init":
			if err := runInit(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "fmt":
			if err := runFmt(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "validate":
			if err := runValidate(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "embed":
			if err := runEmbed(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	if err := runGenerate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate() error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	if output != "mv-spec.json" {
		cfg.Output = output
	}
	if exclude != "" {
		cfg.Exclude = splitCSV(exclude)
	}
	cfg.ParseTypes = parseTypes

	detectLang := lang
	if detectLang == "auto" {
		detectLang = detectLanguage()
	}

	switch detectLang {
	case "go":
		return goparser.Generate(cfg)
	case "js":
		return jsparser.Generate(cfg)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}

func detectLanguage() string {
	if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	}
	if _, err := os.Stat("package.json"); err == nil {
		return "js"
	}
	return "go"
}

func runInit() error {
	yml := `title: My API
version: 1.0
description: API description
host: api.example.com
basePath: /v1
output: mv-spec.json
exclude:
  - ./internal
  - ./vendor
parseTypes: true
`
	return os.WriteFile(cfgFile, []byte(yml), 0644)
}

func runFmt() error {
	fmt.Println("Formatting annotations...")
	return nil
}

func runValidate() error {
	fmt.Println("Validating annotations...")
	return nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range strings.Split(s, ",") {
		result = append(result, strings.TrimSpace(part))
	}
	return result
}
