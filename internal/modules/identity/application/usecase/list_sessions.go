// Package usecase Files: internal/modules/identity/application/usecase/list_sessions.go
package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

type ListSessionsUseCase struct {
	sessions repository.SessionRepository
}

func NewListSessionsUseCase(
	sessions repository.SessionRepository,
) *ListSessionsUseCase {
	if sessions == nil {
		panic("list sessions use case: session repository nil")
	}

	return &ListSessionsUseCase{
		sessions: sessions,
	}
}

func (uc *ListSessionsUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	currentSessionID uuid.UUID,
) (dto.SessionListResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.SessionListResponse{}, err
	}

	if userID == uuid.Nil {
		return dto.SessionListResponse{},
			identitydomain.ErrUserNotFound
	}

	sessions, err := uc.sessions.FindByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return dto.SessionListResponse{}, fmt.Errorf(
			"gagal mengambil daftar session: %w",
			err,
		)
	}

	items := make(
		[]dto.SessionResponse,
		0,
		len(sessions),
	)

	for _, session := range sessions {
		if session == nil {
			continue
		}

		fingerprint := session.Fingerprint()

		items = append(
			items,
			dto.SessionResponse{
				ID: session.ID(),

				UserAgent: fingerprint.UserAgent(),

				IPAddress: fingerprint.IPAddress(),

				ExpiresAt: session.ExpiresAt(),

				LastUsedAt: session.LastUsedAt(),

				CreatedAt: session.CreatedAt(),

				IsCurrent: session.ID() ==
					currentSessionID,
			},
		)
	}

	return dto.SessionListResponse{
		Items: items,
	}, nil
}
