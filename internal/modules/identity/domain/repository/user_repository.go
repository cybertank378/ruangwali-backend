package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type UserRepository interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.User, error)

	FindByEmail(
		ctx context.Context,
		email valueobject.Email,
	) (*entity.User, error)

	ExistsByEmail(
		ctx context.Context,
		email valueobject.Email,
	) (bool, error)

	Create(
		ctx context.Context,
		user *entity.User,
	) error

	Update(
		ctx context.Context,
		user *entity.User,
	) error
}
