package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nnanto/puregen/generator"
	"github.com/nnanto/puregen/idl"
	"github.com/spf13/cobra"
)

var version = "dev" // This will be set during build time

var (
	inputFile     string
	templatePaths string
	outputDir     string
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

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)

	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input YAML IDL file (required)")
	generateCmd.Flags().StringVarP(&templatePaths, "templates", "t", "", "Template file paths (comma-separated for multiple templates) (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "generated", "Output directory for generated files")

	generateCmd.MarkFlagRequired("input")
	generateCmd.MarkFlagRequired("templates")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
