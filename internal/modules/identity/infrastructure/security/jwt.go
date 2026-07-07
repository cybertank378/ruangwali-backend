package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

const (
	tenantIDClaimKey     = "tenant_id"
	membershipIDClaimKey = "membership_id"
)

var (
	ErrInvalidToken = errors.New(
		"token tidak valid",
	)

	ErrMissingSubject = errors.New(
		"subject claim tidak tersedia",
	)

	ErrMissingTenantID = errors.New(
		"tenant_id claim tidak tersedia",
	)

	ErrInvalidTenantID = errors.New(
		"tenant_id claim tidak valid",
	)

	ErrMissingMembershipID = errors.New(
		"membership_id claim tidak tersedia",
	)

	ErrInvalidMembershipID = errors.New(
		"membership_id claim tidak valid",
	)
)

type TokenService struct {
	issuer   string
	audience string
	secret   []byte
	ttl      time.Duration
}

type Claims struct {
	UserID       string
	TenantID     string
	MembershipID string
}

func NewTokenService(
	issuer string,
	audience string,
	secret string,
	ttl time.Duration,
) *TokenService {
	return &TokenService{
		issuer: strings.TrimSpace(
			issuer,
		),
		audience: strings.TrimSpace(
			audience,
		),
		secret: []byte(
			secret,
		),
		ttl: ttl,
	}
}

func (s *TokenService) Issue(
	ctx context.Context,
	claims Claims,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if err := s.validateClaims(claims); err != nil {
		return "", err
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
			strings.TrimSpace(
				claims.UserID,
			),
		).
		IssuedAt(
			now,
		).
		Expiration(
			now.Add(s.ttl),
		).
		Claim(
			tenantIDClaimKey,
			strings.TrimSpace(
				claims.TenantID,
			),
		).
		Claim(
			membershipIDClaimKey,
			strings.TrimSpace(
				claims.MembershipID,
			),
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
) (Claims, error) {
	if err := ctx.Err(); err != nil {
		return Claims{}, err
	}

	raw = strings.TrimSpace(raw)

	if raw == "" {
		return Claims{}, ErrInvalidToken
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
		return Claims{}, fmt.Errorf(
			"%w: %v",
			ErrInvalidToken,
			err,
		)
	}

	return extractClaims(token)
}

func (s *TokenService) validateClaims(
	claims Claims,
) error {
	if strings.TrimSpace(
		claims.UserID,
	) == "" {
		return ErrMissingSubject
	}

	if strings.TrimSpace(
		claims.TenantID,
	) == "" {
		return ErrMissingTenantID
	}

	if strings.TrimSpace(
		claims.MembershipID,
	) == "" {
		return ErrMissingMembershipID
	}

	return nil
}

func extractClaims(
	token jwt.Token,
) (Claims, error) {
	userID, ok := token.Subject()
	if !ok {
		return Claims{}, ErrMissingSubject
	}

	userID = strings.TrimSpace(userID)

	if userID == "" {
		return Claims{}, ErrMissingSubject
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return Claims{}, fmt.Errorf(
			"gagal membaca token claims: %w",
			err,
		)
	}

	var customClaims struct {
		TenantID     string `json:"tenant_id"`
		MembershipID string `json:"membership_id"`
	}

	if err := json.Unmarshal(payload, &customClaims); err != nil {
		return Claims{}, fmt.Errorf(
			"gagal decode token claims: %w",
			err,
		)
	}

	tenantID := strings.TrimSpace(customClaims.TenantID)
	if tenantID == "" {
		return Claims{}, ErrMissingTenantID
	}

	membershipID := strings.TrimSpace(
		customClaims.MembershipID,
	)
	if membershipID == "" {
		return Claims{}, ErrMissingMembershipID
	}

	return Claims{
		UserID:       userID,
		TenantID:     tenantID,
		MembershipID: membershipID,
	}, nil
}

func getStringPrivateClaim(
	privateClaims map[string]any,
	key string,
	missingErr error,
	invalidErr error,
) (string, error) {
	value, ok := privateClaims[key]
	if !ok {
		return "", missingErr
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", invalidErr
	}

	stringValue = strings.TrimSpace(
		stringValue,
	)

	if stringValue == "" {
		return "", invalidErr
	}

	return stringValue, nil
}
