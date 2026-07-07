package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

var (
	ErrInvalidPasswordResetID = fmt.Errorf(
		"password reset ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetUserID = fmt.Errorf(
		"password reset user ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetTokenHash = fmt.Errorf(
		"password reset token hash tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetExpiresAt = fmt.Errorf(
		"waktu kedaluwarsa password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetCreatedAt = fmt.Errorf(
		"waktu pembuatan password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetUpdatedAt = fmt.Errorf(
		"waktu pembaruan password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetUsedAt = fmt.Errorf(
		"waktu penggunaan password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetRevokedAt = fmt.Errorf(
		"waktu pencabutan password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordResetReplacement = fmt.Errorf(
		"replacement password reset tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrPasswordResetExpired = fmt.Errorf(
		"password reset telah kedaluwarsa: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrPasswordResetUsed = fmt.Errorf(
		"password reset telah digunakan: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrPasswordResetRevoked = fmt.Errorf(
		"password reset telah dicabut: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrPasswordResetTerminalState = fmt.Errorf(
		"password reset telah mencapai terminal state: %w",
		shareddomain.ErrConflict,
	)
)

type PasswordReset struct {
	id uuid.UUID

	userID uuid.UUID

	tokenHash []byte

	expiresAt time.Time

	usedAt *time.Time

	revokedAt *time.Time

	replacedByID *uuid.UUID

	createdAt time.Time

	updatedAt time.Time
}

type RehydratePasswordResetParams struct {
	ID uuid.UUID

	UserID uuid.UUID

	TokenHash []byte

	ExpiresAt time.Time

	UsedAt *time.Time

	RevokedAt *time.Time

	ReplacedByID *uuid.UUID

	CreatedAt time.Time

	UpdatedAt time.Time
}

func NewPasswordReset(
	userID uuid.UUID,
	tokenHash []byte,
	expiresAt time.Time,
	now time.Time,
) (*PasswordReset, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidPasswordResetUserID
	}

	if len(tokenHash) == 0 {
		return nil, ErrInvalidPasswordResetTokenHash
	}

	if now.IsZero() {
		return nil, ErrInvalidPasswordResetCreatedAt
	}

	now = now.UTC()
	expiresAt = expiresAt.UTC()

	if expiresAt.IsZero() ||
		!expiresAt.After(now) {
		return nil, ErrInvalidPasswordResetExpiresAt
	}

	return &PasswordReset{
		id: uuid.New(),

		userID: userID,

		tokenHash: clonePasswordResetBytes(
			tokenHash,
		),

		expiresAt: expiresAt,

		usedAt: nil,

		revokedAt: nil,

		replacedByID: nil,

		createdAt: now,

		updatedAt: now,
	}, nil
}

func RehydratePasswordReset(
	params RehydratePasswordResetParams,
) (*PasswordReset, error) {
	if err := validatePasswordResetRehydration(
		params,
	); err != nil {
		return nil, err
	}

	return &PasswordReset{
		id: params.ID,

		userID: params.UserID,

		tokenHash: clonePasswordResetBytes(
			params.TokenHash,
		),

		expiresAt: params.ExpiresAt.UTC(),

		usedAt: clonePasswordResetTimePointer(
			params.UsedAt,
		),

		revokedAt: clonePasswordResetTimePointer(
			params.RevokedAt,
		),

		replacedByID: clonePasswordResetUUIDPointer(
			params.ReplacedByID,
		),

		createdAt: params.CreatedAt.UTC(),

		updatedAt: params.UpdatedAt.UTC(),
	}, nil
}

func (r *PasswordReset) EnsureUsable(
	now time.Time,
) error {
	if r == nil {
		return ErrInvalidPasswordResetID
	}

	if r.usedAt != nil {
		return ErrPasswordResetUsed
	}

	if r.revokedAt != nil {
		return ErrPasswordResetRevoked
	}

	if now.IsZero() {
		return ErrInvalidPasswordResetUpdatedAt
	}

	if !now.UTC().Before(
		r.expiresAt,
	) {
		return ErrPasswordResetExpired
	}

	return nil
}

func (r *PasswordReset) MarkUsed(
	at time.Time,
) error {
	if r == nil {
		return ErrInvalidPasswordResetID
	}

	if err := r.EnsureUsable(at); err != nil {
		return err
	}

	at = at.UTC()

	if at.Before(
		r.createdAt,
	) {
		return ErrInvalidPasswordResetUsedAt
	}

	r.usedAt = &at
	r.updatedAt = at

	return nil
}

func (r *PasswordReset) Revoke(
	at time.Time,
) error {
	if r == nil {
		return ErrInvalidPasswordResetID
	}

	if r.usedAt != nil ||
		r.revokedAt != nil {
		return ErrPasswordResetTerminalState
	}

	if at.IsZero() {
		return ErrInvalidPasswordResetRevokedAt
	}

	at = at.UTC()

	if at.Before(
		r.createdAt,
	) {
		return ErrInvalidPasswordResetRevokedAt
	}

	r.revokedAt = &at
	r.updatedAt = at

	return nil
}

