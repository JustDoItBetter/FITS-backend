package validation

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeString removes potentially dangerous characters and HTML from user input
// Prevents XSS attacks by escaping HTML entities and removing control characters
func SanitizeString(input string) string {
	// Remove null bytes that can bypass validation
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newline and tab (for multiline text support)
	input = removeControlChars(input)

	// Escape HTML entities to prevent script injection
	input = html.EscapeString(input)

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// SanitizeEmail validates and normalizes email addresses
// Ensures email format and prevents injection attacks
func SanitizeEmail(email string) string {
	// Convert to lowercase for case-insensitive comparison
	email = strings.ToLower(strings.TrimSpace(email))

	// Remove dangerous characters
	email = strings.ReplaceAll(email, "\x00", "")

	// Basic email validation regex - prevents most injection attempts
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "" // Return empty string for invalid emails
	}

	return email
}

// removeControlChars removes ASCII control characters except newline and tab
// Control characters can bypass security checks and cause rendering issues
func removeControlChars(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		// Keep printable characters, newline (10), and tab (9)
		if r >= 32 || r == 10 || r == 9 {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// IsValidUUID checks if a string is a valid UUID format
// Prevents injection attacks through UUID parameters
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuid))
}

// SanitizeName sanitizes person names (first name, last name)
// Allows letters, spaces, hyphens, and apostrophes but prevents injection
func SanitizeName(name string) string {
	name = SanitizeString(name)

	// Remove multiple consecutive spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	name = spaceRegex.ReplaceAllString(name, " ")

	return name
}
