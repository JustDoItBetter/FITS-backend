package validation

import (
	"fmt"
	"unicode"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
)

// PasswordRequirements defines the minimum password strength requirements
type PasswordRequirements struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// DefaultPasswordRequirements returns the default password strength requirements
// Follows OWASP recommendations for password complexity
func DefaultPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:      8,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	}
}

// ValidatePasswordStrength validates that a password meets minimum security requirements
// Returns a validation error with details if password is too weak
func ValidatePasswordStrength(password string) error {
	return ValidatePasswordWithRequirements(password, DefaultPasswordRequirements())
}

// ValidatePasswordWithRequirements validates password against custom requirements
func ValidatePasswordWithRequirements(password string, req PasswordRequirements) error {
	if len(password) < req.MinLength {
		return errors.ValidationError(
			fmt.Sprintf("password must be at least %d characters long", req.MinLength))
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var missing []string

	if req.RequireUpper && !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if req.RequireLower && !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if req.RequireNumber && !hasNumber {
		missing = append(missing, "number")
	}
	if req.RequireSpecial && !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		msg := "password must contain at least one "
		for i, item := range missing {
			if i > 0 {
				if i == len(missing)-1 {
					msg += " and "
				} else {
					msg += ", "
				}
			}
			msg += item
		}
		return errors.ValidationError(msg)
	}

	return nil
}

// IsCommonPassword checks if password is in a list of commonly used passwords
// This is a basic implementation - production should use a comprehensive breach database
func IsCommonPassword(password string) bool {
	// Top 20 most common passwords according to security research
	// In production, use a proper breach password database like HaveIBeenPwned API
	commonPasswords := map[string]bool{
		"password": true,
		"123456":   true,
		"12345678": true,
		"qwerty":   true,
		"abc123":   true,
		"monkey":   true,
		"1234567":  true,
		"letmein":  true,
		"trustno1": true,
		"dragon":   true,
		"baseball": true,
		"111111":   true,
		"iloveyou": true,
		"master":   true,
		"sunshine": true,
		"ashley":   true,
		"bailey":   true,
		"passw0rd": true,
		"shadow":   true,
		"123123":   true,
		"654321":   true,
		"superman": true,
		"qazwsx":   true,
		"michael":  true,
		"Football": true,
	}

	return commonPasswords[password]
}
