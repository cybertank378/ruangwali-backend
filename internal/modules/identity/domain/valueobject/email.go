package valueobject

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

const maxEmailLength = 255

var ErrInvalidEmail = fmt.Errorf(
	"email tidak valid: %w",
	shareddomain.ErrValidation,
)

type Email struct {
	value string
}

func NewEmail(
	raw string,
) (Email, error) {
	value := strings.ToLower(
		strings.TrimSpace(raw),
	)

	if value == "" {
		return Email{}, ErrInvalidEmail
	}

	if len(value) > maxEmailLength {
		return Email{}, fmt.Errorf(
			"panjang email melebihi %d karakter: %w",
			maxEmailLength,
			ErrInvalidEmail,
		)
	}

	address, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, fmt.Errorf(
			"format email tidak valid: %w",
			ErrInvalidEmail,
		)
	}

	if address.Address != value {
		return Email{}, fmt.Errorf(
			"format email tidak valid: %w",
			ErrInvalidEmail,
		)
	}

	return Email{
		value: value,
	}, nil
}

func MustEmail(
	raw string,
) Email {
	email, err := NewEmail(raw)
	if err != nil {
		panic(err)
	}

	return email
}

func (e Email) String() string {
	return e.value
}

func (e Email) IsZero() bool {
	return e.value == ""
}

func (e Email) Equal(
	other Email,
) bool {
	return e.value == other.value
}

func IsInvalidEmail(
	err error,
) bool {
	return errors.Is(
		err,
		ErrInvalidEmail,
	)
}
