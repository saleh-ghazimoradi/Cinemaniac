package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

var AppConfig *Config

type Config struct {
	Server    Server
	Database  Database
	CTX       CTX
	RateLimit RateLimit
	SMTP      SMTP
}

type Server struct {
	BodyLimit    int64         `env:"BODY_LIMIT"` // 1024 * 1024
	Env          string        `env:"ENVIRONMENT"`
	Version      string        `env:"VERSION"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"` // 10s
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`  // 5s
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT"`  // 30s
	RateLimit    int           `env:"RATE_LIMIT"`    // 100
	RateLimitExp time.Duration `env:"RATE_EXP"`      // 60s
	Port         string        `env:"PORT"`          // 3000
	Timeout      time.Duration `env:"TIMEOUT"`       // 30s
}

type Database struct {
	DatabaseHost     string        `env:"DATABASE_HOST"`
	DatabasePort     string        `env:"DATABASE_PORT"`
	DatabaseUser     string        `env:"DATABASE_USER"`
	DatabasePassword string        `env:"DATABASE_PASSWORD"`
	DatabaseName     string        `env:"DATABASE_NAME"`
	DatabaseSSLMode  string        `env:"DATABASE_SSLMODE"`
	MaxOpenConn      int           `env:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConn      int           `env:"DB_MAX_IDLE_CONNECTIONS"`
	MaxLifetime      time.Duration `env:"DB_MAX_LIFETIME"`
	MaxIdleTime      time.Duration `env:"DB_MAX_IDLE_TIME"`
	Timeout          time.Duration `env:"DB_TIMEOUT"`
}

type CTX struct {
	Timeout time.Duration `env:"CTX_TIMEOUT"`
}

type RateLimit struct {
	RPS     float64 `env:"RATE_LIMIT_RPS"`
	Burst   int     `env:"RATE_LIMIT_BURST"`
	Enabled bool    `env:"RATE_LIMIT_ENABLED"`
}

type SMTP struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	UserName string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	Sender   string `env:"SMTP_SENDER"`
}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		return err
	}

	AppConfig = config

	return nil
}
