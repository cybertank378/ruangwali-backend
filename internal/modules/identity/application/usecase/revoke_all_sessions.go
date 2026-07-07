// Package usecase Files: internal/modules/identity/application/usecase/revoke_all_sessions.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

const (
	sessionRevocationReasonLogoutAll = "LOGOUT_ALL"
)

type RevokeAllSessionsUseCase struct {
	sessions repository.SessionRepository
	now      func() time.Time
}

func NewRevokeAllSessionsUseCase(
	sessions repository.SessionRepository,
) *RevokeAllSessionsUseCase {
	if sessions == nil {
		panic("revoke all sessions use case: session repository nil")
	}

	return &RevokeAllSessionsUseCase{
		sessions: sessions,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (uc *RevokeAllSessionsUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	exceptSessionID *uuid.UUID,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if userID == uuid.Nil {
		return identitydomain.ErrUserNotFound
	}

	now := uc.now()

	if exceptSessionID != nil &&
		*exceptSessionID != uuid.Nil {

		if err := uc.sessions.RevokeByUserIDExcept(
			ctx,
			userID,
			*exceptSessionID,
			sessionRevocationReasonLogoutAll,
			now,
		); err != nil {
			return fmt.Errorf(
				"gagal mencabut seluruh session selain session aktif: %w",
				err,
			)
		}

		return nil
	}

	if err := uc.sessions.RevokeByUserID(
		ctx,
		userID,
		sessionRevocationReasonLogoutAll,
		now,
	); err != nil {
		return fmt.Errorf(
			"gagal mencabut seluruh session: %w",
			err,
		)
	}

	return nil
}
