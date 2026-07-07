package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Role struct {
	id          uuid.UUID
	code        string
	name        string
	description *string
	isSystem    bool
	isActive    bool
}

func NewRole(
	id uuid.UUID,
	code string,
	name string,
	description *string,
	isSystem bool,
	isActive bool,
) (*Role, error) {
	code = strings.ToUpper(
		strings.TrimSpace(code),
	)

	name = strings.TrimSpace(name)

	if id == uuid.Nil {
		return nil, errors.New(
			"role id tidak valid",
		)
	}

	if code == "" {
		return nil, errors.New(
			"role code wajib tersedia",
		)
	}

	if name == "" {
		return nil, errors.New(
			"role name wajib tersedia",
		)
	}

	description = normalizeOptionalString(
		description,
	)

	return &Role{
		id:          id,
		code:        code,
		name:        name,
		description: description,
		isSystem:    isSystem,
		isActive:    isActive,
	}, nil
}

func (r *Role) ID() uuid.UUID {
	return r.id
}

func (r *Role) Code() string {
	return r.code
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Description() *string {
	return cloneStringPointer(
		r.description,
	)
}

func (r *Role) IsSystem() bool {
	return r.isSystem
}

func (r *Role) IsActive() bool {
	return r.isActive
}

func (r *Role) Activate() {
	r.isActive = true
}

func (r *Role) Deactivate() {
	r.isActive = false
}

func normalizeOptionalString(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(
		*value,
	)

	if normalized == "" {
		return nil
	}

	return &normalized
}

func cloneStringPointer(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	return new(*value)
}
