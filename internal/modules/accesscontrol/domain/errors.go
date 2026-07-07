package domain

import (
	"fmt"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

var (
	ErrRoleNotFound = fmt.Errorf(
		"role tidak ditemukan: %w",
		shareddomain.ErrNotFound,
	)

	ErrPermissionNotFound = fmt.Errorf(
		"permission tidak ditemukan: %w",
		shareddomain.ErrNotFound,
	)

	ErrAccessDenied = fmt.Errorf(
		"akses ditolak: %w",
		shareddomain.ErrForbidden,
	)

	ErrInactiveRole = fmt.Errorf(
		"role tidak aktif: %w",
		shareddomain.ErrForbidden,
	)
)
