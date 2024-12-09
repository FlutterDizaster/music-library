package config

import (
	"errors"
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

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.Any("error", err))

		slog.Info("Using environment variables and default config")
		err = parseEVNs(&config)
		if err != nil {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Failed to unmarshal config", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Config loaded", slog.String("DSN", config.DatabaseDSN))

	return &config, nil
}

// This is needed for viper to work with environment variables without a config file.
// Can be rewrited with reflect package if needed.
//
//nolint:gocognit // This is for viper
func parseEVNs(config *Config) error {
	httpAddr := viper.Get("HTTP_ADDR")
	databaseDSN := viper.Get("DATABASE_DSN")
	dbRetryCount := viper.Get("DB_RETRY_COUNT")
	dbRetryBackoff := viper.Get("DB_RETRY_BACKOFF")
	migrationsPath := viper.Get("MIGRATIONS_PATH")
	detailsServerAddr := viper.Get("DETAILS_SERVER_ADDR")
	detailsServerRetryCount := viper.Get("DETAILS_SERVER_RETRY_COUNT")
	detailsServerRetryBackoff := viper.Get("DETAILS_SERVER_RETRY_BACKOFF")
	detailsServerMaxRetryBackoff := viper.Get("DETAILS_SERVER_MAX_RETRY_BACKOFF")

	var err error

	if httpAddr != nil {
		config.HTTPAddr, err = parseString(httpAddr)
		if err != nil {
			return err
		}
	}

	if databaseDSN != nil {
		config.DatabaseDSN, err = parseString(databaseDSN)
		if err != nil {
			return err
		}
	}

	if dbRetryCount != nil {
		config.DBRetryCount, err = parseInt(dbRetryCount)
		if err != nil {
			return err
		}
	}

	if dbRetryBackoff != nil {
		config.DBRetryBackoff, err = parseDuration(dbRetryBackoff)
		if err != nil {
			return err
		}
	}

	if migrationsPath != nil {
		config.MigrationsPath, err = parseString(migrationsPath)
		if err != nil {
			return err
		}
	}

	if detailsServerAddr != nil {
		config.DetailsServerAddr, err = parseString(detailsServerAddr)
		if err != nil {
			return err
		}
	}

	if detailsServerRetryCount != nil {
		config.DetailsServerRetryCount, err = parseInt(detailsServerRetryCount)
		if err != nil {
			return err
		}
	}

	if detailsServerRetryBackoff != nil {
		config.DetailsServerRetryBackoff, err = parseDuration(detailsServerRetryBackoff)
		if err != nil {
			return err
		}
	}

	if detailsServerMaxRetryBackoff != nil {
		config.DetailsServerMaxRetryBackoff, err = parseDuration(detailsServerMaxRetryBackoff)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseString(source any) (string, error) {
	str, ok := source.(string)
	if ok {
		return str, nil
	}
	return "", errors.New("type assertion failed: source is not a string")
}

func parseInt(source any) (int, error) {
	i, ok := source.(int)
	if ok {
		return i, nil
	}
	return 0, errors.New("type assertion failed: source is not an int")
}

func parseDuration(source any) (time.Duration, error) {
	durStr, ok := source.(string)
	if !ok {
		return 0, errors.New("type assertion failed: source is not a time.Duration")
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, errors.New("type assertion failed: source is not a time.Duration")
	}

	return dur, nil
}
