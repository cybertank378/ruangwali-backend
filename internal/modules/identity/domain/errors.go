package domain

import (
	"fmt"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

var (
	ErrUserNotFound = fmt.Errorf(
		"user tidak ditemukan: %w",
		shareddomain.ErrNotFound,
	)

	ErrEmailAlreadyExists = fmt.Errorf(
		"email sudah digunakan: %w",
		shareddomain.ErrConflict,
	)

	ErrInvalidCredentials = fmt.Errorf(
		"kredensial tidak valid: %w",
		shareddomain.ErrUnauthorized,
	)

	ErrSessionNotFound = fmt.Errorf(
		"session tidak ditemukan: %w",
		shareddomain.ErrNotFound,
	)

	ErrSessionTokenNotFound = fmt.Errorf(
		"session token tidak ditemukan: %w",
		shareddomain.ErrNotFound,
	)
)
