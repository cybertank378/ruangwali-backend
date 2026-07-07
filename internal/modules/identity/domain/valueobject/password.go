package valueobject

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

var (
	ErrPasswordRequired = fmt.Errorf(
		"password wajib diisi: %w",
		shareddomain.ErrValidation,
	)

	ErrPasswordTooShort = fmt.Errorf(
		"password minimal %d karakter: %w",
		MinPasswordLength,
		shareddomain.ErrValidation,
	)

	ErrPasswordTooLong = fmt.Errorf(
		"password maksimal %d karakter: %w",
		MaxPasswordLength,
		shareddomain.ErrValidation,
	)

	ErrPasswordMissingUppercase = fmt.Errorf(
		"password wajib memiliki minimal satu huruf besar: %w",
		shareddomain.ErrValidation,
	)

	ErrPasswordMissingLowercase = fmt.Errorf(
		"password wajib memiliki minimal satu huruf kecil: %w",
		shareddomain.ErrValidation,
	)

	ErrPasswordMissingDigit = fmt.Errorf(
		"password wajib memiliki minimal satu angka: %w",
		shareddomain.ErrValidation,
	)
)

type Password struct {
	value string
}

func NewPassword(
	raw string,
) (Password, error) {
	if raw == "" {
		return Password{}, ErrPasswordRequired
	}

	length := utf8.RuneCountInString(raw)

	if length < MinPasswordLength {
		return Password{}, ErrPasswordTooShort
	}

	if length > MaxPasswordLength {
		return Password{}, ErrPasswordTooLong
	}

	var (
		hasUppercase bool
		hasLowercase bool
		hasDigit     bool
	)

	for _, character := range raw {
		switch {
		case unicode.IsUpper(character):
			hasUppercase = true

		case unicode.IsLower(character):
			hasLowercase = true

		case unicode.IsDigit(character):
			hasDigit = true
		}
	}

	if !hasUppercase {
		return Password{}, ErrPasswordMissingUppercase
	}

	if !hasLowercase {
		return Password{}, ErrPasswordMissingLowercase
	}

	if !hasDigit {
		return Password{}, ErrPasswordMissingDigit
	}

	return Password{
		value: raw,
	}, nil
}

func (p Password) String() string {
	return p.value
}

func (p Password) IsZero() bool {
	return p.value == ""
}
