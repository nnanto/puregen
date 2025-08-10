package generator

import (
	"path/filepath"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// formatPythonComment formats a comment for Python code
func formatPythonComment(comments protogen.CommentSet) []string {
	var result []string

	// Use leading comments if available, otherwise trailing
	comment := comments.Leading
	if comment == "" && comments.Trailing != "" {
		comment = comments.Trailing
	}

	if comment == "" {
		return result
	}

	// Split by lines and format each line
	lines := strings.Split(strings.TrimSpace(string(comment)), "\n")

	// If single line, use # style
	if len(lines) == 1 {
		result = append(result, "# "+lines[0])
		return result
	}

	// Multi-line comment with """ """ style
	result = append(result, "\"\"\"")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			result = append(result, "")
		} else {
			result = append(result, line)
		}
	}
	result = append(result, "\"\"\"")

	return result
}

// writePythonComment writes formatted comments to the generator
func writePythonComment(g *protogen.GeneratedFile, comments protogen.CommentSet) {
	commentLines := formatPythonComment(comments)
	for _, line := range commentLines {
		g.P(line)
	}
}

// Track created package directories to avoid duplicates
var createdPythonPackages = make(map[string]bool)

// collectSamePackageMessagesPython collects messages from the same package that need to be imported
func collectSamePackageMessagesPython(file *protogen.File) []*protogen.Message {
	visited := make(map[string]bool)
	var samePackageMessages []*protogen.Message

	// Check all messages in the file for same package message references
	for _, message := range file.Messages {
		collectSamePackageFromMessagePython(message, file, visited, &samePackageMessages)
	}

	// Check service methods for same package message references
	for _, service := range file.Services {
		for _, method := range service.Methods {
			collectSamePackageFromMessagePython(method.Input, file, visited, &samePackageMessages)
			collectSamePackageFromMessagePython(method.Output, file, visited, &samePackageMessages)
		}
	}

	return samePackageMessages
}

// collectSamePackageFromMessagePython recursively collects same package messages from a message and its fields
func collectSamePackageFromMessagePython(msg *protogen.Message, currentFile *protogen.File, visited map[string]bool, samePackageMessages *[]*protogen.Message) {
	if msg == nil {
		return
	}

	messageKey := string(msg.Desc.FullName())

	// If this message is from the same package but different file and we haven't seen it before
	if !visited[messageKey] && isSamePackageMessagePython(msg, currentFile) {
		visited[messageKey] = true
		*samePackageMessages = append(*samePackageMessages, msg)
	}

	// Recursively check fields for same package message types
	for _, field := range msg.Fields {
		if field.Message != nil {
			collectSamePackageFromMessagePython(field.Message, currentFile, visited, samePackageMessages)
		}
	}

	// Check nested messages
	for _, nested := range msg.Messages {
		collectSamePackageFromMessagePython(nested, currentFile, visited, samePackageMessages)
	}
}

// isSamePackageMessagePython checks if a message is from the same package but different file
func isSamePackageMessagePython(msg *protogen.Message, currentFile *protogen.File) bool {
	if msg.Desc.ParentFile() == nil {
		return false
	}

	// Check if the message's package is the same as the current file's package
	msgPackage := msg.Desc.ParentFile().Package()
	currentPackage := currentFile.Desc.Package()

	// Check if the message's file is different from the current file
	msgFile := msg.Desc.ParentFile()
	currentFileDesc := currentFile.Desc

	return msgPackage == currentPackage && msgFile != currentFileDesc
}

// collectImportedMessagesPython recursively collects all messages that are imported from other packages
func collectImportedMessagesPython(file *protogen.File) []*protogen.Message {
	visited := make(map[string]bool)
	var importedMessages []*protogen.Message

	// Check all messages in the file for imported message references
	for _, message := range file.Messages {
		collectImportedFromMessagePython(message, file, visited, &importedMessages)
	}

	// Check service methods for imported message references
	for _, service := range file.Services {
		for _, method := range service.Methods {
			collectImportedFromMessagePython(method.Input, file, visited, &importedMessages)
			collectImportedFromMessagePython(method.Output, file, visited, &importedMessages)
		}
	}

	return importedMessages
}

