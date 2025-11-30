package config

import (
	"go-ai/pkg/logger"
	"strings"

	"github.com/spf13/viper"
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
	MinioEndPoint       string `mapstructure:"MINIO_END_POINT"`
	MinioPort           string `mapstructure:"MINIO_PORT"`
	MinioAccessKey      string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey      string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket              string `mapstructure:"MINIO_BUCKET"`
	MinioUseSSL         bool   `mapstructure:"MINIO_USE_SSL"`
}

func LoadConfig() (*Config, error) {
	logger := logger.NewLogger()
	logger.With().Str("component", "config").Logger()

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
	viper.SetDefault("JWT_SECRET", "your-access-secret-key")
	viper.SetDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key")
	viper.SetDefault("JWT_EXPIRES_IN", 3000)
	viper.SetDefault("JWT_REFRESH_EXPIRES_IN", 6480000)

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_PASSWORD", "")

	// Database defaults
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_DB", "go-ai")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")

	// Environment
	viper.SetDefault("ENVIRONMENT", "development")

	// Minio defaults
	viper.SetDefault("MINIO_END_POINT", "localhost")
	viper.SetDefault("MINIO_PORT", "9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("MINIO_BUCKET", "uploads")
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
