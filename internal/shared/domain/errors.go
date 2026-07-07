package domain

import "errors"

var (
	ErrNotFound  = errors.New("resource tidak ditemukan")
	ErrConflict  = errors.New("resource conflict")
	ErrForbidden = errors.New("akses ditolak")
	ErrInvalid   = errors.New("data tidak valid")
)