// collectImportedFromMessagePython recursively collects imported messages from a message and its fields
func collectImportedFromMessagePython(msg *protogen.Message, currentFile *protogen.File, visited map[string]bool, importedMessages *[]*protogen.Message) {
	if msg == nil {
		return
	}

	messageKey := string(msg.Desc.FullName())

	// If this message is from an imported file and we haven't seen it before
	if !visited[messageKey] && isImportedMessagePython(msg, currentFile) {
		visited[messageKey] = true
		*importedMessages = append(*importedMessages, msg)
	}

	// Recursively check fields for imported message types
	for _, field := range msg.Fields {
		if field.Message != nil {
			collectImportedFromMessagePython(field.Message, currentFile, visited, importedMessages)
		}
	}

	// Check nested messages
	for _, nested := range msg.Messages {
		collectImportedFromMessagePython(nested, currentFile, visited, importedMessages)
	}
}

// isImportedMessagePython checks if a message is imported from another package
func isImportedMessagePython(msg *protogen.Message, currentFile *protogen.File) bool {
	if msg.Desc.ParentFile() == nil {
		return false
	}

	// Check if the message's package is different from the current file's package
	msgPackage := msg.Desc.ParentFile().Package()
	currentPackage := currentFile.Desc.Package()

	return msgPackage != currentPackage
}

// GeneratePythonFile generates Python code for the given protobuf file
func GeneratePythonFile(gen *protogen.Plugin, file *protogen.File, commonNamespace string) {
	if len(file.Messages) == 0 && len(file.Services) == 0 {
		return
	}

	// Generate global transport if namespace is provided and we have services
	if commonNamespace != "" && len(file.Services) > 0 {
		generateGlobalTransport(gen, commonNamespace)
	}

	// Get Python module name and create package structure
	moduleName := getPythonModuleName(file)

	// Create package directories with __init__.py files
	createPythonPackageStructure(gen, moduleName)

	// Create the main module file in the package directory
	baseFilename := filepath.Base(strings.ReplaceAll(*file.Proto.Name, ".proto", ""))
	filename := strings.ReplaceAll(moduleName, ".", "/") + "/" + baseFilename + ".py"
	g := gen.NewGeneratedFile(filename, "")

	// Generate file header
	g.P("# Code generated by protoc-gen-puregen. DO NOT EDIT.")
	g.P()
	g.P("from dataclasses import dataclass, field")
	g.P("from typing import Optional, List, Dict, Any")
	g.P("from abc import ABC, abstractmethod")
	g.P("import json")

	// Import global transport if namespace is provided and we have services
	if commonNamespace != "" && len(file.Services) > 0 {
		g.P("from ", commonNamespace, " import Transport")
	}

	// Collect and generate imports for same package messages
	samePackageMessages := collectSamePackageMessagesPython(file)
	if len(samePackageMessages) > 0 {
		for _, message := range samePackageMessages {
			// Get the module name for the message (based on its file)
			msgFilename := strings.TrimSuffix(filepath.Base(message.Desc.ParentFile().Path()), ".proto")
			// Use full package path with dots for import
			packagePath := getPythonImportModuleName(file)
			g.P("from ", packagePath, ".", msgFilename, " import ", message.GoIdent.GoName)
		}
	}
	g.P()

	// Collect and generate imported messages first
	importedMessages := collectImportedMessagesPython(file)
	if len(importedMessages) > 0 {
		g.P("# Imported Messages (redefined locally)")
		g.P()
		for _, message := range importedMessages {
			generatePythonMessage(g, message)
		}
	}

	// Generate messages
	if len(file.Messages) > 0 {
		g.P("# Messages")
		g.P()
	}
	for _, message := range file.Messages {
		generatePythonMessage(g, message)
	}

	// Generate services
	if len(file.Services) > 0 {
		g.P("# Services")
		g.P()
	}
	for _, service := range file.Services {
		generatePythonService(g, service)
	}

	// Generate method name constants
	if len(file.Services) > 0 {
		g.P("# Method name constants")
		g.P()
		for _, service := range file.Services {
			generatePythonMethodConstants(g, service)
		}
	}

	if len(file.Services) > 0 {
		g.P("# Client")
		g.P()
	}
	for _, service := range file.Services {
		generatePythonClient(g, service, commonNamespace)
	}
}

