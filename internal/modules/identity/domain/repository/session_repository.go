package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type SessionRepository interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.RefreshSession, error)

	FindByTokenHash(
		ctx context.Context,
		tokenHash []byte,
	) (*entity.RefreshSession, error)

	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*entity.RefreshSession, error)

	Create(
		ctx context.Context,
		session *entity.RefreshSession,
	) error

	Update(
		ctx context.Context,
		session *entity.RefreshSession,
	) error

	RevokeByUserID(
		ctx context.Context,
		userID uuid.UUID,
		reason string,
		revokedAt time.Time,
	) error

	RevokeByUserIDExcept(
		ctx context.Context,
		userID uuid.UUID,
		exceptSessionID uuid.UUID,
		reason string,
		revokedAt time.Time,
	) error

	RevokeByFamilyID(
		ctx context.Context,
		familyID uuid.UUID,
		reason string,
		revokedAt time.Time,
	) error

	DeleteExpiredBefore(
		ctx context.Context,
		before time.Time,
	) (int64, error)
}
