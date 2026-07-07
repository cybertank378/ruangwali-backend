package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(
	db *pgxpool.Pool,
) *SessionRepository {
	if db == nil {
		panic("postgres session repository: db nil")
	}

	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.RefreshSession, error) {
	const query = `
		SELECT
			id,
			user_id,
			family_id,
			token_hash,
			user_agent,
			host(ip_address),
			expires_at,
			last_used_at,
			revoked_at,
			revoke_reason,
			replaced_by_id,
			created_at,
			updated_at
		FROM refresh_sessions
		WHERE id = $1
		LIMIT 1
	`

	return r.findOne(
		ctx,
		query,
		id,
	)
}

func (r *SessionRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash []byte,
) (*entity.RefreshSession, error) {
	const query = `
		SELECT
			id,
			user_id,
			family_id,
			token_hash,
			user_agent,
			host(ip_address),
			expires_at,
			last_used_at,
			revoked_at,
			revoke_reason,
			replaced_by_id,
			created_at,
			updated_at
		FROM refresh_sessions
		WHERE token_hash = $1
		LIMIT 1
	`

	return r.findOne(
		ctx,
		query,
		tokenHash,
	)
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *entity.RefreshSession,
) error {
	if session == nil {
		return errors.New(
			"refresh session tidak boleh nil",
		)
	}

	const query = `
		INSERT INTO refresh_sessions (
			id,
			user_id,
			family_id,
			token_hash,
			user_agent,
			ip_address,
			expires_at,
			last_used_at,
			revoked_at,
			revoke_reason,
			replaced_by_id,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			NULLIF($5, ''),
			NULLIF($6, '')::INET,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		)
	`

	fingerprint := session.Fingerprint()

	_, err := r.db.Exec(
		ctx,
		query,
		session.ID(),
		session.UserID(),
		session.FamilyID(),
		session.TokenHash(),
		fingerprint.UserAgent(),
		fingerprint.IPAddress(),
		session.ExpiresAt(),
		session.LastUsedAt(),
		session.RevokedAt(),
		session.RevokeReason(),
		session.ReplacedByID(),
		session.CreatedAt(),
		session.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat refresh session: %w",
			err,
		)
	}

	return nil
}

func (r *SessionRepository) Update(
	ctx context.Context,
	session *entity.RefreshSession,
) error {
	if session == nil {
		return errors.New(
			"refresh session tidak boleh nil",
		)
	}

	const query = `
		UPDATE refresh_sessions
		SET
			last_used_at = $2,
			revoked_at = $3,
			revoke_reason = $4,
			replaced_by_id = $5,
			updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		session.ID(),
		session.LastUsedAt(),
		session.RevokedAt(),
		session.RevokeReason(),
		session.ReplacedByID(),
		session.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal memperbarui refresh session: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return identitydomain.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) RevokeByUserID(
	ctx context.Context,
	userID uuid.UUID,
	reason string,
	revokedAt time.Time,
) error {
	const query = `
		UPDATE refresh_sessions
		SET
			revoked_at = $2,
			revoke_reason = NULLIF(BTRIM($3), ''),
			updated_at = $2
		WHERE user_id = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		revokedAt.UTC(),
		reason,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal mencabut session user: %w",
			err,
		)
	}

	return nil
}

func (r *SessionRepository) RevokeByFamilyID(
	ctx context.Context,
	familyID uuid.UUID,
	reason string,
	revokedAt time.Time,
) error {
	const query = `
		UPDATE refresh_sessions
		SET
			revoked_at = $2,
			revoke_reason = NULLIF(BTRIM($3), ''),
			updated_at = $2
		WHERE family_id = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(
		ctx,
		query,
		familyID,
		revokedAt.UTC(),
		reason,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal mencabut session family: %w",
			err,
		)
	}

	return nil
}

func (r *SessionRepository) DeleteExpiredBefore(
	ctx context.Context,
	before time.Time,
) (int64, error) {
	const query = `
		DELETE FROM refresh_sessions
		WHERE expires_at < $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		before.UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf(
			"gagal menghapus session kedaluwarsa: %w",
			err,
		)
	}

	return result.RowsAffected(), nil
}

func (r *SessionRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (*entity.RefreshSession, error) {
	var row sessionRow

	err := r.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&row.ID,
		&row.UserID,
		&row.FamilyID,
		&row.TokenHash,
		&row.UserAgent,
		&row.IPAddress,
		&row.ExpiresAt,
		&row.LastUsedAt,
		&row.RevokedAt,
		&row.RevokeReason,
		&row.ReplacedByID,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, identitydomain.ErrSessionNotFound
		}

		return nil, fmt.Errorf(
			"gagal mengambil refresh session: %w",
			err,
		)
	}

	return row.toEntity()
}
