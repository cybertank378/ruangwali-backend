package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
	shareddomain "github.com/ruangwali/internal/shared/domain"
)

var (
	ErrInvalidSessionID = fmt.Errorf(
		"session ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidSessionUserID = fmt.Errorf(
		"session user ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidSessionFamilyID = fmt.Errorf(
		"session family ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidTokenHash = fmt.Errorf(
		"token hash tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrSessionExpired = fmt.Errorf(
		"session telah kedaluwarsa: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrSessionRevoked = fmt.Errorf(
		"session telah dicabut: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrSessionAlreadyReplaced = fmt.Errorf(
		"session telah diganti: %w",
		shareddomain.ErrConflict,
	)

	ErrInvalidSessionReplacement = fmt.Errorf(
		"replacement session tidak valid: %w",
		shareddomain.ErrValidation,
	)
)

type RefreshSession struct {
	id uuid.UUID

	userID uuid.UUID

	familyID uuid.UUID

	tokenHash []byte

	fingerprint valueobject.SessionFingerprint

	expiresAt time.Time

	lastUsedAt *time.Time

	revokedAt    *time.Time
	revokeReason *string

	replacedByID *uuid.UUID

	createdAt time.Time
	updatedAt time.Time
}

type RehydrateRefreshSessionParams struct {
	ID uuid.UUID

	UserID uuid.UUID

	FamilyID uuid.UUID

	TokenHash []byte

	Fingerprint valueobject.SessionFingerprint

	ExpiresAt time.Time

	LastUsedAt *time.Time

	RevokedAt    *time.Time
	RevokeReason *string

	ReplacedByID *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRefreshSession(
	userID uuid.UUID,
	familyID uuid.UUID,
	tokenHash []byte,
	fingerprint valueobject.SessionFingerprint,
	expiresAt time.Time,
	now time.Time,
) (*RefreshSession, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidSessionUserID
	}

	if familyID == uuid.Nil {
		return nil, ErrInvalidSessionFamilyID
	}

	if len(tokenHash) == 0 {
		return nil, ErrInvalidTokenHash
	}

	now = now.UTC()
	expiresAt = expiresAt.UTC()

	if !expiresAt.After(now) {
		return nil, ErrSessionExpired
	}

	return &RefreshSession{
		id: uuid.New(),

		userID: userID,

		familyID: familyID,

		tokenHash: cloneBytes(
			tokenHash,
		),

		fingerprint: fingerprint,

		expiresAt: expiresAt,

		createdAt: now,
		updatedAt: now,
	}, nil
}

func RehydrateRefreshSession(
	params RehydrateRefreshSessionParams,
) (*RefreshSession, error) {
	if params.ID == uuid.Nil {
		return nil, ErrInvalidSessionID
	}

	if params.UserID == uuid.Nil {
		return nil, ErrInvalidSessionUserID
	}

	if params.FamilyID == uuid.Nil {
		return nil, ErrInvalidSessionFamilyID
	}

	if len(params.TokenHash) == 0 {
		return nil, ErrInvalidTokenHash
	}

	if params.CreatedAt.IsZero() {
		return nil, fmt.Errorf(
			"createdAt wajib tersedia: %w",
			shareddomain.ErrValidation,
		)
	}

	if params.UpdatedAt.IsZero() {
		return nil, fmt.Errorf(
			"updatedAt wajib tersedia: %w",
			shareddomain.ErrValidation,
		)
	}

	if !params.ExpiresAt.After(
		params.CreatedAt,
	) {
		return nil, fmt.Errorf(
			"expiresAt harus setelah createdAt: %w",
			shareddomain.ErrValidation,
		)
	}

	if params.UpdatedAt.Before(
		params.CreatedAt,
	) {
		return nil, fmt.Errorf(
			"updatedAt tidak boleh sebelum createdAt: %w",
			shareddomain.ErrValidation,
		)
	}

	if params.ReplacedByID != nil &&
		*params.ReplacedByID == params.ID {
		return nil, ErrInvalidSessionReplacement
	}

	return &RefreshSession{
		id: params.ID,

		userID: params.UserID,

		familyID: params.FamilyID,

		tokenHash: cloneBytes(
			params.TokenHash,
		),

		fingerprint: params.Fingerprint,

		expiresAt: params.ExpiresAt.UTC(),

		lastUsedAt: cloneTimePointer(
			params.LastUsedAt,
		),

		revokedAt: cloneTimePointer(
			params.RevokedAt,
		),

		revokeReason: cloneStringPointer(
			params.RevokeReason,
		),

		replacedByID: cloneUUIDPointer(
			params.ReplacedByID,
		),

		createdAt: params.CreatedAt.UTC(),
		updatedAt: params.UpdatedAt.UTC(),
	}, nil
}

