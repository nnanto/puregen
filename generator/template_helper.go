package generator

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TemplateFuncs returns a map of template helper functions
func TemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"title":      Title,
		"capitalize": Capitalize,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"snake":      SnakeCase,
		"camel":      CamelCase,
		"pascal":     PascalCase,
		"kebab":      KebabCase,
		"split":      strings.Split,
		"trim":       strings.TrimSpace,
		"join":       strings.Join,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"replace":    strings.ReplaceAll,
		"eq":         Equal,
		"ne":         NotEqual,
		"index":      Index,
	}
}

// Title converts a string to title case (first letter of each word capitalized)
func Title(s string) string {
	return cases.Title(language.English).String(strings.ToLower(s))
}

// Capitalize converts the first character to uppercase
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// SnakeCase converts a string to snake_case
func SnakeCase(s string) string {
	// Insert underscore before uppercase letters (except the first one)
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = reg.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}

// CamelCase converts a string to camelCase
func CamelCase(s string) string {
	if s == "" {
		return s
	}

	// Split by common delimiters
	words := regexp.MustCompile(`[_\s-]+`).Split(s, -1)
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if words[i] != "" {
			result += Capitalize(strings.ToLower(words[i]))
		}
	}
	return result
}

// PascalCase converts a string to PascalCase
func PascalCase(s string) string {
	if s == "" {
		return s
	}

	// Split by common delimiters
	words := regexp.MustCompile(`[_\s-]+`).Split(s, -1)
	var result strings.Builder

	for _, word := range words {
		if word != "" {
			result.WriteString(Capitalize(strings.ToLower(word)))
		}
	}
	return result.String()
}

// KebabCase converts a string to kebab-case
func KebabCase(s string) string {
	// Insert hyphen before uppercase letters (except the first one)
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = reg.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(s)
}

// Equal checks if two values are equal
func Equal(a, b interface{}) bool {
	return a == b
}

// NotEqual checks if two values are not equal
func NotEqual(a, b interface{}) bool {
	return a != b
}

// Index safely gets a value from a map with a string key
func Index(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}
