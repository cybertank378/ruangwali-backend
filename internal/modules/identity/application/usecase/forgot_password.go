// Package usecase Files: internal/modules/identity/application/usecase/forgot_password.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/application/ports"
	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/entity"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

type ForgotPasswordUseCase struct {
	users          repository.UserRepository
	passwordResets repository.PasswordResetRepository
	resetTokens    ports.PasswordResetTokenService
	notifier       ports.PasswordResetNotifier
	now            func() time.Time
}

func NewForgotPasswordUseCase(
	users repository.UserRepository,
	passwordResets repository.PasswordResetRepository,
	resetTokens ports.PasswordResetTokenService,
	notifier ports.PasswordResetNotifier,
) *ForgotPasswordUseCase {
	if users == nil {
		panic("forgot password use case: user repository nil")
	}

	if passwordResets == nil {
		panic("forgot password use case: password reset repository nil")
	}

	if resetTokens == nil {
		panic("forgot password use case: password reset token service nil")
	}

	if notifier == nil {
		panic("forgot password use case: password reset notifier nil")
	}

	return &ForgotPasswordUseCase{
		users:          users,
		passwordResets: passwordResets,
		resetTokens:    resetTokens,
		notifier:       notifier,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (uc *ForgotPasswordUseCase) Execute(
	ctx context.Context,
	request dto.ForgotPasswordRequest,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	user, err := uc.users.FindByEmail(
		ctx,
		request.Email,
	)
	if err != nil {
		if errors.Is(
			err,
			identitydomain.ErrUserNotFound,
		) {
			return nil
		}

		return fmt.Errorf(
			"gagal mencari user untuk password reset: %w",
			err,
		)
	}

	generatedToken, err := uc.resetTokens.Generate(ctx)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat password reset token: %w",
			err,
		)
	}

	now := uc.now()
	expiresAt := now.Add(
		uc.resetTokens.TTL(),
	)

	passwordReset, err := entity.NewPasswordReset(
		user.ID(),
		generatedToken.Hash,
		expiresAt,
		now,
	)
	if err != nil {
		return err
	}

	activeReset, err := uc.passwordResets.FindActiveByUserID(
		ctx,
		user.ID(),
	)
	if err != nil &&
		!errors.Is(
			err,
			identitydomain.ErrPasswordResetNotFound,
		) {
		return fmt.Errorf(
			"gagal mencari password reset aktif: %w",
			err,
		)
	}

	if activeReset != nil {
		if err := activeReset.ReplaceWith(
			passwordReset.ID(),
			now,
		); err != nil {
			return err
		}

		if err := uc.passwordResets.Update(
			ctx,
			activeReset,
		); err != nil {
			return fmt.Errorf(
				"gagal mencabut password reset sebelumnya: %w",
				err,
			)
		}
	}

	if err := uc.passwordResets.Create(
		ctx,
		passwordReset,
	); err != nil {
		return fmt.Errorf(
			"gagal menyimpan password reset: %w",
			err,
		)
	}

	if err := uc.notifier.SendPasswordReset(
		ctx,
		ports.PasswordResetNotification{
			UserID:    user.ID(),
			Email:     request.Email,
			Token:     generatedToken.Raw,
			ExpiresAt: expiresAt,
		},
	); err != nil {
		return fmt.Errorf(
			"gagal mengirim notifikasi password reset: %w",
			err,
		)
	}

	return nil
}
