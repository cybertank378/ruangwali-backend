package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
	shareddomain "github.com/ruangwali/internal/shared/domain"
)

var (
	ErrInvalidUserID = fmt.Errorf(
		"user ID tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrInvalidPasswordHash = fmt.Errorf(
		"password hash tidak valid: %w",
		shareddomain.ErrValidation,
	)

	ErrUserInactive = fmt.Errorf(
		"user tidak aktif: %w",
		shareddomain.ErrForbidden,
	)

	ErrUserSuspended = fmt.Errorf(
		"user ditangguhkan: %w",
		shareddomain.ErrForbidden,
	)
)

type User struct {
	id uuid.UUID

	email valueobject.Email

	passwordHash string

	status valueobject.UserStatus

	lastLoginAt *time.Time

	createdAt time.Time
	updatedAt time.Time
}

type RehydrateUserParams struct {
	ID uuid.UUID

	Email valueobject.Email

	PasswordHash string

	Status valueobject.UserStatus

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(
	email valueobject.Email,
	passwordHash string,
	now time.Time,
) (*User, error) {
	passwordHash = strings.TrimSpace(
		passwordHash,
	)

	if email.IsZero() {
		return nil, valueobject.ErrInvalidEmail
	}

	if passwordHash == "" {
		return nil, ErrInvalidPasswordHash
	}

	now = now.UTC()

	return &User{
		id: uuid.New(),

		email: email,

		passwordHash: passwordHash,

		status: valueobject.UserStatusActive,

		createdAt: now,
		updatedAt: now,
	}, nil
}

func RehydrateUser(
	params RehydrateUserParams,
) (*User, error) {
	if params.ID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if params.Email.IsZero() {
		return nil, valueobject.ErrInvalidEmail
	}

	passwordHash := strings.TrimSpace(
		params.PasswordHash,
	)

	if passwordHash == "" {
		return nil, ErrInvalidPasswordHash
	}

	if !params.Status.IsValid() {
		return nil, valueobject.ErrInvalidUserStatus
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

	if params.UpdatedAt.Before(
		params.CreatedAt,
	) {
		return nil, fmt.Errorf(
			"updatedAt tidak boleh sebelum createdAt: %w",
			shareddomain.ErrValidation,
		)
	}

	return &User{
		id: params.ID,

		email: params.Email,

		passwordHash: passwordHash,

		status: params.Status,

		lastLoginAt: cloneTimePointer(
			params.LastLoginAt,
		),

		createdAt: params.CreatedAt.UTC(),
		updatedAt: params.UpdatedAt.UTC(),
	}, nil
}

func (u *User) EnsureCanAuthenticate() error {
	switch u.status {
	case valueobject.UserStatusActive:
		return nil

	case valueobject.UserStatusInactive:
		return ErrUserInactive

	case valueobject.UserStatusSuspended:
		return ErrUserSuspended

	default:
		return valueobject.ErrInvalidUserStatus
	}
}

func (u *User) ChangeEmail(
	email valueobject.Email,
	now time.Time,
) error {
	if email.IsZero() {
		return valueobject.ErrInvalidEmail
	}

	if u.email.Equal(email) {
		return nil
	}

	u.email = email
	u.touch(now)

	return nil
}

func (u *User) ChangePasswordHash(
	passwordHash string,
	now time.Time,
) error {
	passwordHash = strings.TrimSpace(
		passwordHash,
	)

	if passwordHash == "" {
		return ErrInvalidPasswordHash
	}

	if u.passwordHash == passwordHash {
		return nil
	}

	u.passwordHash = passwordHash
	u.touch(now)

	return nil
}

func (u *User) Activate(
	now time.Time,
) {
	if u.status == valueobject.UserStatusActive {
		return
	}

	u.status = valueobject.UserStatusActive
	u.touch(now)
}

func (u *User) Deactivate(
	now time.Time,
) {
	if u.status == valueobject.UserStatusInactive {
		return
	}

	u.status = valueobject.UserStatusInactive
	u.touch(now)
}

func (u *User) Suspend(
	now time.Time,
) {
	if u.status == valueobject.UserStatusSuspended {
		return
	}

	u.status = valueobject.UserStatusSuspended
	u.touch(now)
}

func (u *User) RecordLogin(
	at time.Time,
) error {
	if err := u.EnsureCanAuthenticate(); err != nil {
		return err
	}

	at = at.UTC()

	u.lastLoginAt = &at
	u.touch(at)

	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) Status() valueobject.UserStatus {
	return u.status
}

func (u *User) LastLoginAt() *time.Time {
	return cloneTimePointer(
		u.lastLoginAt,
	)
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) touch(
	now time.Time,
) {
	now = now.UTC()

	if now.Before(u.updatedAt) {
		return
	}

	u.updatedAt = now
}

func cloneTimePointer(
	value *time.Time,
) *time.Time {
	if value == nil {
		return nil
	}

	return new(value.UTC())
}

func IsUserAccessDenied(
	err error,
) bool {
	return errors.Is(
		err,
		shareddomain.ErrForbidden,
	)
}
