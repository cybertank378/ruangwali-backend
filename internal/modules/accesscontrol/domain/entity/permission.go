package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Permission struct {
	id          uuid.UUID
	code        string
	resource    string
	action      string
	description *string
}

func NewPermission(
	id uuid.UUID,
	code string,
	resource string,
	action string,
	description *string,
) (*Permission, error) {
	code = strings.ToLower(
		strings.TrimSpace(code),
	)

	resource = strings.ToLower(
		strings.TrimSpace(resource),
	)

	action = strings.ToLower(
		strings.TrimSpace(action),
	)

	if id == uuid.Nil {
		return nil, errors.New(
			"permission id tidak valid",
		)
	}

	if code == "" {
		return nil, errors.New(
			"permission code wajib tersedia",
		)
	}

	if resource == "" {
		return nil, errors.New(
			"permission resource wajib tersedia",
		)
	}

	if action == "" {
		return nil, errors.New(
			"permission action wajib tersedia",
		)
	}

	expectedCode :=
		resource + "." + action

	if code != expectedCode {
		return nil, errors.New(
			"permission code tidak sesuai resource dan action",
		)
	}

	description = normalizeOptionalString(
		description,
	)

	return &Permission{
		id:          id,
		code:        code,
		resource:    resource,
		action:      action,
		description: description,
	}, nil
}

func (p *Permission) ID() uuid.UUID {
	return p.id
}

func (p *Permission) Code() string {
	return p.code
}

func (p *Permission) Resource() string {
	return p.resource
}

func (p *Permission) Action() string {
	return p.action
}

func (p *Permission) Description() *string {
	return cloneStringPointer(
		p.description,
	)
}
