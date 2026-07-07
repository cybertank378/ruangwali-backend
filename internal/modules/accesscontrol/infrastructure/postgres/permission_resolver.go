package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionResolver struct {
	db *pgxpool.Pool
}

func NewPermissionResolver(db *pgxpool.Pool) *PermissionResolver {
	return &PermissionResolver{db: db}
}

func (r *PermissionResolver) ResolveEffectivePermissions(
	ctx context.Context,
	userID string,
	tenantID string,
) ([]string, error) {
	const query = `
SELECT DISTINCT p.code
FROM memberships m
JOIN membership_roles mr ON mr.membership_id = m.id
JOIN roles r ON r.id = mr.role_id AND r.is_active = TRUE
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE m.user_id = $1
  AND m.tenant_id = $2
  AND m.status = 'ACTIVE'
  AND (r.tenant_id IS NULL OR r.tenant_id = $2)
ORDER BY p.code`

	rows, err := r.db.Query(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		result = append(result, code)
	}
	return result, rows.Err()
}
