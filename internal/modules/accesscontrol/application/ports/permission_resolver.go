package ports

import (
	"context"

	"github.com/google/uuid"
)

// PermissionResolver menyelesaikan permission efektif milik user.
type PermissionResolver interface {
	ResolveByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]string, error)
}
