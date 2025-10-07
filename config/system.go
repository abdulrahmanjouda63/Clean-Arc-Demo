package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	} `mapstructure:"db"`
	Redis struct {
		Enabled  bool   `mapstructure:"enabled"`
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret          string `mapstructure:"secret"`
		ExpirationHours int    `mapstructure:"expiration_hours"`
	} `mapstructure:"jwt"`
}

// LoadConfig reads configuration from file
func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

    // Set environment variable prefix and replace dots with underscores
    v.SetEnvPrefix("")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Allow overriding server.port via SERVER_PORT
    // Example: SERVER_PORT=":8082" make run
    _ = v.BindEnv("server.port", "SERVER_PORT")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Set defaults if not provided
	if cfg.Server.Port == "" {
		cfg.Server.Port = ":8080"
	}

	if cfg.DB.Driver == "" {
		cfg.DB.Driver = "mysql"
	}

	return &cfg, nil
}
