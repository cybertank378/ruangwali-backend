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
	defaultHTTPAddr              = ":8080"
	defaultJWTIssuer             = "ruangwali-api"
	defaultJWTAudience           = "ruangwali-web"
	defaultAccessTokenTTLMinutes = 30
	defaultRefreshTokenTTLHours  = 168
	minJWTSecretLength           = 32
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string

	JWTIssuer   string
	JWTAudience string
	JWTSecret   string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() (Config, error) {
	accessTokenTTLMinutes, err := envInt(
		"ACCESS_TOKEN_TTL_MINUTES",
		defaultAccessTokenTTLMinutes,
	)
	if err != nil {
		return Config{}, err
	}

	refreshTokenTTLHours, err := envInt(
		"REFRESH_TOKEN_TTL_HOURS",
		defaultRefreshTokenTTLHours,
	)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTPAddr: strings.TrimSpace(
			env(
				"HTTP_ADDR",
				defaultHTTPAddr,
			),
		),

		DatabaseURL: strings.TrimSpace(
			os.Getenv("DATABASE_URL"),
		),

		JWTIssuer: strings.TrimSpace(
			env(
				"JWT_ISSUER",
				defaultJWTIssuer,
			),
		),

		JWTAudience: strings.TrimSpace(
			env(
				"JWT_AUDIENCE",
				defaultJWTAudience,
			),
		),

		JWTSecret: strings.TrimSpace(
			os.Getenv("JWT_SECRET"),
		),

		AccessTokenTTL: time.Duration(
			accessTokenTTLMinutes,
		) * time.Minute,

		RefreshTokenTTL: time.Duration(
			refreshTokenTTLHours,
		) * time.Hour,
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.HTTPAddr == "" {
		return errors.New(
			"HTTP_ADDR wajib diisi",
		)
	}

	if c.DatabaseURL == "" {
		return errors.New(
			"DATABASE_URL wajib diisi",
		)
	}

	if c.JWTIssuer == "" {
		return errors.New(
			"JWT_ISSUER wajib diisi",
		)
	}

	if c.JWTAudience == "" {
		return errors.New(
			"JWT_AUDIENCE wajib diisi",
		)
	}

	if c.JWTSecret == "" {
		return errors.New(
			"JWT_SECRET wajib diisi",
		)
	}

	if len(c.JWTSecret) < minJWTSecretLength {
		return fmt.Errorf(
			"JWT_SECRET minimal %d karakter",
			minJWTSecretLength,
		)
	}

	if c.AccessTokenTTL <= 0 {
		return errors.New(
			"ACCESS_TOKEN_TTL_MINUTES harus lebih besar dari 0",
		)
	}

	if c.RefreshTokenTTL <= 0 {
		return errors.New(
			"REFRESH_TOKEN_TTL_HOURS harus lebih besar dari 0",
		)
	}

	if c.RefreshTokenTTL <= c.AccessTokenTTL {
		return errors.New(
			"REFRESH_TOKEN_TTL_HOURS harus lebih lama dari ACCESS_TOKEN_TTL_MINUTES",
		)
	}

	return nil
}

func env(
	key string,
	fallback string,
) string {
	value := strings.TrimSpace(
		os.Getenv(key),
	)

	if value != "" {
		return value
	}

	return fallback
}

func envInt(
	key string,
	fallback int,
) (int, error) {
	rawValue := strings.TrimSpace(
		os.Getenv(key),
	)

	if rawValue == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, fmt.Errorf(
			"%s tidak valid: harus berupa integer",
			key,
		)
	}

	if value <= 0 {
		return 0, fmt.Errorf(
			"%s harus lebih besar dari 0",
			key,
		)
	}

	return value, nil
}
