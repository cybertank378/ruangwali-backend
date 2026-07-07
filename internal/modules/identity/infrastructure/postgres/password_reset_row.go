package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type passwordResetRow struct {
	ID uuid.UUID

	UserID uuid.UUID

	TokenHash []byte

	ExpiresAt time.Time

	UsedAt *time.Time

	RevokedAt *time.Time

	ReplacedByID *uuid.UUID

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (r passwordResetRow) toEntity() (
	*entity.PasswordReset,
	error,
) {
	return entity.RehydratePasswordReset(
		entity.RehydratePasswordResetParams{
			ID: r.ID,

			UserID: r.UserID,

			TokenHash: r.TokenHash,

			ExpiresAt: r.ExpiresAt,

			UsedAt: r.UsedAt,

			RevokedAt: r.RevokedAt,

			ReplacedByID: r.ReplacedByID,

			CreatedAt: r.CreatedAt,

			UpdatedAt: r.UpdatedAt,
		},
	)
}

func passwordResetRowFromEntity(
	passwordReset *entity.PasswordReset,
) passwordResetRow {
	return passwordResetRow{
		ID: passwordReset.ID(),

		UserID: passwordReset.UserID(),

		TokenHash: passwordReset.TokenHash(),

		ExpiresAt: passwordReset.ExpiresAt(),

		UsedAt: passwordReset.UsedAt(),

		RevokedAt: passwordReset.RevokedAt(),

		ReplacedByID: passwordReset.ReplacedByID(),

		CreatedAt: passwordReset.CreatedAt(),

		UpdatedAt: passwordReset.UpdatedAt(),
	}
}
