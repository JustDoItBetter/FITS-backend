package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("hashes password successfully", func(t *testing.T) {
		password := "MySecurePassword123!"

		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash, "hash should not equal plain password")
	})

	t.Run("generates different hashes for same password", func(t *testing.T) {
		password := "SamePassword123"

		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "bcrypt should use different salts")
	})

	t.Run("rejects empty password", func(t *testing.T) {
		hash, err := HashPassword("")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
		assert.Empty(t, hash)
	})

	t.Run("rejects password exceeding 72 bytes", func(t *testing.T) {
		// Bcrypt has a hard limit of 72 bytes
		password := strings.Repeat("a", 100)

		hash, err := HashPassword(password)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
		assert.Empty(t, hash)
	})

	t.Run("hash starts with bcrypt prefix", func(t *testing.T) {
		password := "TestPassword123"

		hash, err := HashPassword(password)

		require.NoError(t, err)
		// Bcrypt hashes start with $2a$, $2b$, or $2y$
		assert.True(t, strings.HasPrefix(hash, "$2"), "bcrypt hash should start with $2")
	})

	t.Run("hash contains cost factor", func(t *testing.T) {
		password := "TestPassword123"

		hash, err := HashPassword(password)

		require.NoError(t, err)
		// Extract cost from hash and verify it matches BcryptCost (12)
		cost, err := bcrypt.Cost([]byte(hash))
		require.NoError(t, err)
		assert.Equal(t, BcryptCost, cost)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("verifies correct password", func(t *testing.T) {
		password := "CorrectPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(password, hash)

		assert.NoError(t, err)
	})

	t.Run("rejects incorrect password", func(t *testing.T) {
		correctPassword := "CorrectPassword123"
		wrongPassword := "WrongPassword123"
		hash, err := HashPassword(correctPassword)
		require.NoError(t, err)

		err = VerifyPassword(wrongPassword, hash)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password verification failed")
		assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
	})

	t.Run("rejects empty password in verification", func(t *testing.T) {
		password := "ValidPassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Try to verify with empty password
		err = VerifyPassword("", hash)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
	})

	t.Run("rejects empty password against non-empty hash", func(t *testing.T) {
		password := "NonEmptyPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword("", hash)

		assert.Error(t, err)
	})

	t.Run("rejects password with invalid hash", func(t *testing.T) {
		password := "TestPassword"
		invalidHash := "not-a-valid-bcrypt-hash"

		err := VerifyPassword(password, invalidHash)

		assert.Error(t, err)
	})

	t.Run("rejects password with empty hash", func(t *testing.T) {
		err := VerifyPassword("password", "")

		assert.Error(t, err)
	})

	t.Run("is case sensitive", func(t *testing.T) {
		password := "CaseSensitive123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Correct case should work
		err = VerifyPassword(password, hash)
		assert.NoError(t, err)

		// Wrong case should fail
		err = VerifyPassword("casesensitive123", hash)
		assert.Error(t, err)

		err = VerifyPassword("CASESENSITIVE123", hash)
		assert.Error(t, err)
	})

	t.Run("is whitespace sensitive", func(t *testing.T) {
		password := "password123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword("password123 ", hash)
		assert.Error(t, err)

		err = VerifyPassword(" password123", hash)
		assert.Error(t, err)
	})

	t.Run("handles special characters", func(t *testing.T) {
		password := "P@ssw0rd!#$%^&*()_+-=[]{}|;:',.<>?/~`"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(password, hash)
		assert.NoError(t, err)
	})

	t.Run("handles unicode characters", func(t *testing.T) {
		password := "パスワード123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		err = VerifyPassword(password, hash)
		assert.NoError(t, err)
	})
}

func TestHashAndVerifyIntegration(t *testing.T) {
	testPasswords := []string{
		"SimplePassword",
		"Complex!P@ssw0rd#2024",
		"12345678",
		"a",
		strings.Repeat("long", 10), // 40 chars, well under 72 byte limit
		"Ümlaut Pässwörd",
		"unicode_password",
		"   spaces   ",
		"\ttabs\t\t",
		"\nnewlines\n",
	}

	for _, password := range testPasswords {
		t.Run("hash_and_verify_"+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			require.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Verify correct password works
			err = VerifyPassword(password, hash)
			assert.NoError(t, err)

			// Verify wrong password fails
			wrongPassword := password + "wrong"
			err = VerifyPassword(wrongPassword, hash)
			assert.Error(t, err)
		})
	}
}

func TestPasswordHashSecurity(t *testing.T) {
	t.Run("same password produces different hashes", func(t *testing.T) {
		password := "SecurityTest123"

		// Generate multiple hashes
		hashes := make([]string, 10)
		for i := 0; i < 10; i++ {
			hash, err := HashPassword(password)
			require.NoError(t, err)
			hashes[i] = hash
		}

		// All hashes should be unique (due to random salt)
		uniqueHashes := make(map[string]bool)
		for _, hash := range hashes {
			uniqueHashes[hash] = true
		}

		assert.Equal(t, 10, len(uniqueHashes), "all hashes should be unique")
	})

	t.Run("all hashes are valid for the same password", func(t *testing.T) {
		password := "ValidateAll123"

		// Generate multiple hashes
		hashes := make([]string, 5)
		for i := 0; i < 5; i++ {
			hash, err := HashPassword(password)
			require.NoError(t, err)
			hashes[i] = hash
		}

		// All hashes should verify the same password
		for i, hash := range hashes {
			err := VerifyPassword(password, hash)
			assert.NoError(t, err, "hash %d should verify password", i)
		}
	})

	t.Run("hash length is consistent", func(t *testing.T) {
		passwords := []string{"short", "medium-length-password", strings.Repeat("long-", 10)} // 50 chars, under 72 limit

		var hashLength int
		for i, password := range passwords {
			hash, err := HashPassword(password)
			require.NoError(t, err)

			if i == 0 {
				hashLength = len(hash)
			} else {
				assert.Equal(t, hashLength, len(hash), "all bcrypt hashes should have same length")
			}
		}

		// Bcrypt hashes are always 60 characters
		assert.Equal(t, 60, hashLength)
	})
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hash)
	}
}

func BenchmarkHashPasswordParallel(b *testing.B) {
	password := "ParallelPassword123!"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = HashPassword(password)
		}
	})
}

func BenchmarkVerifyPasswordParallel(b *testing.B) {
	password := "ParallelPassword123!"
	hash, _ := HashPassword(password)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = VerifyPassword(password, hash)
		}
	})
}

// Test password strength (documentation purposes - not actual validation)
func TestPasswordExamples(t *testing.T) {
	validExamples := []string{
		"weak",                  // Weak but valid
		"StrongerPass123",       // Better
		"Very!Strong@Pass#2024", // Strong
		"12345678",              // Numeric
		strings.Repeat("a", 60), // Long but under 72 byte limit
	}

	invalidExamples := map[string]string{
		"":                       "password cannot be empty",         // Empty password
		strings.Repeat("a", 100): "password length exceeds 72 bytes", // Too long
	}

	// Test valid passwords
	for _, password := range validExamples {
		t.Run("valid_"+password, func(t *testing.T) {
			hash, err := HashPassword(password)
			require.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Verify it works
			err = VerifyPassword(password, hash)
			assert.NoError(t, err)
		})
	}

	// Test invalid passwords
	for password, expectedError := range invalidExamples {
		t.Run("invalid_"+password, func(t *testing.T) {
			hash, err := HashPassword(password)
			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedError)
			assert.Empty(t, hash)
		})
	}
}
