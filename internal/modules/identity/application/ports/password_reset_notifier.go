// Package ports FilePath: internal/modules/identity/application/ports/password_reset_notifier.go
package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type PasswordResetNotification struct {
	UserID uuid.UUID

	Email valueobject.Email

	Token string

	ExpiresAt time.Time
}

type PasswordResetNotifier interface {
	SendPasswordReset(
		ctx context.Context,
		notification PasswordResetNotification,
	) error
}
