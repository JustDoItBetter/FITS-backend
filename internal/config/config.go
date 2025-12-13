// Package config provides configuration management for the FITS backend application.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config represents the main application configuration.
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Secrets  SecretsConfig  `toml:"secrets"`
	Storage  StorageConfig  `toml:"storage"`
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	JWT      JWTConfig      `toml:"jwt"`
}

// ServerConfig holds the HTTP server configuration.
type ServerConfig struct {
	Port           int    `toml:"port"`
	Host           string `toml:"host"`
	ReadTimeout    string `toml:"read_timeout"`
	WriteTimeout   string `toml:"write_timeout"`
	AllowedOrigins string `toml:"allowed_origins"` // Comma-separated list of allowed CORS origins
	RateLimit      int    `toml:"rate_limit"`      // Max requests per minute per IP (0 = unlimited)

	// TLS/HTTPS Configuration
	TLSEnabled      bool   `toml:"tls_enabled"`       // Enable HTTPS/TLS
	TLSCertFile     string `toml:"tls_cert_file"`     // Path to TLS certificate file
	TLSKeyFile      string `toml:"tls_key_file"`      // Path to TLS private key file
	TLSAutoRedirect bool   `toml:"tls_auto_redirect"` // Auto-redirect HTTP to HTTPS
}

// SecretsConfig holds the API secrets configuration.
type SecretsConfig struct {
	MetricsSecret      string `toml:"metrics_secret"`
	RegistrationSecret string `toml:"registration_secret"`
	DeletionSecret     string `toml:"deletion_secret"`
	UpdateSecret       string `toml:"update_secret"`
}

// StorageConfig holds the file storage configuration.
type StorageConfig struct {
	UploadDir   string `toml:"upload_dir"`
	MaxFileSize int64  `toml:"max_file_size"`
}

// LoggingConfig holds the logging configuration.
type LoggingConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"ssl_mode"`
	MaxConns int    `toml:"max_conns"`
	MinConns int    `toml:"min_conns"`
}

// JWTConfig contains JWT token settings
type JWTConfig struct {
	Secret             string `toml:"secret"`
	AccessTokenExpiry  string `toml:"access_token_expiry"`
	RefreshTokenExpiry string `toml:"refresh_token_expiry"`
	InvitationExpiry   string `toml:"invitation_expiry"`
	AdminKeyPath       string `toml:"admin_key_path"`
	AdminPubKeyPath    string `toml:"admin_pub_key_path"`
}

// GetAccessTokenExpiry returns the access token expiry duration
// Panics if the duration cannot be parsed (should never happen after Validate())
func (j *JWTConfig) GetAccessTokenExpiry() time.Duration {
	d, err := time.ParseDuration(j.AccessTokenExpiry)
	if err != nil {
		// This should never happen if Validate() was called properly
		// Panic is appropriate here because invalid token expiry is a critical security issue
		panic(fmt.Sprintf("invalid access_token_expiry '%s': %v (config should have been validated)", j.AccessTokenExpiry, err))
	}
	return d
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
// Panics if the duration cannot be parsed (should never happen after Validate())
func (j *JWTConfig) GetRefreshTokenExpiry() time.Duration {
	d, err := time.ParseDuration(j.RefreshTokenExpiry)
	if err != nil {
		// This should never happen if Validate() was called properly
		panic(fmt.Sprintf("invalid refresh_token_expiry '%s': %v (config should have been validated)", j.RefreshTokenExpiry, err))
	}
	return d
}

// GetInvitationExpiry returns the invitation expiry duration
// Panics if the duration cannot be parsed (should never happen after Validate())
func (j *JWTConfig) GetInvitationExpiry() time.Duration {
	d, err := time.ParseDuration(j.InvitationExpiry)
	if err != nil {
		// This should never happen if Validate() was called properly
		panic(fmt.Sprintf("invalid invitation_expiry '%s': %v (config should have been validated)", j.InvitationExpiry, err))
	}
	return d
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		if _, err := fmt.Sscanf(port, "%d", &cfg.Server.Port); err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT value: %w", err)
		}
	}

	if metricsSecret := os.Getenv("METRICS_SECRET"); metricsSecret != "" {
		cfg.Secrets.MetricsSecret = metricsSecret
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	// Server validation
	if c.Server.Port == 0 {
		return fmt.Errorf("server.port must be set")
	}

	if _, err := time.ParseDuration(c.Server.ReadTimeout); err != nil {
		return fmt.Errorf("invalid server.read_timeout: %w", err)
	}

	if _, err := time.ParseDuration(c.Server.WriteTimeout); err != nil {
		return fmt.Errorf("invalid server.write_timeout: %w", err)
	}

	// TLS validation (if enabled)
	if c.Server.TLSEnabled {
		if c.Server.TLSCertFile == "" {
			return fmt.Errorf("server.tls_cert_file must be set when TLS is enabled")
		}
		if c.Server.TLSKeyFile == "" {
			return fmt.Errorf("server.tls_key_file must be set when TLS is enabled")
		}
	}

	// Secrets validation
	if c.Secrets.MetricsSecret == "" {
		return fmt.Errorf("secrets.metrics_secret must be set")
	}

	// Database validation
	if c.Database.Host == "" {
		return fmt.Errorf("database.host must be set")
	}

	if c.Database.Port == 0 {
		return fmt.Errorf("database.port must be set")
	}

	if c.Database.Database == "" {
		return fmt.Errorf("database.database must be set")
	}

	// JWT validation
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret must be set")
	}

	// JWT secret must be at least 32 characters (256 bits) for HS256 security
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("jwt.secret must be at least 32 characters for HS256 security (current length: %d)", len(c.JWT.Secret))
	}

	if _, err := time.ParseDuration(c.JWT.AccessTokenExpiry); err != nil {
		return fmt.Errorf("invalid jwt.access_token_expiry: %w", err)
	}

	if _, err := time.ParseDuration(c.JWT.RefreshTokenExpiry); err != nil {
		return fmt.Errorf("invalid jwt.refresh_token_expiry: %w", err)
	}

	if _, err := time.ParseDuration(c.JWT.InvitationExpiry); err != nil {
		return fmt.Errorf("invalid jwt.invitation_expiry: %w", err)
	}

	return nil
}

// GetReadTimeout returns the server read timeout duration
// Panics if the duration cannot be parsed (should never happen after Validate())
func (c *Config) GetReadTimeout() time.Duration {
	d, err := time.ParseDuration(c.Server.ReadTimeout)
	if err != nil {
		// This should never happen if Validate() was called properly
		// Panic is appropriate because invalid timeouts can cause server hangs
		panic(fmt.Sprintf("invalid read_timeout '%s': %v (config should have been validated)", c.Server.ReadTimeout, err))
	}
	return d
}

// GetWriteTimeout returns the server write timeout duration
// Panics if the duration cannot be parsed (should never happen after Validate())
func (c *Config) GetWriteTimeout() time.Duration {
	d, err := time.ParseDuration(c.Server.WriteTimeout)
	if err != nil {
		// This should never happen if Validate() was called properly
		panic(fmt.Sprintf("invalid write_timeout '%s': %v (config should have been validated)", c.Server.WriteTimeout, err))
	}
	return d
}
