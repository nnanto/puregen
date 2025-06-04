package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nnanto/puregen/generator"
	"github.com/nnanto/puregen/idl"
	"github.com/spf13/cobra"
)

var version = "dev" // This will be set during build time

var (
	inputFile             string
	templatePaths         string
	outputDir             string
	validateFile          string
	additionalContextJSON string
)

var rootCmd = &cobra.Command{
	Use:   "puregen",
	Short: "Generate code from IDL files using templates",
	Long:  `Puregen is a code generator that uses YAML IDL files and templates to generate code.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from IDL files using templates",
	Long:  `Generate code from YAML IDL files using one or more templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			fmt.Fprintf(os.Stderr, "Error: input file is required\n")
			os.Exit(1)
		}
		if templatePaths == "" {
			fmt.Fprintf(os.Stderr, "Error: template paths are required\n")
			os.Exit(1)
		}

		parser := idl.NewParser()
		schema, err := parser.ParseFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing IDL file: %v\n", err)
			os.Exit(1)
		}

		// Parse additional context JSON if provided
		if additionalContextJSON != "" {
			var additionalContext map[string]interface{}
			if err := json.Unmarshal([]byte(additionalContextJSON), &additionalContext); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing additional context JSON: %v\n", err)
				os.Exit(1)
			}
			schema.AdditionalContext = additionalContext
		}

		gen := generator.New()

		// Process each template
		for _, templatePath := range strings.Split(templatePaths, ",") {
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
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("puregen version %s\n", version)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate YAML IDL file",
	Long:  `Validate a YAML IDL file and check for potential issues with field types.`,
	Run: func(cmd *cobra.Command, args []string) {
		if validateFile == "" {
			fmt.Fprintf(os.Stderr, "Error: input file is required\n")
			os.Exit(1)
		}

		parser := idl.NewParser()
		schema, err := parser.ParseFile(validateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing IDL file: %v\n", err)
			os.Exit(1)
		}

		// Define primitive types
		primitiveTypes := map[string]bool{
			"string":  true,
			"int":     true,
			"int32":   true,
			"int64":   true,
			"float":   true,
			"float32": true,
			"float64": true,
			"bool":    true,
			"byte":    true,
			"bytes":   true,
		}

		// Collect all defined message names
		definedMessages := make(map[string]bool)
		for _, message := range schema.Messages {
			definedMessages[message.Name] = true
		}

		// Validate field types
		hasWarnings := false
		for _, message := range schema.Messages {
			for fieldName, field := range message.Fields {
				fieldType := field.Type
				// Remove array notation if present
				fieldType = strings.TrimPrefix(fieldType, "[]")

				// Check if type is primitive or defined in Messages
				if !primitiveTypes[fieldType] && !definedMessages[fieldType] {
					fmt.Printf("Warning: Field '%s' in message '%s' has type '%s' which is not primitive and not defined in Messages\n",
						fieldName, message.Name, field.Type)
					hasWarnings = true
				}
			}
		}

		if !hasWarnings {
			fmt.Println("Validation completed successfully - no issues found")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(creatorCmd)

	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input YAML IDL file (required)")
	generateCmd.Flags().StringVarP(&templatePaths, "templates", "t", "", "Template file paths (comma-separated for multiple templates) (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "generated", "Output directory for generated files")
	generateCmd.Flags().StringVar(&additionalContextJSON, "additional-context-json", "", "Additional context as JSON to pass to the template")

	generateCmd.MarkFlagRequired("input")
	generateCmd.MarkFlagRequired("templates")

	validateCmd.Flags().StringVarP(&validateFile, "input", "i", "", "Input YAML IDL file to validate (required)")
	validateCmd.MarkFlagRequired("input")

	creatorCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Output YAML file path (required)")
	creatorCmd.MarkFlagRequired("output-file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
