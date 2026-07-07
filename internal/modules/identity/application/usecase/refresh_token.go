package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/application/ports"
	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type RefreshTokenUseCase struct {
	users repository.UserRepository

	sessions repository.SessionRepository

	accessTokens ports.AccessTokenService

	refreshTokens ports.RefreshTokenService

	refreshTokenTTL time.Duration
}

func NewRefreshTokenUseCase(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	accessTokens ports.AccessTokenService,
	refreshTokens ports.RefreshTokenService,
	refreshTokenTTL time.Duration,
) *RefreshTokenUseCase {
	if users == nil {
		panic("refresh token use case: user repository nil")
	}

	if sessions == nil {
		panic("refresh token use case: session repository nil")
	}

	if accessTokens == nil {
		panic("refresh token use case: access token service nil")
	}

	if refreshTokens == nil {
		panic("refresh token use case: refresh token service nil")
	}

	if refreshTokenTTL <= 0 {
		panic(
			"refresh token use case: refresh token TTL tidak valid",
		)
	}

	return &RefreshTokenUseCase{
		users: users,

		sessions: sessions,

		accessTokens: accessTokens,

		refreshTokens: refreshTokens,

		refreshTokenTTL: refreshTokenTTL,
	}
}

func (uc *RefreshTokenUseCase) Execute(
	ctx context.Context,
	request dto.RefreshTokenRequest,
) (dto.RefreshTokenResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	if err := uc.refreshTokens.Validate(
		request.RefreshToken,
	); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	currentTokenHash := uc.refreshTokens.Hash(
		request.RefreshToken,
	)

	currentSession, err :=
		uc.sessions.FindByTokenHash(
			ctx,
			currentTokenHash,
		)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	now := time.Now().UTC()

	if err := currentSession.EnsureUsable(
		now,
	); err != nil {
		if currentSession.IsRevoked() {
			_ = uc.sessions.RevokeByFamilyID(
				ctx,
				currentSession.FamilyID(),
				"REFRESH_TOKEN_REUSE_DETECTED",
				now,
			)
		}

		return dto.RefreshTokenResponse{}, err
	}

	user, err := uc.users.FindByID(
		ctx,
		currentSession.UserID(),
	)
	if err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal mengambil pemilik session: %w",
				err,
			)
	}

	if err := user.EnsureCanAuthenticate(); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	fingerprint, err := valueobject.NewSessionFingerprint(
		request.UserAgent,
		request.IPAddress,
	)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	generatedRefreshToken, err :=
		uc.refreshTokens.Generate()
	if err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal membuat refresh token baru: %w",
				err,
			)
	}

	refreshExpiresAt := now.Add(
		uc.refreshTokenTTL,
	)

	replacementSession, err :=
		entity.NewRefreshSession(
			user.ID(),
			currentSession.FamilyID(),
			generatedRefreshToken.Hash,
			fingerprint,
			refreshExpiresAt,
			now,
		)
	if err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal membuat replacement session: %w",
				err,
			)
	}

	accessToken, err := uc.accessTokens.Issue(
		ctx,
		ports.AccessTokenClaims{
			UserID: user.ID(),
		},
	)
	if err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal membuat access token baru: %w",
				err,
			)
	}

	if err := uc.sessions.Create(
		ctx,
		replacementSession,
	); err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal menyimpan replacement session: %w",
				err,
			)
	}

	if err := currentSession.ReplaceWith(
		replacementSession.ID(),
		now,
	); err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	if err := uc.sessions.Update(
		ctx,
		currentSession,
	); err != nil {
		return dto.RefreshTokenResponse{},
			fmt.Errorf(
				"gagal menyimpan rotasi session: %w",
				err,
			)
	}

	return dto.RefreshTokenResponse{
		Tokens: dto.TokenResponse{
			AccessToken: accessToken,

			RefreshToken: generatedRefreshToken.Raw,

			TokenType: "Bearer",

			AccessTokenExpiresAt: now.Add(
				uc.accessTokens.TTL(),
			),

			RefreshTokenExpiresAt: refreshExpiresAt,
		},
	}, nil
}
