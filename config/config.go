package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"time"
)

var AppConfig *Config

type Config struct {
	Server             Server
	Database           Database
	AuthConfig         AuthConfig
	JWTConfig          JWTConfig
	CtxConfig          CtxConfig
	NotificationConfig NotificationConfig
	Order              Order
	Stripe             Stripe
}

type Stripe struct {
	Secret         string `env:"STRIPE_SECRET"`
	SuccessUrl     string `env:"SUCCESS_URL"`
	CancelUrl      string `env:"CANCEL_URL"`
	PublishableKey string `env:"PUBLISHABLE_KEY"`
}

type Order struct {
	Code int `env:"SERVER_ORDER_CODE"`
}

type Server struct {
	BodyLimit    int           `env:"BODY_LIMIT"`    // 1024 * 1024
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

type CtxConfig struct {
	Timeout time.Duration `env:"CTX_TIMEOUT"`
}

type AuthConfig struct {
	Secret     string        `env:"SECRET"`
	Code       int           `env:"CODE"`
	CodeExpiry time.Duration `env:"CODE_EXPIRY"`
}

type JWTConfig struct {
	TokenExp time.Duration `env:"TOKEN_EXP"`
}

type NotificationConfig struct {
	AccountSMSId string `env:"ACCOUNT_SMS_ID"`
	AuthToken    string `env:"AUTH_TOKEN"`
	SetTo        string `env:"SET_TO"`
	SetFrom      string `env:"SET_FROM"`
}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		slg.Logger.Error("error loading config", "error", err)
		return err
	}

	AppConfig = config

	return nil
}
