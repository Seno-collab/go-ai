package config

import (
	"strings"

	"github.com/spf13/viper"

	"go-ai/pkg/common"
)

type Config struct {
	JwtAccessSecret     string `mapstructure:"JWT_SECRET"`
	JwtRefreshSecret    string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtExpiresIn        int    `mapstructure:"JWT_EXPIRES_IN"`
	JwtRefreshExpiresIn int    `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
	RedisHost           string `mapstructure:"REDIS_HOST"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	RedisPort           int    `mapstructure:"REDIS_PORT"`
	RedisDB             int    `mapstructure:"REDIS_DB"`
	DBName              string `mapstructure:"POSTGRES_DB"`
	DBHost              string `mapstructure:"POSTGRES_HOST"`
	DBPort              string `mapstructure:"POSTGRES_PORT"`
	DBUser              string `mapstructure:"POSTGRES_USER"`
	DBPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	DBSSLMode           string `mapstructure:"db_sslmode"`
	ServerPort          string `mapstructure:"PORT"`
	ServerHost          string `mapstructure:"server_host"`
	Environment         string `mapstructure:"ENVIRONMENT"`
}

func LoadConfig() (*Config, error) {
	logger := common.Logger.With().Str("component", "config").Logger()

	// Set default values
	setDefaults()

	// Set environment variable key replacer to convert dots to underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable automatic environment variable reading with highest priority
	viper.AutomaticEnv()

	// Ensure environment variables take precedence over config files
	viper.SetEnvPrefix("")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warn().Err(err).Msg("error reading .env file")
		}
	} else {
		logger.Info().Msg(".env file loaded for backward compatibility")
	}

	// Create config struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		logger.Error().Err(err).Msg("unable to decode config")
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	// JWT defaults
	viper.SetDefault("jwt_access_secret", "your-access-secret-key")
	viper.SetDefault("jwt_refresh_secret", "your-refresh-secret-key")
	viper.SetDefault("jwt_expires_in", 3000)            // in second
	viper.SetDefault("jwt_refresh_expires_in", 6480000) // in second (30 days)
	// Redis defaults
	viper.SetDefault("redis_host", "localhost")
	viper.SetDefault("redis_port", 6379)
	viper.SetDefault("redis_db", 0)
	viper.SetDefault("redis_password", "")

	// Database defaults
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", "5432")
	viper.SetDefault("db_name", "go-ai")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "")
	viper.SetDefault("db_sslmode", "disable")

	// Server defaults
	viper.SetDefault("server_host", "0.0.0.0")
	viper.SetDefault("server_port", "8080")

	// Environment
	viper.SetDefault("environment", "development")
}

// GetString returns a string value from config
func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an int value from config
func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a bool value from config
func (c *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
