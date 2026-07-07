package requestcontext

import (
	"context"

	"github.com/google/uuid"
)

type contextKey uint8

const (
	authenticatedUserIDKey contextKey = iota
)

func WithUserID(
	ctx context.Context,
	userID uuid.UUID,
) context.Context {
	return context.WithValue(
		ctx,
		authenticatedUserIDKey,
		userID,
	)
}

func UserID(
	ctx context.Context,
) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}

	userID, ok := ctx.Value(
		authenticatedUserIDKey,
	).(uuid.UUID)

	if !ok || userID == uuid.Nil {
		return uuid.Nil, false
	}

	return userID, true
}

func RequireUserID(
	ctx context.Context,
) (uuid.UUID, bool) {
	return UserID(ctx)
}
