// Package ports FilePath: /internal/modules/identity/application/ports/token_service.go
package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID uuid.UUID
}

type AccessTokenService interface {
	Issue(
		ctx context.Context,
		claims AccessTokenClaims,
	) (string, error)

	Parse(
		ctx context.Context,
		raw string,
	) (AccessTokenClaims, error)

	TTL() time.Duration
}

type GeneratedRefreshToken struct {
	Raw string

	Hash []byte
}

type RefreshTokenService interface {
	Generate() (
		GeneratedRefreshToken,
		error,
	)

	Hash(
		raw string,
	) []byte

	Validate(
		raw string,
	) error
}
