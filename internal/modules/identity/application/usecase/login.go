// Package usecase Files: internal/modules/identity/application/usecase/login.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/application/ports"
	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type LoginUseCase struct {
	users repository.UserRepository

	sessions repository.SessionRepository

	passwordHasher ports.PasswordHasher

	accessTokens ports.AccessTokenService

	refreshTokens ports.RefreshTokenService

	refreshTokenTTL time.Duration
}

func NewLoginUseCase(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	passwordHasher ports.PasswordHasher,
	accessTokens ports.AccessTokenService,
	refreshTokens ports.RefreshTokenService,
	refreshTokenTTL time.Duration,
) *LoginUseCase {
	if users == nil {
		panic("login use case: user repository nil")
	}

	if sessions == nil {
		panic("login use case: session repository nil")
	}

	if passwordHasher == nil {
		panic("login use case: password hasher nil")
	}

	if accessTokens == nil {
		panic("login use case: access token service nil")
	}

	if refreshTokens == nil {
		panic("login use case: refresh token service nil")
	}

	if refreshTokenTTL <= 0 {
		panic("login use case: refresh token TTL tidak valid")
	}

	return &LoginUseCase{
		users: users,

		sessions: sessions,

		passwordHasher: passwordHasher,

		accessTokens: accessTokens,

		refreshTokens: refreshTokens,

		refreshTokenTTL: refreshTokenTTL,
	}
}

func (uc *LoginUseCase) Execute(
	ctx context.Context,
	request dto.LoginRequest,
) (dto.LoginResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.LoginResponse{}, err
	}

	email, err := valueobject.NewEmail(
		request.Email,
	)
	if err != nil {
		return dto.LoginResponse{},
			identitydomain.ErrInvalidCredentials
	}

	user, err := uc.users.FindByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(
			err,
			identitydomain.ErrUserNotFound,
		) {
			return dto.LoginResponse{},
				identitydomain.ErrInvalidCredentials
		}

		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal mencari user login: %w",
				err,
			)
	}

	if err := user.EnsureCanAuthenticate(); err != nil {
		return dto.LoginResponse{}, err
	}

	valid, err := uc.passwordHasher.Verify(
		request.Password,
		user.PasswordHash(),
	)
	if err != nil {
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal memverifikasi password: %w",
				err,
			)
	}

	if !valid {
		return dto.LoginResponse{},
			identitydomain.ErrInvalidCredentials
	}

	fingerprint, err := valueobject.NewSessionFingerprint(
		request.UserAgent,
		request.IPAddress,
	)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	now := time.Now().UTC()

	generatedRefreshToken, err :=
		uc.refreshTokens.Generate()
	if err != nil {
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal membuat refresh token: %w",
				err,
			)
	}

	refreshExpiresAt := now.Add(
		uc.refreshTokenTTL,
	)

	session, err := entity.NewRefreshSession(
		user.ID(),
		uuid.New(),
		generatedRefreshToken.Hash,
		fingerprint,
		refreshExpiresAt,
		now,
	)
	if err != nil {
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal membuat refresh session: %w",
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
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal membuat access token: %w",
				err,
			)
	}

	if err := uc.sessions.Create(
		ctx,
		session,
	); err != nil {
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal menyimpan refresh session: %w",
				err,
			)
	}

	if err := user.RecordLogin(now); err != nil {
		return dto.LoginResponse{}, err
	}

	if err := uc.users.Update(
		ctx,
		user,
	); err != nil {
		return dto.LoginResponse{},
			fmt.Errorf(
				"gagal memperbarui login user: %w",
				err,
			)
	}

	return dto.LoginResponse{
		User: mapAuthUser(
			user,
		),

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

func mapAuthUser(
	user *entity.User,
) dto.AuthUserResponse {
	return dto.AuthUserResponse{
		ID: user.ID().String(),

		Email: user.Email().String(),

		Status: user.Status().String(),
	}
}