func generatePythonMethodConstants(g *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName

	g.P("class ", serviceName, "Methods:")
	g.P("    \"\"\"Method name constants for ", serviceName, "\"\"\"")

	for _, method := range service.Methods {
		constName := serviceName + "_" + method.GoName
		g.P("    ", constName, " = \"", constName, "\"")
	}
	g.P()

	// Generate method metadata dictionary
	g.P("    METHOD_METADATA = {")
	for _, method := range service.Methods {
		constName := serviceName + "_" + method.GoName
		metadata := parseMethodMetadata(method.Comments)
		if metadata != nil {
			g.P("        ", constName, ": {")

			// Sort keys for consistent output
			var keys []string
			for key := range metadata {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := metadata[key]
				g.P("            \"", key, "\": \"", value, "\",")
			}
			g.P("        },")
		}
	}
	g.P("    }")
	g.P()
}

func generatePythonMessage(g *protogen.GeneratedFile, msg *protogen.Message) {
	// Generate message comment
	writePythonComment(g, msg.Comments)

	// Generate dataclass
	g.P("@dataclass")
	g.P("class ", msg.GoIdent.GoName, ":")
	g.P("    \"\"\"Generated message class for ", msg.GoIdent.GoName, "\"\"\"")

	// Generate fields with default values
	if len(msg.Fields) == 0 {
		g.P("    pass")
	} else {
		for _, field := range msg.Fields {
			// Generate field comment
			if commentLines := formatPythonComment(field.Comments); len(commentLines) > 0 {
				for _, line := range commentLines {
					g.P("    ", line)
				}
			}

			fieldType := getPythonFieldType(field)
			fieldName := getPythonFieldName(field.GoName)
			defaultValue := getPythonDefaultValue(field)
			g.P("    ", fieldName, ": ", fieldType, " = ", defaultValue)
		}
	}
	g.P()

	// Generate validation method
	g.P("    def validate(self) -> bool:")
	g.P("        \"\"\"Validate the message fields\"\"\"")
	g.P("        # Add custom validation logic here")
	g.P("        return True")
	g.P()

	// Generate JSON serialization methods
	g.P("    def to_json(self) -> str:")
	g.P("        \"\"\"Convert message to JSON string\"\"\"")
	g.P("        return json.dumps(self.to_dict())")
	g.P()

	g.P("    def to_dict(self) -> Dict[str, Any]:")
	g.P("        \"\"\"Convert message to dictionary\"\"\"")
	g.P("        result = {}")
	for _, field := range msg.Fields {
		fieldName := getPythonFieldName(field.GoName)
		jsonName := field.Desc.JSONName()
		g.P("        if self.", fieldName, " is not None:")
		if field.Desc.IsList() {
			if field.Message != nil {
				g.P("            result['", jsonName, "'] = [item.to_dict() if hasattr(item, 'to_dict') else item for item in self.", fieldName, "]")
			} else {
				g.P("            result['", jsonName, "'] = self.", fieldName)
			}
		} else if field.Message != nil {
			g.P("            result['", jsonName, "'] = self.", fieldName, ".to_dict() if hasattr(self.", fieldName, ", 'to_dict') else self.", fieldName)
		} else {
			g.P("            result['", jsonName, "'] = self.", fieldName)
		}
	}
	g.P("        return result")
	g.P()

	g.P("    @classmethod")
	g.P("    def from_json(cls, json_str: str) -> '", msg.GoIdent.GoName, "':")
	g.P("        \"\"\"Create message from JSON string\"\"\"")
	g.P("        data = json.loads(json_str)")
	g.P("        return cls.from_dict(data)")
	g.P()

	g.P("    @classmethod")
	g.P("    def from_dict(cls, data: Dict[str, Any]) -> '", msg.GoIdent.GoName, "':")
	g.P("        \"\"\"Create message from dictionary\"\"\"")
	g.P("        kwargs = {}")
	for _, field := range msg.Fields {
		fieldName := getPythonFieldName(field.GoName)
		jsonName := field.Desc.JSONName()
		if field.Desc.IsList() {
			if field.Message != nil {
				g.P("        if '", jsonName, "' in data:")
				g.P("            kwargs['", fieldName, "'] = [", field.Message.GoIdent.GoName, ".from_dict(item) if isinstance(item, dict) else item for item in data['", jsonName, "']]")
			} else {
				g.P("        if '", jsonName, "' in data:")
				g.P("            kwargs['", fieldName, "'] = data['", jsonName, "']")
			}
		} else if field.Message != nil {
			g.P("        if '", jsonName, "' in data:")
			g.P("            kwargs['", fieldName, "'] = ", field.Message.GoIdent.GoName, ".from_dict(data['", jsonName, "']) if isinstance(data['", jsonName, "'], dict) else data['", jsonName, "']")
		} else {
			g.P("        if '", jsonName, "' in data:")
			g.P("            kwargs['", fieldName, "'] = data['", jsonName, "']")
		}
	}
	g.P("        return cls(**kwargs)")
	g.P()

	// Generate nested messages
	for _, nested := range msg.Messages {
		generatePythonMessage(g, nested)
	}
}

