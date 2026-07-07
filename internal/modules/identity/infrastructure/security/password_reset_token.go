// =========================================================
// File: internal/modules/identity/infrastructure/security/password_reset_token.go.go
// =========================================================

package security

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/ports"
)

const passwordResetTokenBytes = 32

var ErrInvalidPasswordResetToken = errors.New(
	"password reset token tidak valid",
)

type PasswordResetTokenService struct {
	ttl time.Duration
}

func NewPasswordResetTokenService(
	ttl time.Duration,
) (*PasswordResetTokenService, error) {
	if ttl <= 0 {
		return nil, errors.New(
			"password reset token TTL harus lebih besar dari 0",
		)
	}

	return &PasswordResetTokenService{
		ttl: ttl,
	}, nil
}

func (s *PasswordResetTokenService) Generate(
	ctx context.Context,
) (ports.PasswordResetToken, error) {
	if err := ctx.Err(); err != nil {
		return ports.PasswordResetToken{}, err
	}

	randomBytes := make(
		[]byte,
		passwordResetTokenBytes,
	)

	if _, err := rand.Read(
		randomBytes,
	); err != nil {
		return ports.PasswordResetToken{},
			fmt.Errorf(
				"gagal membuat password reset token: %w",
				err,
			)
	}

	raw := base64.RawURLEncoding.EncodeToString(
		randomBytes,
	)

	hash, err := s.Hash(
		ctx,
		raw,
	)
	if err != nil {
		return ports.PasswordResetToken{}, err
	}

	return ports.PasswordResetToken{
		Raw:  raw,
		Hash: hash,
	}, nil
}

func (s *PasswordResetTokenService) Hash(
	ctx context.Context,
	raw string,
) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	raw = strings.TrimSpace(
		raw,
	)
	if raw == "" {
		return nil, ErrInvalidPasswordResetToken
	}

	decoded, err := base64.RawURLEncoding.DecodeString(
		raw,
	)
	if err != nil {
		return nil, ErrInvalidPasswordResetToken
	}

	if len(decoded) != passwordResetTokenBytes {
		return nil, ErrInvalidPasswordResetToken
	}

	sum := sha256.Sum256(
		[]byte(raw),
	)

	hash := make(
		[]byte,
		len(sum),
	)

	copy(
		hash,
		sum[:],
	)

	return hash, nil
}

func (s *PasswordResetTokenService) TTL() time.Duration {
	return s.ttl
}

var _ ports.PasswordResetTokenService = (*PasswordResetTokenService)(nil)
