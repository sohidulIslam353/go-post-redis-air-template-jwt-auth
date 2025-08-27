package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ToTitleCase returns a string with first letter uppercase and rest lowercase
func ToTitleCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}

	return string(runes)
}

// MakeSlug converts a string to a URL-friendly slug
func MakeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = re.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	return s
}

func Asset(path string) string {
	return "/static/" + path
}
