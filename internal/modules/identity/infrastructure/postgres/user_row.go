package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type userRow struct {
	ID uuid.UUID

	Email string

	PasswordHash string

	Status string

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r userRow) toEntity() (*entity.User, error) {
	email, err := valueobject.NewEmail(
		r.Email,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal rehydrate email user: %w",
			err,
		)
	}

	status, err := valueobject.ParseUserStatus(
		r.Status,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal rehydrate status user: %w",
			err,
		)
	}

	user, err := entity.RehydrateUser(
		entity.RehydrateUserParams{
			ID: r.ID,

			Email: email,

			PasswordHash: r.PasswordHash,

			Status: status,

			LastLoginAt: r.LastLoginAt,

			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal rehydrate user: %w",
			err,
		)
	}

	return user, nil
}
