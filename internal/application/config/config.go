package config

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPAddr = ":8080"

	defaultDBRetryCount   = 3
	defaultDBRetryBackoff = "1s"

	defaultMigrationsPath = "./migrations"

	defaultDetailsServerAddr            = "http://localhost:8081"
	defaultDdetailsServerRetryCount     = 3
	defaultDetailsServerRetryBackoff    = "1s"
	defaultDetailsServerMaxRetryBackoff = "10s"
)

// Config represents the application configuration.
type Config struct {
	HTTPAddr string `mapstructure:"HTTP_ADDR"`

	DatabaseDSN    string        `mapstructure:"DATABASE_DSN"`
	DBRetryCount   int           `mapstructure:"DB_RETRY_COUNT"`
	DBRetryBackoff time.Duration `mapstructure:"DB_RETRY_BACKOFF"`

	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`

	DetailsServerAddr            string        `mapstructure:"DETAILS_SERVER_ADDR"`
	DetailsServerRetryCount      int           `mapstructure:"DETAILS_SERVER_RETRY_COUNT"`
	DetailsServerRetryBackoff    time.Duration `mapstructure:"DETAILS_SERVER_RETRY_BACKOFF"`
	DetailsServerMaxRetryBackoff time.Duration `mapstructure:"DETAILS_SERVER_MAX_RETRY_BACKOFF"`
}

// LoadConfig loads the application configuration from a .env file.
func LoadConfig() (*Config, error) {
	slog.Debug("Loading config")

	viper.SetDefault("HTTP_ADDR", defaultHTTPAddr)

	viper.SetDefault("DB_RETRY_COUNT", defaultDBRetryCount)
	viper.SetDefault("DB_RETRY_BACKOFF", defaultDBRetryBackoff)

	viper.SetDefault("MIGRATIONS_PATH", defaultMigrationsPath)

	viper.SetDefault("DETAILS_SERVER_ADDR", defaultDetailsServerAddr)
	viper.SetDefault("DETAILS_SERVER_RETRY_COUNT", defaultDdetailsServerRetryCount)
	viper.SetDefault("DETAILS_SERVER_RETRY_BACKOFF", defaultDetailsServerRetryBackoff)
	viper.SetDefault("DETAILS_SERVER_MAX_RETRY_BACKOFF", defaultDetailsServerMaxRetryBackoff)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.Any("error", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Failed to unmarshal config", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Config loaded", slog.String("DSN", config.DatabaseDSN))

	return &config, nil
}
