package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

type GetCurrentUserUseCase struct {
	users repository.UserRepository
}

func NewGetCurrentUserUseCase(
	users repository.UserRepository,
) *GetCurrentUserUseCase {
	if users == nil {
		panic(
			"get current user use case: user repository nil",
		)
	}

	return &GetCurrentUserUseCase{
		users: users,
	}
}

func (uc *GetCurrentUserUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
) (dto.CurrentUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.CurrentUserResponse{}, err
	}

	if userID == uuid.Nil {
		return dto.CurrentUserResponse{},
			fmt.Errorf(
				"user ID tidak valid",
			)
	}

	user, err := uc.users.FindByID(
		ctx,
		userID,
	)
	if err != nil {
		return dto.CurrentUserResponse{}, err
	}

	if err := user.EnsureCanAuthenticate(); err != nil {
		return dto.CurrentUserResponse{}, err
	}

	return dto.CurrentUserResponse{
		ID: user.ID().String(),

		Email: user.Email().String(),

		Status: user.Status().String(),

		LastLoginAt: user.LastLoginAt(),

		CreatedAt: user.CreatedAt(),

		UpdatedAt: user.UpdatedAt(),
	}, nil
}
