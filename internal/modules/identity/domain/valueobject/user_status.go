package valueobject

import (
	"fmt"
	"strings"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

type UserStatus string

const (
	UserStatusActive UserStatus = "ACTIVE"

	UserStatusInactive UserStatus = "INACTIVE"

	UserStatusSuspended UserStatus = "SUSPENDED"
)

var ErrInvalidUserStatus = fmt.Errorf(
	"status user tidak valid: %w",
	shareddomain.ErrValidation,
)

func ParseUserStatus(
	raw string,
) (UserStatus, error) {
	status := UserStatus(
		strings.ToUpper(
			strings.TrimSpace(raw),
		),
	)

	if !status.IsValid() {
		return "", fmt.Errorf(
			"%q: %w",
			raw,
			ErrInvalidUserStatus,
		)
	}

	return status, nil
}

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsValid() bool {
	switch s {
	case
		UserStatusActive,
		UserStatusInactive,
		UserStatusSuspended:
		return true

	default:
		return false
	}
}

func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

func (s UserStatus) IsInactive() bool {
	return s == UserStatusInactive
}

func (s UserStatus) IsSuspended() bool {
	return s == UserStatusSuspended
}
