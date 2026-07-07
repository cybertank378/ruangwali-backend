package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, student *Student) error
	FindByID(ctx context.Context, tenantID, studentID uuid.UUID) (*Student, error)
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]Student, error)
}
