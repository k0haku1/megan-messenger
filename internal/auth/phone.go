package auth

import (
	"regexp"
	"strings"
	"unicode"
)

var e164Pattern = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func NormalizePhone(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ErrInvalidPhone
	}

	var b strings.Builder
	for i, r := range trimmed {
		if i == 0 && r == '+' {
			b.WriteRune(r)
			continue
		}
		if unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		if unicode.IsSpace(r) || r == '-' || r == '(' || r == ')' {
			continue
		}
		return "", ErrInvalidPhone
	}

	phone := b.String()
	if !strings.HasPrefix(phone, "+") || !e164Pattern.MatchString(phone) {
		return "", ErrInvalidPhone
	}
	return phone, nil
}
