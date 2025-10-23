package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid strong password",
			password:    "MyP@ssw0rd123!",
			expectError: false,
		},
		{
			name:        "valid password with special chars",
			password:    "Secure#Pass1",
			expectError: false,
		},
		{
			name:        "password too short",
			password:    "Pass1!",
			expectError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "missing uppercase",
			password:    "myp@ssw0rd123!",
			expectError: true,
			errorMsg:    "uppercase",
		},
		{
			name:        "missing lowercase",
			password:    "MYP@SSW0RD123!",
			expectError: true,
			errorMsg:    "lowercase",
		},
		{
			name:        "missing number",
			password:    "MyP@ssword!",
			expectError: true,
			errorMsg:    "number",
		},
		{
			name:        "missing special character",
			password:    "MyPassword123",
			expectError: true,
			errorMsg:    "special character",
		},
		{
			name:        "missing multiple requirements",
			password:    "password",
			expectError: true,
			errorMsg:    "uppercase",
		},
		{
			name:        "empty password",
			password:    "",
			expectError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "exactly 8 characters valid",
			password:    "Pass1@rd",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)

			if tt.expectError {
				assert.Error(t, err, "Expected error for password: %s", tt.password)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg,
						"Error message should contain '%s'", tt.errorMsg)
				}
			} else {
				assert.NoError(t, err, "Expected no error for password: %s", tt.password)
			}
		})
	}
}

func TestValidatePasswordWithRequirements(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		requirements PasswordRequirements
		expectError  bool
	}{
		{
			name:     "custom min length 12",
			password: "Short1!Aa",
			requirements: PasswordRequirements{
				MinLength:      12,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: true,
			},
			expectError: true,
		},
		{
			name:     "no special char required",
			password: "MyPassword123",
			requirements: PasswordRequirements{
				MinLength:      8,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: false,
			},
			expectError: false,
		},
		{
			name:     "only length check",
			password: "abcdefgh",
			requirements: PasswordRequirements{
				MinLength:      8,
				RequireUpper:   false,
				RequireLower:   false,
				RequireNumber:  false,
				RequireSpecial: false,
			},
			expectError: false,
		},
		{
			name:     "very strict requirements",
			password: "MyP@ssw0rd123!",
			requirements: PasswordRequirements{
				MinLength:      14,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordWithRequirements(tt.password, tt.requirements)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsCommonPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		isCommon bool
	}{
		{
			name:     "password - common",
			password: "password",
			isCommon: true,
		},
		{
			name:     "123456 - common",
			password: "123456",
			isCommon: true,
		},
		{
			name:     "qwerty - common",
			password: "qwerty",
			isCommon: true,
		},
		{
			name:     "MyP@ssw0rd123! - not common",
			password: "MyP@ssw0rd123!",
			isCommon: false,
		},
		{
			name:     "letmein - common",
			password: "letmein",
			isCommon: true,
		},
		{
			name:     "secure_password_123 - not common",
			password: "secure_password_123",
			isCommon: false,
		},
		{
			name:     "passw0rd - common variant",
			password: "passw0rd",
			isCommon: true,
		},
		{
			name:     "case sensitive - Password not in list",
			password: "Password",
			isCommon: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCommonPassword(tt.password)
			assert.Equal(t, tt.isCommon, result,
				"Password '%s' common check mismatch", tt.password)
		})
	}
}

func TestDefaultPasswordRequirements(t *testing.T) {
	req := DefaultPasswordRequirements()

	assert.Equal(t, 8, req.MinLength, "Default min length should be 8")
	assert.True(t, req.RequireUpper, "Should require uppercase by default")
	assert.True(t, req.RequireLower, "Should require lowercase by default")
	assert.True(t, req.RequireNumber, "Should require number by default")
	assert.True(t, req.RequireSpecial, "Should require special char by default")
}

// TestPasswordValidationErrorMessages ensures error messages are user-friendly
func TestPasswordValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		expectedMsgs []string
	}{
		{
			name:     "missing all requirements",
			password: "abc",
			expectedMsgs: []string{
				"at least 8 characters",
			},
		},
		{
			name:     "missing uppercase and number",
			password: "mypassword!",
			expectedMsgs: []string{
				"uppercase",
				"number",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			assert.Error(t, err)

			errMsg := err.Error()
			for _, expectedMsg := range tt.expectedMsgs {
				assert.Contains(t, errMsg, expectedMsg,
					"Error message should contain '%s'", expectedMsg)
			}
		})
	}
}
