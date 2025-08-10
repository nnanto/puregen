package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/nnanto/puregen/internal/generator"
)

const version = "1.0.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-puregen %v\n", version)
		return
	}

	var flags flag.FlagSet
	languageFlag := flags.String("language", "all", "target language: go, java, python, or all")
	commonNamespaceFlag := flags.String("common_namespace", "", "namespace for common classes/interfaces (e.g., 'common' or 'shared.transport')")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		language := *languageFlag
		commonNamespace := *commonNamespaceFlag

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			switch language {
			case "go":
				generator.GenerateGoFile(gen, f, commonNamespace)
			case "java":
				generator.GenerateJavaFile(gen, f, commonNamespace)
			case "python":
				generator.GeneratePythonFile(gen, f, commonNamespace)
			case "all":
				generator.GenerateGoFile(gen, f, commonNamespace)
				generator.GenerateJavaFile(gen, f, commonNamespace)
				generator.GeneratePythonFile(gen, f, commonNamespace)
			default:
				return fmt.Errorf("unsupported language: %s", language)
			}
		}
		return nil
	})
}
