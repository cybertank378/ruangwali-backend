// Package usecase Files: internal/modules/identity/application/usecase/logout.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/application/ports"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

type LogoutUseCase struct {
	sessions repository.SessionRepository

	refreshTokens ports.RefreshTokenService
}

func NewLogoutUseCase(
	sessions repository.SessionRepository,
	refreshTokens ports.RefreshTokenService,
) *LogoutUseCase {
	if sessions == nil {
		panic("logout use case: session repository nil")
	}

	if refreshTokens == nil {
		panic("logout use case: refresh token service nil")
	}

	return &LogoutUseCase{
		sessions: sessions,

		refreshTokens: refreshTokens,
	}
}

func (uc *LogoutUseCase) Execute(
	ctx context.Context,
	request dto.LogoutRequest,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := uc.refreshTokens.Validate(
		request.RefreshToken,
	); err != nil {
		return err
	}

	tokenHash := uc.refreshTokens.Hash(
		request.RefreshToken,
	)

	session, err := uc.sessions.FindByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	session.Revoke(
		"LOGOUT",
		now,
	)

	if err := uc.sessions.Update(
		ctx,
		session,
	); err != nil {
		return fmt.Errorf(
			"gagal menyimpan logout session: %w",
			err,
		)
	}

	return nil
}
