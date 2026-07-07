package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAppEnv = "development"

	defaultHTTPAddr = ":8080"

	defaultJWTIssuer   = "ruangwali-api"
	defaultJWTAudience = "ruangwali-web"

	defaultAccessTokenTTLMinutes        = 15
	defaultRefreshTokenTTLDays          = 30
	defaultPasswordResetTokenTTLMinutes = 30

	defaultDatabaseMaxConns = 20
	defaultDatabaseMinConns = 2

	defaultDatabaseMaxConnLifetime = 30 * time.Minute
	defaultDatabaseMaxConnIdleTime = 5 * time.Minute
	defaultDatabaseHealthTimeout   = 5 * time.Second
)

type Config struct {
	App         AppConfig
	HTTP        HTTPConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Integration IntegrationConfig
}

type AppConfig struct {
	Env string
}

type HTTPConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL string

	MaxConns int32
	MinConns int32

	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	HealthTimeout time.Duration
}

type AuthConfig struct {
	JWTIssuer   string
	JWTAudience string
	JWTSecret   string

	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	PasswordResetTokenTTL time.Duration
}

type IntegrationConfig struct {
	GoogleAppsScript GoogleAppsScriptConfig
}

type GoogleAppsScriptConfig struct {
	BaseURL       string
	WebhookSecret string
}

func Load() (Config, error) {
	databaseMaxConns, err := envInt(
		"DATABASE_MAX_CONNS",
		defaultDatabaseMaxConns,
	)
	if err != nil {
		return Config{}, err
	}

	databaseMinConns, err := envInt(
		"DATABASE_MIN_CONNS",
		defaultDatabaseMinConns,
	)
	if err != nil {
		return Config{}, err
	}

	databaseMaxConnLifetime, err := envDuration(
		"DATABASE_MAX_CONN_LIFETIME",
		defaultDatabaseMaxConnLifetime,
	)
	if err != nil {
		return Config{}, err
	}

	databaseMaxConnIdleTime, err := envDuration(
		"DATABASE_MAX_CONN_IDLE_TIME",
		defaultDatabaseMaxConnIdleTime,
	)
	if err != nil {
		return Config{}, err
	}

	databaseHealthTimeout, err := envDuration(
		"DATABASE_HEALTH_TIMEOUT",
		defaultDatabaseHealthTimeout,
	)
	if err != nil {
		return Config{}, err
	}

	accessTokenTTLMinutes, err := envInt(
		"ACCESS_TOKEN_TTL_MINUTES",
		defaultAccessTokenTTLMinutes,
	)
	if err != nil {
		return Config{}, err
	}

	refreshTokenTTLDays, err := envInt(
		"REFRESH_TOKEN_TTL_DAYS",
		defaultRefreshTokenTTLDays,
	)
	if err != nil {
		return Config{}, err
	}

	passwordResetTokenTTLMinutes, err := envInt(
		"PASSWORD_RESET_TOKEN_TTL_MINUTES",
		defaultPasswordResetTokenTTLMinutes,
	)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App: AppConfig{
			Env: env(
				"APP_ENV",
				defaultAppEnv,
			),
		},

		HTTP: HTTPConfig{
			Addr: env(
				"HTTP_ADDR",
				defaultHTTPAddr,
			),
		},

		Database: DatabaseConfig{
			URL: strings.TrimSpace(
				os.Getenv(
					"DATABASE_URL",
				),
			),

			MaxConns: int32(
				databaseMaxConns,
			),

			MinConns: int32(
				databaseMinConns,
			),

			MaxConnLifetime: databaseMaxConnLifetime,

			MaxConnIdleTime: databaseMaxConnIdleTime,

			HealthTimeout: databaseHealthTimeout,
		},

		Auth: AuthConfig{
			JWTIssuer: env(
				"JWT_ISSUER",
				defaultJWTIssuer,
			),

			JWTAudience: env(
				"JWT_AUDIENCE",
				defaultJWTAudience,
			),

			JWTSecret: strings.TrimSpace(
				os.Getenv(
					"JWT_SECRET",
				),
			),

			AccessTokenTTL: time.Duration(
				accessTokenTTLMinutes,
			) * time.Minute,

			RefreshTokenTTL: time.Duration(
				refreshTokenTTLDays,
			) * 24 * time.Hour,

			PasswordResetTokenTTL: time.Duration(
				passwordResetTokenTTLMinutes,
			) * time.Minute,
		},

		Integration: IntegrationConfig{
			GoogleAppsScript: GoogleAppsScriptConfig{
				BaseURL: strings.TrimSpace(
					os.Getenv(
						"GOOGLE_APPS_SCRIPT_BASE_URL",
					),
				),

				WebhookSecret: strings.TrimSpace(
					os.Getenv(
						"GOOGLE_APPS_SCRIPT_WEBHOOK_SECRET",
					),
				),
			},
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}

	if err := c.HTTP.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	if err := c.Auth.Validate(); err != nil {
		return err
	}

	if err := c.Integration.Validate(); err != nil {
		return err
	}

	return nil
}

