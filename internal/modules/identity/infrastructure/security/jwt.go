package security

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwt"

	"github.com/ruangwali/internal/modules/identity/application/ports"
)

var (
	ErrInvalidAccessToken = errors.New(
		"access token tidak valid",
	)

	ErrMissingSubject = errors.New(
		"subject claim tidak tersedia",
	)

	ErrInvalidSubject = errors.New(
		"subject claim tidak valid",
	)
)

type TokenService struct {
	issuer string

	audience string

	secret []byte

	ttl time.Duration
}

func NewTokenService(
	issuer string,
	audience string,
	secret string,
	ttl time.Duration,
) (*TokenService, error) {
	issuer = strings.TrimSpace(
		issuer,
	)

	audience = strings.TrimSpace(
		audience,
	)

	secret = strings.TrimSpace(
		secret,
	)

	if issuer == "" {
		return nil, errors.New(
			"JWT issuer wajib diisi",
		)
	}

	if audience == "" {
		return nil, errors.New(
			"JWT audience wajib diisi",
		)
	}

	if len(secret) < 32 {
		return nil, errors.New(
			"JWT secret minimal 32 karakter",
		)
	}

	if ttl <= 0 {
		return nil, errors.New(
			"JWT TTL harus lebih besar dari 0",
		)
	}

	return &TokenService{
		issuer: issuer,

		audience: audience,

		secret: []byte(secret),

		ttl: ttl,
	}, nil
}

func (s *TokenService) Issue(
	ctx context.Context,
	claims ports.AccessTokenClaims,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if claims.UserID == uuid.Nil {
		return "", ErrInvalidSubject
	}

	now := time.Now().UTC()

	token, err := jwt.NewBuilder().
		Issuer(
			s.issuer,
		).
		Audience(
			[]string{
				s.audience,
			},
		).
		Subject(
			claims.UserID.String(),
		).
		IssuedAt(
			now,
		).
		Expiration(
			now.Add(s.ttl),
		).
		Build()
	if err != nil {
		return "", fmt.Errorf(
			"gagal membuat access token: %w",
			err,
		)
	}

	signed, err := jwt.Sign(
		token,
		jwt.WithKey(
			jwa.HS256(),
			s.secret,
		),
	)
	if err != nil {
		return "", fmt.Errorf(
			"gagal menandatangani access token: %w",
			err,
		)
	}

	return string(signed), nil
}

func (s *TokenService) Parse(
	ctx context.Context,
	raw string,
) (ports.AccessTokenClaims, error) {
	if err := ctx.Err(); err != nil {
		return ports.AccessTokenClaims{}, err
	}

	raw = strings.TrimSpace(
		raw,
	)

	if raw == "" {
		return ports.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	token, err := jwt.Parse(
		[]byte(raw),
		jwt.WithKey(
			jwa.HS256(),
			s.secret,
		),
		jwt.WithValidate(true),
		jwt.WithIssuer(
			s.issuer,
		),
		jwt.WithAudience(
			s.audience,
		),
	)
	if err != nil {
		return ports.AccessTokenClaims{},
			fmt.Errorf(
				"%w: %v",
				ErrInvalidAccessToken,
				err,
			)
	}

	subject, ok := token.Subject()
	if !ok {
		return ports.AccessTokenClaims{},
			ErrMissingSubject
	}

	subject = strings.TrimSpace(
		subject,
	)

	if subject == "" {
		return ports.AccessTokenClaims{},
			ErrMissingSubject
	}

	userID, err := uuid.Parse(
		subject,
	)
	if err != nil {
		return ports.AccessTokenClaims{},
			fmt.Errorf(
				"%w: %v",
				ErrInvalidSubject,
				err,
			)
	}

	return ports.AccessTokenClaims{
		UserID: userID,
	}, nil
}

func (s *TokenService) TTL() time.Duration {
	return s.ttl
}
