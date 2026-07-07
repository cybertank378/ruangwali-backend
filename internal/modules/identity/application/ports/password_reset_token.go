package ports

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	Raw string

	Hash []byte
}

type PasswordResetTokenService interface {
	Generate(
		ctx context.Context,
	) (PasswordResetToken, error)

	Hash(
		ctx context.Context,
		raw string,
	) ([]byte, error)

	TTL() time.Duration
}
