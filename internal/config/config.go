package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Default values
const (
	defaultHTTPHost               = "localhost"
	defaultHTTPPort               = "8080"
	defaultHTTPWriteTimeout       = 10 * time.Second
	defaultHTTPReadTimeout        = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	defaultLimiterRPS   = 10
	defaultLimiterBurst = 2
	defaultLimiterTTL   = 10 * time.Minute

	defaultSSLMode         = "disable"
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 5 * time.Minute
)

type (
	Config struct {
		App      AppConfig
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Limiter  LimiterConfig
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}

	AppConfig struct {
		AppSecretKey string
		Environment  string
	}

	PostgresConfig struct {
		Host            string
		Port            string
		User            string
		Password        string
		DBName          string
		SSLMode         string        `mapstructure:"sslMode"`
		MaxOpenConns    int           `mapstructure:"maxOpenConns"`
		MaxIdleConns    int           `mapstructure:"maxIdleConns"`
		ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func newCfg() Config {
	cfg := Config{
		App:      AppConfig{},
		Limiter:  LimiterConfig{},
		Postgres: PostgresConfig{},
		HTTP:     HTTPConfig{},
	}
	return cfg
}

// Init loads config from file and environment variables.
// Priority: defaults -> config file -> env vars
func Init(configfile, envfile string) (*Config, error) {
	populateDefault()

	if err := parseConfigFile(configfile); err != nil {
		return nil, err
	}

	cfg := newCfg()

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := setFromEnv(envfile, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func populateDefault() {
	// http config
	viper.SetDefault("http.host", defaultHTTPHost)
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderMegabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPReadTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPWriteTimeout)

	// limiter config
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)

	// postgres config
	viper.SetDefault("postgres.sslMode", defaultSSLMode)
	viper.SetDefault("postgres.maxOpenConns", defaultMaxOpenConns)
	viper.SetDefault("postgres.maxIdleConns", defaultMaxIdleConns)
	viper.SetDefault("postgres.connMaxLifetime", defaultConnMaxLifetime)
}

func parseConfigFile(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// setFromEnv loads sensitive data from environment variables
func setFromEnv(envpath string, cfg *Config) error {
	err := godotenv.Load(envpath)
	if err != nil {
		return err
	}

	cfg.App.Environment = os.Getenv("ENV")
	cfg.App.AppSecretKey = os.Getenv("APP_SECRET")

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	return nil
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.App.Environment),
		slog.Group("http",
			slog.String("http_address", c.HTTP.Host+":"+c.HTTP.Port),
			slog.Duration("read_timeout", c.HTTP.ReadTimeout),
			slog.Duration("write_timeout", c.HTTP.WriteTimeout),
			slog.Int("maxHeaderMegabytes", c.HTTP.MaxHeaderMegabytes),
		),
		slog.Group("postgres",
			slog.String("host", c.Postgres.Host),
			slog.String("port", c.Postgres.Port),
			slog.String("db", c.Postgres.DBName),
			slog.Int("max_conns", c.Postgres.MaxOpenConns),
		),
		slog.Group("limiter",
			slog.Int("rps", c.Limiter.RPS),
			slog.Int("burst", c.Limiter.Burst),
			slog.Duration("ttl", c.Limiter.TTL),
		),
	)
}
