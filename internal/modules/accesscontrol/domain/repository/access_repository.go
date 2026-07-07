package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/accesscontrol/domain/entity"
)

type AccessRepository interface {
	FindRoleByID(
		ctx context.Context,
		roleID uuid.UUID,
	) (*entity.Role, error)

	FindRoleByCode(
		ctx context.Context,
		code string,
	) (*entity.Role, error)

	FindPermissionByID(
		ctx context.Context,
		permissionID uuid.UUID,
	) (*entity.Permission, error)

	FindPermissionByCode(
		ctx context.Context,
		code string,
	) (*entity.Permission, error)

	FindRolesByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*entity.Role, error)

	FindPermissionsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*entity.Permission, error)
}
