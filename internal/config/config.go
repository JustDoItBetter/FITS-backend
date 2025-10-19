package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Secrets  SecretsConfig  `toml:"secrets"`
	Storage  StorageConfig  `toml:"storage"`
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	JWT      JWTConfig      `toml:"jwt"`
}

type ServerConfig struct {
	Port           int    `toml:"port"`
	Host           string `toml:"host"`
	ReadTimeout    string `toml:"read_timeout"`
	WriteTimeout   string `toml:"write_timeout"`
	AllowedOrigins string `toml:"allowed_origins"` // Comma-separated list of allowed CORS origins
	RateLimit      int    `toml:"rate_limit"`      // Max requests per minute per IP (0 = unlimited)
}

type SecretsConfig struct {
	MetricsSecret      string `toml:"metrics_secret"`
	RegistrationSecret string `toml:"registration_secret"`
	DeletionSecret     string `toml:"deletion_secret"`
	UpdateSecret       string `toml:"update_secret"`
}

type StorageConfig struct {
	UploadDir   string `toml:"upload_dir"`
	MaxFileSize int64  `toml:"max_file_size"`
}

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
func (j *JWTConfig) GetAccessTokenExpiry() time.Duration {
	d, _ := time.ParseDuration(j.AccessTokenExpiry)
	return d
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (j *JWTConfig) GetRefreshTokenExpiry() time.Duration {
	d, _ := time.ParseDuration(j.RefreshTokenExpiry)
	return d
}

// GetInvitationExpiry returns the invitation expiry duration
func (j *JWTConfig) GetInvitationExpiry() time.Duration {
	d, _ := time.ParseDuration(j.InvitationExpiry)
	return d
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Server.Port)
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

func (c *Config) GetReadTimeout() time.Duration {
	d, _ := time.ParseDuration(c.Server.ReadTimeout)
	return d
}

func (c *Config) GetWriteTimeout() time.Duration {
	d, _ := time.ParseDuration(c.Server.WriteTimeout)
	return d
}
