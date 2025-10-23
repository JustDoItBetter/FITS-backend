package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTConfig_GetAccessTokenExpiry(t *testing.T) {
	tests := []struct {
		name        string
		expiry      string
		expectPanic bool
		expected    time.Duration
	}{
		{
			name:        "valid duration 1h",
			expiry:      "1h",
			expectPanic: false,
			expected:    time.Hour,
		},
		{
			name:        "valid duration 30m",
			expiry:      "30m",
			expectPanic: false,
			expected:    30 * time.Minute,
		},
		{
			name:        "invalid duration - should panic",
			expiry:      "invalid",
			expectPanic: true,
		},
		{
			name:        "invalid duration with space - should panic",
			expiry:      "1 hour",
			expectPanic: true,
		},
		{
			name:        "empty duration - should panic",
			expiry:      "",
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &JWTConfig{
				AccessTokenExpiry: tt.expiry,
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					cfg.GetAccessTokenExpiry()
				}, "Expected panic for invalid duration: %s", tt.expiry)
			} else {
				result := cfg.GetAccessTokenExpiry()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestJWTConfig_GetRefreshTokenExpiry(t *testing.T) {
	tests := []struct {
		name        string
		expiry      string
		expectPanic bool
		expected    time.Duration
	}{
		{
			name:        "valid duration 7d",
			expiry:      "168h", // 7 days
			expectPanic: false,
			expected:    168 * time.Hour,
		},
		{
			name:        "invalid duration",
			expiry:      "7 days",
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &JWTConfig{
				RefreshTokenExpiry: tt.expiry,
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					cfg.GetRefreshTokenExpiry()
				})
			} else {
				result := cfg.GetRefreshTokenExpiry()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestConfig_GetReadTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     string
		expectPanic bool
		expected    time.Duration
	}{
		{
			name:        "valid timeout 30s",
			timeout:     "30s",
			expectPanic: false,
			expected:    30 * time.Second,
		},
		{
			name:        "valid timeout 1m",
			timeout:     "1m",
			expectPanic: false,
			expected:    time.Minute,
		},
		{
			name:        "invalid timeout",
			timeout:     "invalid",
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					ReadTimeout: tt.timeout,
				},
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					cfg.GetReadTimeout()
				})
			} else {
				result := cfg.GetReadTimeout()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestConfig_GetWriteTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     string
		expectPanic bool
		expected    time.Duration
	}{
		{
			name:        "valid timeout 30s",
			timeout:     "30s",
			expectPanic: false,
			expected:    30 * time.Second,
		},
		{
			name:        "invalid timeout",
			timeout:     "30 seconds",
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					WriteTimeout: tt.timeout,
				},
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					cfg.GetWriteTimeout()
				})
			} else {
				result := cfg.GetWriteTimeout()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "this-is-a-very-secure-secret-key-with-32-chars",
					AccessTokenExpiry:  "1h",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: false,
		},
		{
			name: "missing server port",
			config: &Config{
				Server: ServerConfig{
					Port:         0,
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "this-is-a-very-secure-secret-key-with-32-chars",
					AccessTokenExpiry:  "1h",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: true,
			errorMsg:    "server.port must be set",
		},
		{
			name: "JWT secret too short",
			config: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "too-short",
					AccessTokenExpiry:  "1h",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: true,
			errorMsg:    "jwt.secret must be at least 32 characters",
		},
		{
			name: "invalid read timeout",
			config: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  "invalid",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "this-is-a-very-secure-secret-key-with-32-chars",
					AccessTokenExpiry:  "1h",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: true,
			errorMsg:    "invalid server.read_timeout",
		},
		{
			name: "invalid JWT access token expiry",
			config: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "this-is-a-very-secure-secret-key-with-32-chars",
					AccessTokenExpiry:  "1 hour",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: true,
			errorMsg:    "invalid jwt.access_token_expiry",
		},
		{
			name: "empty JWT secret",
			config: &Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Secrets: SecretsConfig{
					MetricsSecret: "test-secret",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "test_db",
				},
				JWT: JWTConfig{
					Secret:             "",
					AccessTokenExpiry:  "1h",
					RefreshTokenExpiry: "168h",
					InvitationExpiry:   "168h",
				},
			},
			expectError: true,
			errorMsg:    "jwt.secret must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Load_EnvironmentOverride(t *testing.T) {
	// Create a temporary config file for testing
	tmpFile, err := os.CreateTemp("", "config-test-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
[server]
port = 8080
host = "0.0.0.0"
read_timeout = "30s"
write_timeout = "30s"
allowed_origins = "*"
rate_limit = 100

[database]
host = "localhost"
port = 5432
user = "test"
password = "test"
database = "test_db"
ssl_mode = "disable"
max_conns = 10
min_conns = 2

[jwt]
secret = "this-is-a-very-secure-secret-key-with-32-characters-long"
access_token_expiry = "1h"
refresh_token_expiry = "168h"
invitation_expiry = "168h"
admin_key_path = "./keys/admin.key"
admin_pub_key_path = "./keys/admin.pub"

[secrets]
metrics_secret = "test-metrics-secret"
registration_secret = "test-registration-secret"
deletion_secret = "test-deletion-secret"
update_secret = "test-update-secret"

[storage]
upload_dir = "./uploads"
max_file_size = 10485760

[logging]
level = "info"
format = "json"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Test environment variable override for SERVER_PORT
	t.Run("SERVER_PORT override", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9090")
		defer os.Unsetenv("SERVER_PORT")

		cfg, err := Load(tmpFile.Name())
		require.NoError(t, err)
		assert.Equal(t, 9090, cfg.Server.Port)
	})

	// Test environment variable override for METRICS_SECRET
	t.Run("METRICS_SECRET override", func(t *testing.T) {
		os.Setenv("METRICS_SECRET", "env-metrics-secret")
		defer os.Unsetenv("METRICS_SECRET")

		cfg, err := Load(tmpFile.Name())
		require.NoError(t, err)
		assert.Equal(t, "env-metrics-secret", cfg.Secrets.MetricsSecret)
	})

	// Test loading without environment overrides
	t.Run("no environment overrides", func(t *testing.T) {
		cfg, err := Load(tmpFile.Name())
		require.NoError(t, err)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "test-metrics-secret", cfg.Secrets.MetricsSecret)
	})
}
