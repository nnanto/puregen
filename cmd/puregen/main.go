package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nnanto/puregen/generator"
	"github.com/nnanto/puregen/idl"
)

var version = "dev" // This will be set during build time

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	if showVersion {
		fmt.Printf("puregen version %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <yaml-file> <template-path(s)> [output-dir]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Template paths can be comma-separated for multiple templates\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := args[0]
	templatePaths := strings.Split(args[1], ",")

	outputDir := "generated"
	if len(args) >= 3 {
		outputDir = args[2]
	}

	parser := idl.NewParser()
	schema, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing IDL file: %v\n", err)
		os.Exit(1)
	}

	gen := generator.New()

	// Process each template
	for _, templatePath := range templatePaths {
		templatePath = strings.TrimSpace(templatePath)

		// Set up template reader from file
		templateFile, err := os.Open(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening template file %s: %v\n", templatePath, err)
			os.Exit(1)
		}

		err = gen.Generate(schema, templateFile, outputDir)
		templateFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating code from template %s: %v\n", templatePath, err)
			os.Exit(1)
		}
	}
}
