package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App    AppConfig
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	Log    LogConfig
	CORS   CORSConfig
}

type AppConfig struct {
	Env string `mapstructure:"APP_ENV"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"SERVER_HOST"`
	Port            string        `mapstructure:"SERVER_PORT"`
	ReadTimeout     time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
}

type DBConfig struct {
	Host            string        `mapstructure:"DB_HOST"`
	Port            string        `mapstructure:"DB_PORT"`
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	Name            string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

// DSN returns the PostgreSQL connection string.
func (d DBConfig) DSN() string {
	return "host=" + d.Host +
		" port=" + d.Port +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" sslmode=" + d.SSLMode
}

type JWTConfig struct {
	Secret string        `mapstructure:"JWT_SECRET"`
	Expiry time.Duration `mapstructure:"JWT_EXPIRY"`
}

type LogConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"`
	Format string `mapstructure:"LOG_FORMAT"`
}

type CORSConfig struct {
	AllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

// Load reads configuration from environment variables and .env file.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Ignore missing .env file in production (env vars will be set directly).
	_ = viper.ReadInConfig()

	setDefaults()

	cfg := &Config{
		App: AppConfig{
			Env: viper.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Host:            viper.GetString("SERVER_HOST"),
			Port:            viper.GetString("SERVER_PORT"),
			ReadTimeout:     viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout:    viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: viper.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
		},
		DB: DBConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSL_MODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: viper.GetDuration("JWT_EXPIRY"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
		},
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "15s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "30s")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "5m")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
}
