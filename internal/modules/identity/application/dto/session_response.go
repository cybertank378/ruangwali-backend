package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionResponse struct {
	ID uuid.UUID

	UserAgent string

	IPAddress string

	ExpiresAt time.Time

	LastUsedAt *time.Time

	CreatedAt time.Time

	IsCurrent bool
}

type SessionListResponse struct {
	Items []SessionResponse
}
