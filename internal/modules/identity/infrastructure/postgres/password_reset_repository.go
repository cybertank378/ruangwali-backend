package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	identitydomain "github.com/ruangwali/internal/modules/identity/domain"
	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type PasswordResetRepository struct {
	pool *pgxpool.Pool
}

func NewPasswordResetRepository(
	pool *pgxpool.Pool,
) *PasswordResetRepository {
	if pool == nil {
		panic(
			"password reset repository: pool nil",
		)
	}

	return &PasswordResetRepository{
		pool: pool,
	}
}

func (r *PasswordResetRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.PasswordReset, error) {
	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			used_at,
			revoked_at,
			replaced_by_id,
			created_at,
			updated_at
		FROM password_reset_tokens
		WHERE id = $1
		LIMIT 1
	`

	row, err := scanPasswordResetRow(
		r.pool.QueryRow(
			ctx,
			query,
			id,
		),
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil,
				identitydomain.ErrPasswordResetNotFound
		}

		return nil, fmt.Errorf(
			"gagal mencari password reset berdasarkan ID: %w",
			err,
		)
	}

	passwordReset, err := row.toEntity()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal merehidrasi password reset: %w",
			err,
		)
	}

	return passwordReset, nil
}

func (r *PasswordResetRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash []byte,
) (*entity.PasswordReset, error) {
	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			used_at,
			revoked_at,
			replaced_by_id,
			created_at,
			updated_at
		FROM password_reset_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	row, err := scanPasswordResetRow(
		r.pool.QueryRow(
			ctx,
			query,
			tokenHash,
		),
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil,
				identitydomain.ErrPasswordResetTokenNotFound
		}

		return nil, fmt.Errorf(
			"gagal mencari password reset berdasarkan token hash: %w",
			err,
		)
	}

	passwordReset, err := row.toEntity()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal merehidrasi password reset: %w",
			err,
		)
	}

	return passwordReset, nil
}

func (r *PasswordResetRepository) FindActiveByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.PasswordReset, error) {
	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			used_at,
			revoked_at,
			replaced_by_id,
			created_at,
			updated_at
		FROM password_reset_tokens
		WHERE user_id = $1
		  AND used_at IS NULL
		  AND revoked_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	row, err := scanPasswordResetRow(
		r.pool.QueryRow(
			ctx,
			query,
			userID,
		),
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil,
				identitydomain.ErrPasswordResetNotFound
		}

		return nil, fmt.Errorf(
			"gagal mencari password reset aktif berdasarkan user ID: %w",
			err,
		)
	}

	passwordReset, err := row.toEntity()
	if err != nil {
		return nil, fmt.Errorf(
			"gagal merehidrasi password reset aktif: %w",
			err,
		)
	}

	return passwordReset, nil
}

func (r *PasswordResetRepository) Create(
	ctx context.Context,
	passwordReset *entity.PasswordReset,
) error {
	if passwordReset == nil {
		return errors.New(
			"password reset wajib tersedia",
		)
	}

	row := passwordResetRowFromEntity(
		passwordReset,
	)

	const query = `
		INSERT INTO password_reset_tokens (
			id,
			user_id,
			token_hash,
			expires_at,
			used_at,
			revoked_at,
			replaced_by_id,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		row.ID,
		row.UserID,
		row.TokenHash,
		row.ExpiresAt,
		row.UsedAt,
		row.RevokedAt,
		row.ReplacedByID,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat password reset: %w",
			err,
		)
	}

	return nil
}

func (r *PasswordResetRepository) Update(
	ctx context.Context,
	passwordReset *entity.PasswordReset,
) error {
	if passwordReset == nil {
		return errors.New(
			"password reset wajib tersedia",
		)
	}

	row := passwordResetRowFromEntity(
		passwordReset,
	)

	const query = `
		UPDATE password_reset_tokens
		SET
			used_at = $2,
			revoked_at = $3,
			replaced_by_id = $4,
			updated_at = $5
		WHERE id = $1
	`

	commandTag, err := r.pool.Exec(
		ctx,
		query,
		row.ID,
		row.UsedAt,
		row.RevokedAt,
		row.ReplacedByID,
		row.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal memperbarui password reset: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return identitydomain.ErrPasswordResetNotFound
	}

	return nil
}

type passwordResetScanner interface {
	Scan(
		dest ...any,
	) error
}

func scanPasswordResetRow(
	scanner passwordResetScanner,
) (passwordResetRow, error) {
	var row passwordResetRow

	err := scanner.Scan(
		&row.ID,
		&row.UserID,
		&row.TokenHash,
		&row.ExpiresAt,
		&row.UsedAt,
		&row.RevokedAt,
		&row.ReplacedByID,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return passwordResetRow{}, err
	}

	return row, nil
}