func generatePythonService(g *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName

	// Generate service comment
	writePythonComment(g, service.Comments)

	// Generate abstract service interface
	g.P("class ", serviceName, "Service(ABC):")
	g.P("    \"\"\"Abstract service interface for ", serviceName, "\"\"\"")
	g.P()

	for _, method := range service.Methods {
		// Generate method comment
		if commentLines := formatPythonComment(method.Comments); len(commentLines) > 0 {
			for _, line := range commentLines {
				g.P("    ", line)
			}
		}

		inputType := method.Input.GoIdent.GoName
		outputType := method.Output.GoIdent.GoName
		methodName := getPythonMethodName(method.GoName)

		g.P("    @abstractmethod")
		g.P("    def ", methodName, "(self, ctx: Dict[str, Any], request: ", inputType, ") -> ", outputType, ":")
		g.P("        \"\"\"", method.GoName, " method\"\"\"")
		g.P("        pass")
		g.P()
	}

	// Generate default implementation
	g.P("class Default", serviceName, "Service(", serviceName, "Service):")
	g.P("    \"\"\"Default implementation of ", serviceName, "Service\"\"\"")
	g.P()

	for _, method := range service.Methods {
		// Generate method comment
		if commentLines := formatPythonComment(method.Comments); len(commentLines) > 0 {
			for _, line := range commentLines {
				g.P("    ", line)
			}
		}

		inputType := method.Input.GoIdent.GoName
		outputType := method.Output.GoIdent.GoName
		methodName := getPythonMethodName(method.GoName)

		g.P("    def ", methodName, "(self, ctx: Dict[str, Any], request: ", inputType, ") -> ", outputType, ":")
		g.P("        \"\"\"", method.GoName, " method implementation\"\"\"")
		g.P("        # TODO: Implement ", methodName)
		g.P("        raise NotImplementedError(\"Method ", methodName, " not implemented\")")
		g.P()
	}
}

