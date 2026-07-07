package domain

import "errors"

var (
	ErrValidation = errors.New(
		"validation error",
	)

	ErrNotFound = errors.New(
		"resource not found",
	)

	ErrConflict = errors.New(
		"resource conflict",
	)

	ErrUnauthorized = errors.New(
		"authentication required",
	)

	ErrForbidden = errors.New(
		"access forbidden",
	)
)
