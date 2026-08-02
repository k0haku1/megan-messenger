package project

import (
	"regexp"
	"strings"
)

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func SlugFromTitle(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = slugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "doc"
	}
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimRight(slug, "-")
	}
	return slug
}