func generatePythonClient(g *protogen.GeneratedFile, service *protogen.Service, commonNamespace string) {
	serviceName := service.GoName

	// Only generate Transport interface if no global namespace is provided
	if commonNamespace == "" {
		// Generate Transport interface
		g.P("class Transport(ABC):")
		g.P("    \"\"\"Abstract transport interface for client communication\"\"\"")
		g.P()
		g.P("    @abstractmethod")
		g.P("    def send(self, ctx: Dict[str, Any], method_name: str, input_data: Any, output_type: type) -> Any:")
		g.P("        \"\"\"Send request and return response\"\"\"")
		g.P("        pass")
		g.P()
	}

	// Generate client class
	g.P("class ", serviceName, "Client:")
	g.P("    \"\"\"Client for ", serviceName, " service\"\"\"")
	g.P()
	g.P("    def __init__(self, transport: Transport):")
	g.P("        self.transport = transport")
	g.P()

	// Generate client methods
	for _, method := range service.Methods {
		inputType := method.Input.GoIdent.GoName
		outputType := method.Output.GoIdent.GoName
		methodName := getPythonMethodName(method.GoName)
		constName := serviceName + "Methods." + serviceName + "_" + method.GoName

		g.P("    def ", methodName, "(self, ctx: Dict[str, Any], request: ", inputType, ") -> ", outputType, ":")
		g.P("        \"\"\"", method.GoName, " client method\"\"\"")
		g.P("        result = self.transport.send(ctx, ", constName, ", request, ", outputType, ")")
		g.P("        if isinstance(result, ", outputType, "):")
		g.P("            return result")
		g.P("        if isinstance(result, dict):")
		g.P("            return ", outputType, ".from_dict(result)")
		g.P("        raise ValueError(f\"Invalid response type for ", methodName, ": {type(result)}\")")
		g.P()
	}
}

// createPythonPackageStructure creates directories and __init__.py files for the package hierarchy
func createPythonPackageStructure(gen *protogen.Plugin, moduleName string) {
	// For single level package, create __init__.py in the module directory
	if !strings.Contains(moduleName, ".") {
		initFile := moduleName + "/__init__.py"
		if !createdPythonPackages[initFile] {
			initGen := gen.NewGeneratedFile(initFile, "")
			initGen.P("# Package initialization file")
			initGen.P("# Generated by protoc-gen-puregen")
			createdPythonPackages[initFile] = true
		}
		return
	}

	// For multi-level package, only create __init__.py in the final directory
	parts := strings.Split(moduleName, ".")
	finalPath := strings.Join(parts, "/")
	initFile := finalPath + "/__init__.py"

	if !createdPythonPackages[initFile] {
		initGen := gen.NewGeneratedFile(initFile, "")
		initGen.P("# Package initialization file")
		initGen.P("# Generated by protoc-gen-puregen")
		createdPythonPackages[initFile] = true
	}
}

func getPythonModuleName(file *protogen.File) string {
	// Convert proto package to Python module name
	pkg := string(file.Desc.Package())
	if pkg == "" {
		// Use filename without extension
		return strings.TrimSuffix(filepath.Base(file.Desc.Path()), ".proto")
	}

	// Replace dots with slashes for Python package structure
	pkg = strings.ReplaceAll(pkg, ".", "/")
	// Ensure the package name is a valid Python identifier
	pkg = strings.ReplaceAll(pkg, "-", "_")
	// Remove leading slashes if any
	pkg = strings.TrimPrefix(pkg, "/")
	return pkg
}

func getPythonImportModuleName(file *protogen.File) string {
	// Convert proto package to Python import module name (with dots)
	pkg := string(file.Desc.Package())
	if pkg == "" {
		// Use filename without extension
		return strings.TrimSuffix(filepath.Base(file.Desc.Path()), ".proto")
	}

	// Keep dots for Python import statements
	// Ensure the package name is a valid Python identifier
	pkg = strings.ReplaceAll(pkg, "-", "_")
	return pkg
}

