// Package usecase Files: internal/modules/identity/application/usecase/revoke_session.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

const sessionRevocationReasonManual = "MANUAL_REVOKE"

type RevokeSessionUseCase struct {
	sessions repository.SessionRepository
	now      func() time.Time
}

func NewRevokeSessionUseCase(
	sessions repository.SessionRepository,
) *RevokeSessionUseCase {
	if sessions == nil {
		panic("revoke session use case: session repository nil")
	}

	return &RevokeSessionUseCase{
		sessions: sessions,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (uc *RevokeSessionUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	sessionID uuid.UUID,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if userID == uuid.Nil {
		return identitydomain.ErrUserNotFound
	}

	if sessionID == uuid.Nil {
		return identitydomain.ErrSessionNotFound
	}

	session, err := uc.sessions.FindByID(
		ctx,
		sessionID,
	)
	if err != nil {
		return err
	}

	if session.UserID() != userID {
		return identitydomain.ErrSessionNotFound
	}

	now := uc.now()

	session.Revoke(
		sessionRevocationReasonManual,
		now,
	)

	if err := uc.sessions.Update(
		ctx,
		session,
	); err != nil {
		return fmt.Errorf(
			"gagal mencabut session: %w",
			err,
		)
	}

	return nil
}
