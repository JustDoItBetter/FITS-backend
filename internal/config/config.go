package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	Secrets SecretsConfig `toml:"secrets"`
	Storage StorageConfig `toml:"storage"`
	Logging LoggingConfig `toml:"logging"`
}

type ServerConfig struct {
	Port         int    `toml:"port"`
	Host         string `toml:"host"`
	ReadTimeout  string `toml:"read_timeout"`
	WriteTimeout string `toml:"write_timeout"`
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
	if c.Server.Port == 0 {
		return fmt.Errorf("server.port must be set")
	}

	if c.Secrets.MetricsSecret == "" {
		return fmt.Errorf("secrets.metrics_secret must be set")
	}

	if c.Secrets.RegistrationSecret == "" {
		return fmt.Errorf("secrets.registration_secret must be set")
	}

	if c.Secrets.DeletionSecret == "" {
		return fmt.Errorf("secrets.deletion_secret must be set")
	}

	if c.Secrets.UpdateSecret == "" {
		return fmt.Errorf("secrets.update_secret must be set")
	}

	if _, err := time.ParseDuration(c.Server.ReadTimeout); err != nil {
		return fmt.Errorf("invalid server.read_timeout: %w", err)
	}

	if _, err := time.ParseDuration(c.Server.WriteTimeout); err != nil {
		return fmt.Errorf("invalid server.write_timeout: %w", err)
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
