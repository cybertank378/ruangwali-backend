// =========================================================
// File: internal/composition/identity.go
// =========================================================

package composition

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruangwali/internal/modules/identity/application/usecase"
	identitynotification "github.com/ruangwali/internal/modules/identity/infrastructure/notification"
	identitypostgres "github.com/ruangwali/internal/modules/identity/infrastructure/postgres"
	identitysecurity "github.com/ruangwali/internal/modules/identity/infrastructure/security"
	identityhttp "github.com/ruangwali/internal/modules/identity/presentation/http"
	"github.com/ruangwali/internal/platform/config"
)

const (
	defaultPasswordResetNotifierTimeout = 15 * time.Second
)

type IdentityModule struct {
	Handler *identityhttp.Handler

	AuthMiddleware *identityhttp.AuthMiddleware
}

func buildIdentityModule(
	cfg config.Config,
	db *pgxpool.Pool,
) (
	*IdentityModule,
	error,
) {
	if db == nil {
		return nil, fmt.Errorf(
			"identity module: database pool nil",
		)
	}

	// =====================================================
	// REPOSITORIES
	// =====================================================

	userRepository :=
		identitypostgres.NewUserRepository(
			db,
		)

	sessionRepository :=
		identitypostgres.NewSessionRepository(
			db,
		)

	passwordResetRepository :=
		identitypostgres.NewPasswordResetRepository(
			db,
		)

	// =====================================================
	// SECURITY SERVICES
	// =====================================================

	passwordHasher :=
		identitysecurity.NewPasswordHasher()

	accessTokenService, err :=
		identitysecurity.NewTokenService(
			cfg.Auth.JWTIssuer,
			cfg.Auth.JWTAudience,
			cfg.Auth.JWTSecret,
			cfg.Auth.AccessTokenTTL,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membangun identity access token service: %w",
			err,
		)
	}

	refreshTokenService :=
		identitysecurity.NewRefreshTokenService()

	passwordResetTokenService, err :=
		identitysecurity.NewPasswordResetTokenService(
			cfg.Auth.PasswordResetTokenTTL,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membangun password reset token service: %w",
			err,
		)
	}

	// =====================================================
	// HTTP CLIENTS
	// =====================================================

	passwordResetHTTPClient :=
		&http.Client{
			Timeout: defaultPasswordResetNotifierTimeout,
		}

	// =====================================================
	// NOTIFICATION SERVICES
	// =====================================================

	passwordResetNotifier, err :=
		identitynotification.NewPasswordResetNotifier(
			passwordResetHTTPClient,
			cfg.Integration.GoogleAppsScript.BaseURL,
			cfg.Integration.GoogleAppsScript.WebhookSecret,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membangun password reset notifier: %w",
			err,
		)
	}

	// =====================================================
	// USE CASES
	// =====================================================

	loginUseCase :=
		usecase.NewLoginUseCase(
			userRepository,
			sessionRepository,
			passwordHasher,
			accessTokenService,
			refreshTokenService,
			cfg.Auth.RefreshTokenTTL,
		)

	logoutUseCase :=
		usecase.NewLogoutUseCase(
			sessionRepository,
			refreshTokenService,
		)

	refreshTokenUseCase :=
		usecase.NewRefreshTokenUseCase(
			userRepository,
			sessionRepository,
			accessTokenService,
			refreshTokenService,
			cfg.Auth.RefreshTokenTTL,
		)

	getCurrentUserUseCase :=
		usecase.NewGetCurrentUserUseCase(
			userRepository,
		)

	changePasswordUseCase :=
		usecase.NewChangePasswordUseCase(
			userRepository,
			sessionRepository,
			passwordHasher,
		)

	forgotPasswordUseCase :=
		usecase.NewForgotPasswordUseCase(
			userRepository,
			passwordResetRepository,
			passwordResetTokenService,
			passwordResetNotifier,
		)

	resetPasswordUseCase :=
		usecase.NewResetPasswordUseCase(
			userRepository,
			sessionRepository,
			passwordResetRepository,
			passwordResetTokenService,
			passwordHasher,
		)

	listSessionsUseCase :=
		usecase.NewListSessionsUseCase(
			sessionRepository,
		)

	revokeSessionUseCase :=
		usecase.NewRevokeSessionUseCase(
			sessionRepository,
		)

	revokeAllSessionsUseCase :=
		usecase.NewRevokeAllSessionsUseCase(
			sessionRepository,
		)

	// =====================================================
	// PRESENTATION
	// =====================================================

	handler :=
		identityhttp.NewHandler(
			loginUseCase,
			logoutUseCase,
			refreshTokenUseCase,
			getCurrentUserUseCase,
			changePasswordUseCase,
			forgotPasswordUseCase,
			resetPasswordUseCase,
			listSessionsUseCase,
			revokeSessionUseCase,
			revokeAllSessionsUseCase,
		)

	authMiddleware :=
		identityhttp.NewAuthMiddleware(
			accessTokenService,
		)

	return &IdentityModule{
		Handler: handler,

		AuthMiddleware: authMiddleware,
	}, nil
}
