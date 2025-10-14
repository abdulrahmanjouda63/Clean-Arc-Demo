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
		// Sentinel configuration
		UseSentinel      bool     `mapstructure:"use_sentinel"`
		SentinelAddrs    []string `mapstructure:"sentinel_addrs"`
		SentinelMaster   string   `mapstructure:"sentinel_master"`
		SentinelPassword string   `mapstructure:"sentinel_password"`
	} `mapstructure:"redis"`

	JWT struct {
		Secret          string `mapstructure:"secret"`
		ExpirationHours int    `mapstructure:"expiration_hours"`
	} `mapstructure:"jwt"`

	Email struct {
		Enabled  bool   `mapstructure:"enabled"`
		SMTPHost string `mapstructure:"smtp_host"`
		SMTPPort int    `mapstructure:"smtp_port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
	} `mapstructure:"email"`

	CORS struct {
		Enabled        bool     `mapstructure:"enabled"`
		AllowedOrigins []string `mapstructure:"allowed_origins"`
		AllowedMethods []string `mapstructure:"allowed_methods"`
		AllowedHeaders []string `mapstructure:"allowed_headers"`
	} `mapstructure:"cors"`

	Logging struct {
		Level      string `mapstructure:"level"`
		OutputPath string `mapstructure:"output_path"`
	} `mapstructure:"logging"`
}

// LoadConfig reads configuration from file
func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	// Set environment variable prefix and replace dots with underscores
	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Allow overriding via environment variables
	_ = v.BindEnv("server.port", "SERVER_PORT")
	_ = v.BindEnv("db.dsn", "DB_DSN")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("email.password", "EMAIL_PASSWORD")

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

	if cfg.JWT.ExpirationHours == 0 {
		cfg.JWT.ExpirationHours = 24
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	// CORS defaults
	if cfg.CORS.Enabled && len(cfg.CORS.AllowedOrigins) == 0 {
		cfg.CORS.AllowedOrigins = []string{"*"}
	}
	if cfg.CORS.Enabled && len(cfg.CORS.AllowedMethods) == 0 {
		cfg.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if cfg.CORS.Enabled && len(cfg.CORS.AllowedHeaders) == 0 {
		cfg.CORS.AllowedHeaders = []string{"Origin", "Content-Type", "Authorization"}
	}

	return &cfg, nil
}
