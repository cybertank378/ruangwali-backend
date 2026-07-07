// Package usecase Files: internal/modules/identity/application/usecase/change_password.go
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
	"github.com/ruangwali/internal/modules/identity/domain/repository"
)

var (
	ErrCurrentPasswordInvalid = fmt.Errorf(
		"password saat ini tidak valid: %w",
		identitydomain.ErrInvalidCredentials,
	)

	ErrNewPasswordSameAsCurrent = errors.New(
		"password baru tidak boleh sama dengan password saat ini",
	)
)

const sessionRevocationReasonPasswordChanged = "password_changed"

type ChangePasswordUseCase struct {
	users     repository.UserRepository
	sessions  repository.SessionRepository
	passwords ports.PasswordHasher
	now       func() time.Time
}

func NewChangePasswordUseCase(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	passwords ports.PasswordHasher,
) *ChangePasswordUseCase {
	if users == nil {
		panic("change password use case: user repository nil")
	}

	if sessions == nil {
		panic("change password use case: session repository nil")
	}

	if passwords == nil {
		panic("change password use case: password hasher nil")
	}

	return &ChangePasswordUseCase{
		users:     users,
		sessions:  sessions,
		passwords: passwords,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (uc *ChangePasswordUseCase) Execute(
	ctx context.Context,
	request dto.ChangePasswordRequest,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if request.UserID == uuid.Nil {
		return identitydomain.ErrUserNotFound
	}

	user, err := uc.users.FindByID(
		ctx,
		request.UserID,
	)
	if err != nil {
		return err
	}

	valid, err := uc.passwords.Verify(
		request.CurrentPassword,
		user.PasswordHash(),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal memverifikasi password saat ini: %w",
			err,
		)
	}

	if !valid {
		return ErrCurrentPasswordInvalid
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

	now := uc.now()

	if err := user.ChangePasswordHash(
		passwordHash,
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

	if err := uc.sessions.RevokeByUserID(
		ctx,
		request.UserID,
		sessionRevocationReasonPasswordChanged,
		now,
	); err != nil {
		return fmt.Errorf(
			"password berhasil diubah tetapi gagal mencabut session: %w",
			err,
		)
	}

	return nil
}
