package dto

import (
	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type ChangePasswordRequest struct {
	UserID uuid.UUID

	CurrentPassword string

	NewPassword valueobject.Password
}

type ForgotPasswordRequest struct {
	Email valueobject.Email
}

type ResetPasswordRequest struct {
	Token string

	NewPassword valueobject.Password
}
