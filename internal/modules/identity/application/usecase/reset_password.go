// Package usecase Files: internal/modules/identity/application/usecase/reset_password.go
package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/application/ports"
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

const sessionRevocationReasonPasswordReset = "PASSWORD_RESET"

type ResetPasswordUseCase struct {
	users          repository.UserRepository
	sessions       repository.SessionRepository
	passwordResets repository.PasswordResetRepository
	resetTokens    ports.PasswordResetTokenService
	passwords      ports.PasswordHasher
	now            func() time.Time
}

func NewResetPasswordUseCase(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	passwordResets repository.PasswordResetRepository,
	resetTokens ports.PasswordResetTokenService,
	passwords ports.PasswordHasher,
) *ResetPasswordUseCase {
	if users == nil {
		panic("reset password use case: user repository nil")
	}

	if sessions == nil {
		panic("reset password use case: session repository nil")
	}

	if passwordResets == nil {
		panic("reset password use case: password reset repository nil")
	}

	if resetTokens == nil {
		panic("reset password use case: password reset token service nil")
	}

	if passwords == nil {
		panic("reset password use case: password hasher nil")
	}

	return &ResetPasswordUseCase{
		users:          users,
		sessions:       sessions,
		passwordResets: passwordResets,
		resetTokens:    resetTokens,
		passwords:      passwords,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (uc *ResetPasswordUseCase) Execute(
	ctx context.Context,
	request dto.ResetPasswordRequest,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	rawToken := strings.TrimSpace(
		request.Token,
	)
	if rawToken == "" {
		return fmt.Errorf(
			"password reset token wajib diisi",
		)
	}

	tokenHash, err := uc.resetTokens.Hash(
		ctx,
		rawToken,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat hash password reset token: %w",
			err,
		)
	}

	passwordReset, err := uc.passwordResets.FindByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return err
	}

	now := uc.now()

	if err := passwordReset.EnsureUsable(
		now,
	); err != nil {
		return err
	}

	user, err := uc.users.FindByID(
		ctx,
		passwordReset.UserID(),
	)
	if err != nil {
		return err
	}

	samePassword, err := uc.passwords.Verify(
		request.NewPassword.String(),
		user.PasswordHash(),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membandingkan password baru: %w",
			err,
		)
	}

	if samePassword {
		return ErrNewPasswordSameAsCurrent
	}

	passwordHash, err := uc.passwords.Hash(
		request.NewPassword.String(),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat password hash baru: %w",
			err,
		)
	}

	if err := user.ChangePasswordHash(
		passwordHash,
		now,
	); err != nil {
		return err
	}

	if err := passwordReset.MarkUsed(
		now,
	); err != nil {
		return err
	}

	if err := uc.users.Update(
		ctx,
		user,
	); err != nil {
		return fmt.Errorf(
			"gagal memperbarui password user: %w",
			err,
		)
	}

	if err := uc.passwordResets.Update(
		ctx,
		passwordReset,
	); err != nil {
		return fmt.Errorf(
			"gagal menandai password reset telah digunakan: %w",
			err,
		)
	}

	if err := uc.sessions.RevokeByUserID(
		ctx,
		user.ID(),
		sessionRevocationReasonPasswordReset,
		now,
	); err != nil {
		return fmt.Errorf(
			"gagal mencabut session user: %w",
			err,
		)
	}

	return nil
}
