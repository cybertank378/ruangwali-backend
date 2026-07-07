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
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(
	db *pgxpool.Pool,
) *UserRepository {
	if db == nil {
		panic("postgres user repository: db nil")
	}

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.User, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			status,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	return r.findOne(
		ctx,
		query,
		id,
	)
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email valueobject.Email,
) (*entity.User, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			status,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE LOWER(BTRIM(email)) = $1
		LIMIT 1
	`

	return r.findOne(
		ctx,
		query,
		email.String(),
	)
}

func (r *UserRepository) ExistsByEmail(
	ctx context.Context,
	email valueobject.Email,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE LOWER(BTRIM(email)) = $1
		)
	`

	var exists bool

	if err := r.db.QueryRow(
		ctx,
		query,
		email.String(),
	).Scan(
		&exists,
	); err != nil {
		return false, fmt.Errorf(
			"gagal memeriksa keberadaan email: %w",
			err,
		)
	}

	return exists, nil
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *entity.User,
) error {
	if user == nil {
		return errors.New(
			"user tidak boleh nil",
		)
	}

	const query = `
		INSERT INTO users (
			id,
			email,
			password_hash,
			status,
			last_login_at,
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
			$7
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash(),
		user.Status().String(),
		user.LastLoginAt(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return identitydomain.ErrEmailAlreadyExists
		}

		return fmt.Errorf(
			"gagal membuat user: %w",
			err,
		)
	}

	return nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user *entity.User,
) error {
	if user == nil {
		return errors.New(
			"user tidak boleh nil",
		)
	}

	const query = `
		UPDATE users
		SET
			email = $2,
			password_hash = $3,
			status = $4,
			last_login_at = $5,
			updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash(),
		user.Status().String(),
		user.LastLoginAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return identitydomain.ErrEmailAlreadyExists
		}

		return fmt.Errorf(
			"gagal memperbarui user: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return identitydomain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (*entity.User, error) {
	var row userRow

	err := r.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&row.ID,
		&row.Email,
		&row.PasswordHash,
		&row.Status,
		&row.LastLoginAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, identitydomain.ErrUserNotFound
		}

		return nil, fmt.Errorf(
			"gagal mengambil user: %w",
			err,
		)
	}

	return row.toEntity()
}