func getPythonFieldType(field *protogen.Field) string {
	baseType := ""
	switch field.Desc.Kind().String() {
	case "bool":
		baseType = "bool"
	case "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64",
		"uint32", "fixed32", "uint64", "fixed64":
		baseType = "int"
	case "float", "double":
		baseType = "float"
	case "string":
		baseType = "str"
	case "bytes":
		baseType = "bytes"
	case "enum":
		baseType = "int" // Enum as int for simplicity
	case "message":
		baseType = field.Message.GoIdent.GoName
	default:
		baseType = "Any"
	}

	if field.Desc.IsList() {
		if field.Desc.Kind().String() == "message" {
			return "List['" + baseType + "']"
		} else {
			return "List[" + baseType + "]"
		}
	}

	if field.Desc.Kind().String() == "message" {
		return "Optional['" + baseType + "']"
	}

	return baseType
}

func getPythonFieldName(goName string) string {
	// Convert PascalCase to snake_case
	var result strings.Builder
	for i, r := range goName {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func getPythonMethodName(goName string) string {
	return getPythonFieldName(goName)
}

func getPythonDefaultValue(field *protogen.Field) string {
	if field.Desc.IsList() {
		return "field(default_factory=list)"
	}

	switch field.Desc.Kind().String() {
	case "bool":
		return "False"
	case "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64",
		"uint32", "fixed32", "uint64", "fixed64":
		return "0"
	case "float", "double":
		return "0.0"
	case "string":
		return `""`
	case "bytes":
		return "b''"
	case "enum":
		return "0"
	case "message":
		return "None"
	default:
		return "None"
	}
}

// Track created transport namespaces to avoid duplicates
var createdTransportNamespaces = make(map[string]bool)

// generateGlobalTransport creates a global Transport class in the specified namespace
func generateGlobalTransport(gen *protogen.Plugin, commonNamespace string) {
	// Only create once per namespace
	if createdTransportNamespaces[commonNamespace] {
		return
	}
	createdTransportNamespaces[commonNamespace] = true

	// Create package structure for parent directories
	createTransportPackageStructure(gen, commonNamespace)

	// Create the transport module file
	filename := strings.ReplaceAll(commonNamespace, ".", "/") + "/transport.py"
	g := gen.NewGeneratedFile(filename, "")

	// Generate file header
	g.P("# Code generated by protoc-gen-puregen. DO NOT EDIT.")
	g.P("# Global Transport interface")
	g.P()
	g.P("from abc import ABC, abstractmethod")
	g.P("from typing import Dict, Any")
	g.P()

	// Generate Transport interface
	g.P("class Transport(ABC):")
	g.P("    \"\"\"Abstract transport interface for client communication\"\"\"")
	g.P()
	g.P("    @abstractmethod")
	g.P("    def send(self, ctx: Dict[str, Any], method_name: str, input_data: Any, output_type: type) -> Any:")
	g.P("        \"\"\"Send request and return response\"\"\"")
	g.P("        pass")

	// Create a proper __init__.py file to export Transport
	initFilename := strings.ReplaceAll(commonNamespace, ".", "/") + "/__init__.py"
	initG := gen.NewGeneratedFile(initFilename, "")
	initG.P("# Code generated by protoc-gen-puregen. DO NOT EDIT.")
	initG.P("from .transport import Transport")
	initG.P()
	initG.P("__all__ = ['Transport']")
}

// createTransportPackageStructure creates package directories for transport namespace
func createTransportPackageStructure(gen *protogen.Plugin, commonNamespace string) {
	if !strings.Contains(commonNamespace, ".") {
		// Single level package, nothing more to create
		return
	}

	// For multi-level packages, create parent __init__.py files
	parts := strings.Split(commonNamespace, ".")
	for i := 1; i < len(parts); i++ {
		parentPath := strings.Join(parts[:i], "/")
		initFile := parentPath + "/__init__.py"

		if !createdPythonPackages[initFile] {
			parentG := gen.NewGeneratedFile(initFile, "")
			parentG.P("# Code generated by protoc-gen-puregen. DO NOT EDIT.")
			parentG.P("# Package initialization file")
			createdPythonPackages[initFile] = true
		}
	}
}
