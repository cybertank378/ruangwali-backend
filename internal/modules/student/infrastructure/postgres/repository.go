package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	studentdomain "github.com/ruangwali/internal/modules/student/domain"
	shareddomain "github.com/ruangwali/internal/shared/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, s *studentdomain.Student) error {
	const query = `
INSERT INTO students (
	id, tenant_id, full_name, gender, religion,
	birth_place, birth_date, address, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := r.db.Exec(ctx, query,
		s.ID, s.TenantID, s.FullName, s.Gender, s.Religion,
		s.BirthPlace, s.BirthDate, s.Address, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByID(ctx context.Context, tenantID, studentID uuid.UUID) (*studentdomain.Student, error) {
	const query = `
SELECT id, tenant_id, full_name, gender, religion, birth_place,
       birth_date, address, created_at, updated_at
FROM students
WHERE tenant_id = $1 AND id = $2
LIMIT 1`

	var s studentdomain.Student
	err := r.db.QueryRow(ctx, query, tenantID, studentID).Scan(
		&s.ID, &s.TenantID, &s.FullName, &s.Gender, &s.Religion,
		&s.BirthPlace, &s.BirthDate, &s.Address, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, shareddomain.ErrNotFound
	}
	return &s, err
}

func (r *Repository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]studentdomain.Student, error) {
	const query = `
SELECT id, tenant_id, full_name, gender, religion, birth_place,
       birth_date, address, created_at, updated_at
FROM students
WHERE tenant_id = $1
ORDER BY full_name ASC
LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]studentdomain.Student, 0)
	for rows.Next() {
		var s studentdomain.Student
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.FullName, &s.Gender, &s.Religion,
			&s.BirthPlace, &s.BirthDate, &s.Address, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
