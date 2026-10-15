package engine

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

// defaultFuncMap retorna el conjunto de funciones utilitarias inyectadas en las plantillas.
func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"toLower":    strings.ToLower,
		"toUpper":    strings.ToUpper,
		"camelCase":  camelCase,
		"pascalCase": pascalCase,
		"snakeCase":  snakeCase,
		"kebabCase":  kebabCase,
		"contains":   contains,
		"join":       join,
	}
}

func splitWords(s string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			if unicode.IsLower(prev) || unicode.IsDigit(prev) {
				if current.Len() > 0 {
					words = append(words, current.String())
					current.Reset()
				}
			} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) && unicode.IsUpper(prev) {
				if current.Len() > 0 {
					words = append(words, current.String())
					current.Reset()
				}
			}
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

func camelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	var res strings.Builder
	for i, w := range words {
		if i == 0 {
			res.WriteString(strings.ToLower(w))
		} else {
			res.WriteString(titleWord(w))
		}
	}
	return res.String()
}

func pascalCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	var res strings.Builder
	for _, w := range words {
		res.WriteString(titleWord(w))
	}
	return res.String()
}

func snakeCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	var lower []string
	for _, w := range words {
		lower = append(lower, strings.ToLower(w))
	}
	return strings.Join(lower, "_")
}

func kebabCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	var lower []string
	for _, w := range words {
		lower = append(lower, strings.ToLower(w))
	}
	return strings.Join(lower, "-")
}

func titleWord(w string) string {
	if len(w) == 0 {
		return ""
	}
	runes := []rune(w)
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}

func contains(collection any, target any) bool {
	targetStr := fmt.Sprintf("%v", target)
	switch c := collection.(type) {
	case []string:
		for _, s := range c {
			if s == targetStr {
				return true
			}
		}
	case string:
		return strings.Contains(c, targetStr)
	case []any:
		for _, v := range c {
			if fmt.Sprintf("%v", v) == targetStr {
				return true
			}
		}
	}
	return false
}

func join(list any, sep string) string {
	switch c := list.(type) {
	case []string:
		return strings.Join(c, sep)
	case []any:
		var strs []string
		for _, v := range c {
			strs = append(strs, fmt.Sprintf("%v", v))
		}
		return strings.Join(strs, sep)
	default:
		return fmt.Sprintf("%v", list)
	}
}
