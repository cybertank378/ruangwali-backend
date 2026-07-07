package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type sessionRow struct {
	ID uuid.UUID

	UserID uuid.UUID

	FamilyID uuid.UUID

	TokenHash []byte

	UserAgent *string
	IPAddress *string

	ExpiresAt time.Time

	LastUsedAt *time.Time

	RevokedAt    *time.Time
	RevokeReason *string

	ReplacedByID *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r sessionRow) toEntity() (*entity.RefreshSession, error) {
	userAgent := ""

	if r.UserAgent != nil {
		userAgent = *r.UserAgent
	}

	ipAddress := ""

	if r.IPAddress != nil {
		ipAddress = *r.IPAddress
	}

	fingerprint, err := valueobject.NewSessionFingerprint(
		userAgent,
		ipAddress,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal rehydrate session fingerprint: %w",
			err,
		)
	}

	session, err := entity.RehydrateRefreshSession(
		entity.RehydrateRefreshSessionParams{
			ID: r.ID,

			UserID: r.UserID,

			FamilyID: r.FamilyID,

			TokenHash: r.TokenHash,

			Fingerprint: fingerprint,

			ExpiresAt: r.ExpiresAt,

			LastUsedAt: r.LastUsedAt,

			RevokedAt:    r.RevokedAt,
			RevokeReason: r.RevokeReason,

			ReplacedByID: r.ReplacedByID,

			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal rehydrate refresh session: %w",
			err,
		)
	}

	return session, nil
}