func (r *PasswordReset) ReplaceWith(
	replacementID uuid.UUID,
	at time.Time,
) error {
	if r == nil {
		return ErrInvalidPasswordResetID
	}

	if replacementID == uuid.Nil ||
		replacementID == r.id {
		return ErrInvalidPasswordResetReplacement
	}

	if err := r.EnsureUsable(at); err != nil {
		return err
	}

	at = at.UTC()

	if at.Before(
		r.createdAt,
	) {
		return ErrInvalidPasswordResetRevokedAt
	}

	r.revokedAt = &at

	r.replacedByID = clonePasswordResetUUIDPointer(
		&replacementID,
	)

	r.updatedAt = at

	return nil
}

func (r *PasswordReset) IsUsable(
	now time.Time,
) bool {
	return r.EnsureUsable(now) == nil
}

func (r *PasswordReset) IsUsed() bool {
	return r != nil &&
		r.usedAt != nil
}

func (r *PasswordReset) IsRevoked() bool {
	return r != nil &&
		r.revokedAt != nil
}

func (r *PasswordReset) IsExpired(
	now time.Time,
) bool {
	if r == nil ||
		now.IsZero() {
		return false
	}

	return !now.UTC().Before(
		r.expiresAt,
	)
}

func (r *PasswordReset) IsReplaced() bool {
	return r != nil &&
		r.replacedByID != nil
}

func (r *PasswordReset) ID() uuid.UUID {
	if r == nil {
		return uuid.Nil
	}

	return r.id
}

func (r *PasswordReset) UserID() uuid.UUID {
	if r == nil {
		return uuid.Nil
	}

	return r.userID
}

func (r *PasswordReset) TokenHash() []byte {
	if r == nil {
		return nil
	}

	return clonePasswordResetBytes(
		r.tokenHash,
	)
}

func (r *PasswordReset) ExpiresAt() time.Time {
	if r == nil {
		return time.Time{}
	}

	return r.expiresAt
}

func (r *PasswordReset) UsedAt() *time.Time {
	if r == nil {
		return nil
	}

	return clonePasswordResetTimePointer(
		r.usedAt,
	)
}

func (r *PasswordReset) RevokedAt() *time.Time {
	if r == nil {
		return nil
	}

	return clonePasswordResetTimePointer(
		r.revokedAt,
	)
}

func (r *PasswordReset) ReplacedByID() *uuid.UUID {
	if r == nil {
		return nil
	}

	return clonePasswordResetUUIDPointer(
		r.replacedByID,
	)
}

func (r *PasswordReset) CreatedAt() time.Time {
	if r == nil {
		return time.Time{}
	}

	return r.createdAt
}

func (r *PasswordReset) UpdatedAt() time.Time {
	if r == nil {
		return time.Time{}
	}

	return r.updatedAt
}

func validatePasswordResetRehydration(
	params RehydratePasswordResetParams,
) error {
	if params.ID == uuid.Nil {
		return ErrInvalidPasswordResetID
	}

	if params.UserID == uuid.Nil {
		return ErrInvalidPasswordResetUserID
	}

	if len(params.TokenHash) == 0 {
		return ErrInvalidPasswordResetTokenHash
	}

	if params.CreatedAt.IsZero() {
		return ErrInvalidPasswordResetCreatedAt
	}

	if params.UpdatedAt.IsZero() {
		return ErrInvalidPasswordResetUpdatedAt
	}

	createdAt := params.CreatedAt.UTC()
	updatedAt := params.UpdatedAt.UTC()
	expiresAt := params.ExpiresAt.UTC()

	if params.ExpiresAt.IsZero() ||
		!expiresAt.After(createdAt) {
		return ErrInvalidPasswordResetExpiresAt
	}

	if updatedAt.Before(createdAt) {
		return ErrInvalidPasswordResetUpdatedAt
	}

	if params.UsedAt != nil {
		usedAt := params.UsedAt.UTC()

		if usedAt.Before(createdAt) {
			return ErrInvalidPasswordResetUsedAt
		}

		if usedAt.After(updatedAt) {
			return ErrInvalidPasswordResetUsedAt
		}
	}

	if params.RevokedAt != nil {
		revokedAt := params.RevokedAt.UTC()

		if revokedAt.Before(createdAt) {
			return ErrInvalidPasswordResetRevokedAt
		}

		if revokedAt.After(updatedAt) {
			return ErrInvalidPasswordResetRevokedAt
		}
	}

	if params.UsedAt != nil &&
		params.RevokedAt != nil {
		return ErrPasswordResetTerminalState
	}

	if params.ReplacedByID != nil {
		if *params.ReplacedByID == uuid.Nil ||
			*params.ReplacedByID == params.ID {
			return ErrInvalidPasswordResetReplacement
		}

		if params.RevokedAt == nil {
			return ErrInvalidPasswordResetReplacement
		}
	}

	return nil
}

func clonePasswordResetBytes(
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

func clonePasswordResetTimePointer(
	value *time.Time,
) *time.Time {
	if value == nil {
		return nil
	}

	return new(value.UTC())
}

func clonePasswordResetUUIDPointer(
	value *uuid.UUID,
) *uuid.UUID {
	if value == nil {
		return nil
	}

	return new(*value)
}