func (s *RefreshSession) EnsureUsable(
	now time.Time,
) error {
	if s.revokedAt != nil {
		return ErrSessionRevoked
	}

	if !now.UTC().Before(s.expiresAt) {
		return ErrSessionExpired
	}

	return nil
}

func (s *RefreshSession) RecordUse(
	at time.Time,
) error {
	if err := s.EnsureUsable(at); err != nil {
		return err
	}

	at = at.UTC()

	s.lastUsedAt = &at
	s.touch(at)

	return nil
}

func (s *RefreshSession) Revoke(
	reason string,
	at time.Time,
) {
	if s.revokedAt != nil {
		return
	}

	at = at.UTC()

	normalizedReason := strings.TrimSpace(
		reason,
	)

	s.revokedAt = &at

	if normalizedReason != "" {
		s.revokeReason = &normalizedReason
	}

	s.touch(at)
}

func (s *RefreshSession) ReplaceWith(
	replacementID uuid.UUID,
	at time.Time,
) error {
	if replacementID == uuid.Nil ||
		replacementID == s.id {
		return ErrInvalidSessionReplacement
	}

	if s.replacedByID != nil {
		return ErrSessionAlreadyReplaced
	}

	if err := s.EnsureUsable(at); err != nil {
		return err
	}

	replacement := replacementID

	s.replacedByID = &replacement

	s.Revoke(
		"ROTATED",
		at,
	)

	return nil
}

func (s *RefreshSession) IsExpired(
	now time.Time,
) bool {
	return !now.UTC().Before(
		s.expiresAt,
	)
}

func (s *RefreshSession) IsRevoked() bool {
	return s.revokedAt != nil
}

func (s *RefreshSession) ID() uuid.UUID {
	return s.id
}

func (s *RefreshSession) UserID() uuid.UUID {
	return s.userID
}

func (s *RefreshSession) FamilyID() uuid.UUID {
	return s.familyID
}

func (s *RefreshSession) TokenHash() []byte {
	return cloneBytes(
		s.tokenHash,
	)
}

func (s *RefreshSession) Fingerprint() valueobject.SessionFingerprint {
	return s.fingerprint
}

func (s *RefreshSession) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *RefreshSession) LastUsedAt() *time.Time {
	return cloneTimePointer(
		s.lastUsedAt,
	)
}

func (s *RefreshSession) RevokedAt() *time.Time {
	return cloneTimePointer(
		s.revokedAt,
	)
}

func (s *RefreshSession) RevokeReason() *string {
	return cloneStringPointer(
		s.revokeReason,
	)
}

func (s *RefreshSession) ReplacedByID() *uuid.UUID {
	return cloneUUIDPointer(
		s.replacedByID,
	)
}

func (s *RefreshSession) CreatedAt() time.Time {
	return s.createdAt
}

func (s *RefreshSession) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *RefreshSession) touch(
	now time.Time,
) {
	now = now.UTC()

	if now.Before(s.updatedAt) {
		return
	}

	s.updatedAt = now
}

func cloneBytes(
	value []byte,
) []byte {
	if value == nil {
		return nil
	}

	cloned := make(
		[]byte,
		len(value),
	)

	copy(
		cloned,
		value,
	)

	return cloned
}

func cloneStringPointer(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	cloned := *value

	return &cloned
}

func cloneUUIDPointer(
	value *uuid.UUID,
) *uuid.UUID {
	if value == nil {
		return nil
	}

	cloned := *value

	return &cloned
}