func (c AppConfig) Validate() error {
	switch c.Env {
	case
		"development",
		"test",
		"staging",
		"production":
		return nil

	default:
		return fmt.Errorf(
			"APP_ENV tidak valid: %q",
			c.Env,
		)
	}
}

func (c AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}

func (c AppConfig) IsTest() bool {
	return c.Env == "test"
}

func (c AppConfig) IsStaging() bool {
	return c.Env == "staging"
}

func (c AppConfig) IsProduction() bool {
	return c.Env == "production"
}

func (c HTTPConfig) Validate() error {
	if strings.TrimSpace(
		c.Addr,
	) == "" {
		return errors.New(
			"HTTP_ADDR wajib diisi",
		)
	}

	return nil
}

func (c DatabaseConfig) Validate() error {
	if strings.TrimSpace(
		c.URL,
	) == "" {
		return errors.New(
			"DATABASE_URL wajib diisi",
		)
	}

	if c.MaxConns < 1 {
		return errors.New(
			"DATABASE_MAX_CONNS minimal 1",
		)
	}

	if c.MinConns < 0 {
		return errors.New(
			"DATABASE_MIN_CONNS tidak boleh negatif",
		)
	}

	if c.MinConns > c.MaxConns {
		return errors.New(
			"DATABASE_MIN_CONNS tidak boleh melebihi DATABASE_MAX_CONNS",
		)
	}

	if c.MaxConnLifetime <= 0 {
		return errors.New(
			"DATABASE_MAX_CONN_LIFETIME harus lebih besar dari 0",
		)
	}

	if c.MaxConnIdleTime <= 0 {
		return errors.New(
			"DATABASE_MAX_CONN_IDLE_TIME harus lebih besar dari 0",
		)
	}

	if c.HealthTimeout <= 0 {
		return errors.New(
			"DATABASE_HEALTH_TIMEOUT harus lebih besar dari 0",
		)
	}

	return nil
}

func (c AuthConfig) Validate() error {
	if strings.TrimSpace(
		c.JWTIssuer,
	) == "" {
		return errors.New(
			"JWT_ISSUER wajib diisi",
		)
	}

	if strings.TrimSpace(
		c.JWTAudience,
	) == "" {
		return errors.New(
			"JWT_AUDIENCE wajib diisi",
		)
	}

	if len(
		c.JWTSecret,
	) < 32 {
		return errors.New(
			"JWT_SECRET minimal 32 karakter",
		)
	}

	if c.AccessTokenTTL <= 0 {
		return errors.New(
			"ACCESS_TOKEN_TTL_MINUTES harus lebih besar dari 0",
		)
	}

	if c.RefreshTokenTTL <= 0 {
		return errors.New(
			"REFRESH_TOKEN_TTL_DAYS harus lebih besar dari 0",
		)
	}

	if c.PasswordResetTokenTTL <= 0 {
		return errors.New(
			"PASSWORD_RESET_TOKEN_TTL_MINUTES harus lebih besar dari 0",
		)
	}

	return nil
}

func (c IntegrationConfig) Validate() error {
	return c.GoogleAppsScript.Validate()
}

func (c GoogleAppsScriptConfig) Validate() error {
	baseURL := strings.TrimSpace(
		c.BaseURL,
	)

	webhookSecret := strings.TrimSpace(
		c.WebhookSecret,
	)

	if baseURL == "" &&
		webhookSecret == "" {
		return nil
	}

	if baseURL == "" {
		return errors.New(
			"GOOGLE_APPS_SCRIPT_BASE_URL wajib diisi ketika integrasi Google Apps Script digunakan",
		)
	}

	if webhookSecret == "" {
		return errors.New(
			"GOOGLE_APPS_SCRIPT_WEBHOOK_SECRET wajib diisi ketika integrasi Google Apps Script digunakan",
		)
	}

	return nil
}

func env(
	key string,
	fallback string,
) string {
	value := strings.TrimSpace(
		os.Getenv(
			key,
		),
	)

	if value == "" {
		return fallback
	}

	return value
}

func envInt(
	key string,
	fallback int,
) (
	int,
	error,
) {
	raw := strings.TrimSpace(
		os.Getenv(
			key,
		),
	)

	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(
		raw,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"%s tidak valid: %w",
			key,
			err,
		)
	}

	return value, nil
}

func envDuration(
	key string,
	fallback time.Duration,
) (
	time.Duration,
	error,
) {
	raw := strings.TrimSpace(
		os.Getenv(
			key,
		),
	)

	if raw == "" {
		return fallback, nil
	}

	value, err := time.ParseDuration(
		raw,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"%s tidak valid: %w",
			key,
			err,
		)
	}

	return value, nil
}
