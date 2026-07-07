package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type PasswordResetRepository interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.PasswordReset, error)

	FindByTokenHash(
		ctx context.Context,
		tokenHash []byte,
	) (*entity.PasswordReset, error)

	FindActiveByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*entity.PasswordReset, error)

	Create(
		ctx context.Context,
		passwordReset *entity.PasswordReset,
	) error

	Update(
		ctx context.Context,
		passwordReset *entity.PasswordReset,
	) error
}
