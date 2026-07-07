package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/ruangwali/internal/modules/identity/application/ports"
)

const refreshTokenBytes = 32

var ErrInvalidRefreshToken = errors.New(
	"refresh token tidak valid",
)

type RefreshTokenService struct{}

func NewRefreshTokenService() *RefreshTokenService {
	return &RefreshTokenService{}
}

func (s *RefreshTokenService) Generate() (
	ports.GeneratedRefreshToken,
	error,
) {
	randomBytes := make(
		[]byte,
		refreshTokenBytes,
	)

	if _, err := rand.Read(
		randomBytes,
	); err != nil {
		return ports.GeneratedRefreshToken{},
			fmt.Errorf(
				"gagal membuat refresh token: %w",
				err,
			)
	}

	raw := base64.RawURLEncoding.EncodeToString(
		randomBytes,
	)

	return ports.GeneratedRefreshToken{
		Raw: raw,

		Hash: s.Hash(raw),
	}, nil
}

func (s *RefreshTokenService) Hash(
	raw string,
) []byte {
	sum := sha256.Sum256(
		[]byte(
			strings.TrimSpace(raw),
		),
	)

	hash := make(
		[]byte,
		len(sum),
	)

	copy(
		hash,
		sum[:],
	)

	return hash
}

func (s *RefreshTokenService) Validate(
	raw string,
) error {
	raw = strings.TrimSpace(
		raw,
	)

	if raw == "" {
		return ErrInvalidRefreshToken
	}

	decoded, err := base64.RawURLEncoding.DecodeString(
		raw,
	)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	if len(decoded) != refreshTokenBytes {
		return ErrInvalidRefreshToken
	}

	return nil
}
