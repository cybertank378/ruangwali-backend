package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type Student struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	FullName   string
	Gender     Gender
	Religion   *string
	BirthPlace *string
	BirthDate  *time.Time
	Address    *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewStudent(tenantID uuid.UUID, fullName string, gender Gender) (*Student, error) {
	fullName = strings.TrimSpace(fullName)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant id wajib tersedia")
	}
	if fullName == "" {
		return nil, errors.New("nama siswa wajib diisi")
	}
	if gender != GenderMale && gender != GenderFemale {
		return nil, errors.New("jenis kelamin tidak valid")
	}

	now := time.Now().UTC()
	return &Student{
		ID: uuid.New(), TenantID: tenantID, FullName: fullName,
		Gender: gender, CreatedAt: now, UpdatedAt: now,
	}, nil
}
