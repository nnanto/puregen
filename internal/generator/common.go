package generator

import (
	"encoding/json"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// parseMethodMetadata extracts metadata from method comments using puregen:metadata: directive
func parseMethodMetadata(comments protogen.CommentSet) map[string]string {
	// Use leading comments if available, otherwise trailing
	comment := comments.Leading
	if comment == "" && comments.Trailing != "" {
		comment = comments.Trailing
	}

	if comment == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(comment)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "puregen:metadata:") {
			jsonStr := strings.TrimPrefix(line, "puregen:metadata:")
			jsonStr = strings.TrimSpace(jsonStr)

			var metadata map[string]string
			if err := json.Unmarshal([]byte(jsonStr), &metadata); err == nil {
				return metadata
			}
		}
	}
	return nil
}

// PuregenDirective represents a parsed puregen directive from comments
type PuregenDirective struct {
	EnumType string `json:"enumType,omitempty"`
	// Add other directive fields as needed
}

// parsePuregenDirective extracts puregen directives from comments
func parsePuregenDirective(comments protogen.CommentSet) *PuregenDirective {
	// Use leading comments if available, otherwise trailing
	comment := comments.Leading
	if comment == "" && comments.Trailing != "" {
		comment = comments.Trailing
	}

	if comment == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(comment)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "puregen:generate:") {
			jsonStr := strings.TrimPrefix(line, "puregen:generate:")
			jsonStr = strings.TrimSpace(jsonStr)

			var directive PuregenDirective
			if err := json.Unmarshal([]byte(jsonStr), &directive); err == nil {
				return &directive
			}
		}
	}
	return nil
}

// fileExists checks if a file already exists in the plugin's file list
func fileExists(gen *protogen.Plugin, filename string) bool {
	// Check if the file is already being generated in this run
	for _, file := range gen.Files {
		if file.GeneratedFilenamePrefix+".py" == filename ||
			file.GeneratedFilenamePrefix+".go" == filename ||
			file.GeneratedFilenamePrefix+".java" == filename {
			return true
		}
	}
	// Note: We can't easily check the actual filesystem from within protoc-gen,
	// but we can track what we're generating in this session
	return false
}

// collectImportedMessages recursively collects all messages that are imported from other packages
func collectImportedMessages(file *protogen.File) []*protogen.Message {
	visited := make(map[string]bool)
	var importedMessages []*protogen.Message

	// Check all messages in the file for imported message references
	for _, message := range file.Messages {
		collectImportedFromMessage(message, file, visited, &importedMessages)
	}

	// Check service methods for imported message references
	for _, service := range file.Services {
		for _, method := range service.Methods {
			collectImportedFromMessage(method.Input, file, visited, &importedMessages)
			collectImportedFromMessage(method.Output, file, visited, &importedMessages)
		}
	}

	return importedMessages
}

// collectImportedFromMessage recursively collects imported messages from a message and its fields
func collectImportedFromMessage(msg *protogen.Message, currentFile *protogen.File, visited map[string]bool, importedMessages *[]*protogen.Message) {
	if msg == nil {
		return
	}

	messageKey := string(msg.Desc.FullName())

	// If this message is from an imported file and we haven't seen it before
	if !visited[messageKey] && isImportedMessage(msg, currentFile) {
		visited[messageKey] = true
		*importedMessages = append(*importedMessages, msg)
	}

	// Recursively check fields for imported message types
	for _, field := range msg.Fields {
		if field.Message != nil {
			collectImportedFromMessage(field.Message, currentFile, visited, importedMessages)
		}
	}

	// Check nested messages
	for _, nested := range msg.Messages {
		collectImportedFromMessage(nested, currentFile, visited, importedMessages)
	}
}

// isImportedMessage checks if a message is imported from another package
func isImportedMessage(msg *protogen.Message, currentFile *protogen.File) bool {
	if msg.Desc.ParentFile() == nil {
		return false
	}

	// Check if the message's package is different from the current file's package
	msgPackage := msg.Desc.ParentFile().Package()
	currentPackage := currentFile.Desc.Package()

	return msgPackage != currentPackage
}

// collectSamePackageMessages collects messages from the same package that need to be imported
func collectSamePackageMessages(file *protogen.File) []*protogen.Message {
	visited := make(map[string]bool)
	var samePackageMessages []*protogen.Message

	// Check all messages in the file for same package message references
	for _, message := range file.Messages {
		collectSamePackageFromMessage(message, file, visited, &samePackageMessages)
	}

	// Check service methods for same package message references
	for _, service := range file.Services {
		for _, method := range service.Methods {
			collectSamePackageFromMessage(method.Input, file, visited, &samePackageMessages)
			collectSamePackageFromMessage(method.Output, file, visited, &samePackageMessages)
		}
	}

	return samePackageMessages
}

// collectSamePackageFromMessage recursively collects same package messages from a message and its fields
func collectSamePackageFromMessage(msg *protogen.Message, currentFile *protogen.File, visited map[string]bool, samePackageMessages *[]*protogen.Message) {
	if msg == nil {
		return
	}

	messageKey := string(msg.Desc.FullName())

	// If this message is from the same package but different file and we haven't seen it before
	if !visited[messageKey] && isSamePackageMessage(msg, currentFile) {
		visited[messageKey] = true
		*samePackageMessages = append(*samePackageMessages, msg)
	}

	// Recursively check fields for same package message types
	for _, field := range msg.Fields {
		if field.Message != nil {
			collectSamePackageFromMessage(field.Message, currentFile, visited, samePackageMessages)
		}
	}

	// Check nested messages
	for _, nested := range msg.Messages {
		collectSamePackageFromMessage(nested, currentFile, visited, samePackageMessages)
	}
}

// isSamePackageMessage checks if a message is from the same package but different file
func isSamePackageMessage(msg *protogen.Message, currentFile *protogen.File) bool {
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

// collectAllEnums collects all enums from the file including nested enums
func collectAllEnums(file *protogen.File) []*protogen.Enum {
	var allEnums []*protogen.Enum

	// Add file-level enums
	allEnums = append(allEnums, file.Enums...)

	// Add enums from messages (including nested)
	for _, message := range file.Messages {
		allEnums = append(allEnums, collectEnumsFromMessage(message)...)
	}

	return allEnums
}

// collectEnumsFromMessage recursively collects all enums from a message and its nested messages
func collectEnumsFromMessage(msg *protogen.Message) []*protogen.Enum {
	var enums []*protogen.Enum

	// Add enums from this message
	enums = append(enums, msg.Enums...)

	// Recursively collect from nested messages
	for _, nested := range msg.Messages {
		enums = append(enums, collectEnumsFromMessage(nested)...)
	}

	return enums
}
