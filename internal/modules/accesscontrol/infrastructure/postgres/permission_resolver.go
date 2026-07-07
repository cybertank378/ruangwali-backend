package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruangwali/internal/modules/accesscontrol/application/ports"
)

type PermissionResolver struct {
	db *pgxpool.Pool
}

func NewPermissionResolver(
	db *pgxpool.Pool,
) *PermissionResolver {
	if db == nil {
		panic(
			"postgres permission resolver: db nil",
		)
	}

	return &PermissionResolver{
		db: db,
	}
}

func (r *PermissionResolver) ResolveByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]string, error) {
	const query = `
		SELECT DISTINCT
			p.code
		FROM user_roles ur
		INNER JOIN roles r
			ON r.id = ur.role_id
		INNER JOIN role_permissions rp
			ON rp.role_id = r.id
		INNER JOIN permissions p
			ON p.id = rp.permission_id
		WHERE ur.user_id = $1
		  AND r.is_active = TRUE
		ORDER BY p.code ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil permission user: %w",
			err,
		)
	}
	defer rows.Close()

	permissions := make(
		[]string,
		0,
	)

	for rows.Next() {
		var permission string

		if err := rows.Scan(
			&permission,
		); err != nil {
			return nil, fmt.Errorf(
				"gagal memindai permission user: %w",
				err,
			)
		}

		permissions = append(
			permissions,
			permission,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil permission user: %w",
			err,
		)
	}

	return permissions, nil
}

var _ ports.PermissionResolver = (*PermissionResolver)(nil)
