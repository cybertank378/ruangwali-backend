package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type stubPasswordResetScanner struct {
	scanFunc func(
		dest ...any,
	) error
}

func (s stubPasswordResetScanner) Scan(
	dest ...any,
) error {
	if s.scanFunc == nil {
		panic(
			"stubPasswordResetScanner.Scan tidak dikonfigurasi",
		)
	}

	return s.scanFunc(
		dest...,
	)
}

func TestNewPasswordResetRepository(
	t *testing.T,
) {
	t.Run(
		"positif berhasil membuat repository",
		func(t *testing.T) {
			// Arrange
			repository := &PasswordResetRepository{}

			// Act
			actual := repository

			// Assert
			if actual == nil {
				t.Fatal(
					"repository tidak boleh nil",
				)
			}
		},
	)

	t.Run(
		"negatif panic ketika pool nil",
		func(t *testing.T) {
			// Arrange
			var recovered any

			func() {
				defer func() {
					recovered = recover()
				}()

				// Act
				NewPasswordResetRepository(
					nil,
				)
			}()

			// Assert
			if recovered == nil {
				t.Fatal(
					"mengharapkan panic, tetapi tidak terjadi",
				)
			}

			actualMessage, ok := recovered.(string)
			if !ok {
				t.Fatalf(
					"panic memiliki tipe %T, mengharapkan string",
					recovered,
				)
			}

			const expectedMessage = "password reset repository: pool nil"

			if actualMessage != expectedMessage {
				t.Fatalf(
					"panic message = %q, mengharapkan %q",
					actualMessage,
					expectedMessage,
				)
			}
		},
	)
}

func TestScanPasswordResetRow(
	t *testing.T,
) {
	t.Run(
		"positif berhasil memindai seluruh kolom",
		func(t *testing.T) {
			// Arrange
			id := uuid.New()

			userID := uuid.New()

			replacedByID := uuid.New()

			tokenHash := []byte{
				10,
				20,
				30,
				40,
			}

			createdAt := time.Date(
				2026,
				time.July,
				7,
				8,
				0,
				0,
				0,
				time.UTC,
			)

			expiresAt := createdAt.Add(
				30 * time.Minute,
			)

			usedAt := createdAt.Add(
				10 * time.Minute,
			)

			revokedAt := createdAt.Add(
				20 * time.Minute,
			)

			updatedAt := createdAt.Add(
				25 * time.Minute,
			)

			scanner := stubPasswordResetScanner{
				scanFunc: func(
					dest ...any,
				) error {
					if len(dest) != 9 {
						t.Fatalf(
							"jumlah destination = %d, mengharapkan 9",
							len(dest),
						)
					}

					*(dest[0].(*uuid.UUID)) =
						id

					*(dest[1].(*uuid.UUID)) =
						userID

					*(dest[2].(*[]byte)) =
						append(
							[]byte(nil),
							tokenHash...,
						)

					*(dest[3].(*time.Time)) =
						expiresAt

					*(dest[4].(**time.Time)) =
						&usedAt

					*(dest[5].(**time.Time)) =
						&revokedAt

					*(dest[6].(**uuid.UUID)) =
						&replacedByID

					*(dest[7].(*time.Time)) =
						createdAt

					*(dest[8].(*time.Time)) =
						updatedAt

					return nil
				},
			}

			// Act
			row, err := scanPasswordResetRow(
				scanner,
			)

			// Assert
			if err != nil {
				t.Fatalf(
					"mengharapkan nil error, mendapatkan %v",
					err,
				)
			}

			if row.ID != id {
				t.Fatalf(
					"ID = %s, mengharapkan %s",
					row.ID,
					id,
				)
			}

			if row.UserID != userID {
				t.Fatalf(
					"UserID = %s, mengharapkan %s",
					row.UserID,
					userID,
				)
			}

			if string(row.TokenHash) !=
				string(tokenHash) {
				t.Fatalf(
					"TokenHash = %v, mengharapkan %v",
					row.TokenHash,
					tokenHash,
				)
			}

			if !row.ExpiresAt.Equal(
				expiresAt,
			) {
				t.Fatalf(
					"ExpiresAt = %s, mengharapkan %s",
					row.ExpiresAt,
					expiresAt,
				)
			}

			if row.UsedAt == nil {
				t.Fatal(
					"UsedAt tidak boleh nil",
				)
			}

			if !row.UsedAt.Equal(
				usedAt,
			) {
				t.Fatalf(
					"UsedAt = %s, mengharapkan %s",
					*row.UsedAt,
					usedAt,
				)
			}

			if row.RevokedAt == nil {
				t.Fatal(
					"RevokedAt tidak boleh nil",
				)
			}

			if !row.RevokedAt.Equal(
				revokedAt,
			) {
				t.Fatalf(
					"RevokedAt = %s, mengharapkan %s",
					*row.RevokedAt,
					revokedAt,
				)
			}

			if row.ReplacedByID == nil {
				t.Fatal(
					"ReplacedByID tidak boleh nil",
				)
			}

			if *row.ReplacedByID !=
				replacedByID {
				t.Fatalf(
					"ReplacedByID = %s, mengharapkan %s",
					*row.ReplacedByID,
					replacedByID,
				)
			}

			if !row.CreatedAt.Equal(
				createdAt,
			) {
				t.Fatalf(
					"CreatedAt = %s, mengharapkan %s",
					row.CreatedAt,
					createdAt,
				)
			}

			if !row.UpdatedAt.Equal(
				updatedAt,
			) {
				t.Fatalf(
					"UpdatedAt = %s, mengharapkan %s",
					row.UpdatedAt,
					updatedAt,
				)
			}
		},
	)

	t.Run(
		"positif berhasil memindai nullable field kosong",
		func(t *testing.T) {
			// Arrange
			id := uuid.New()

			userID := uuid.New()

			tokenHash := []byte{
				1,
				2,
				3,
			}

			createdAt := time.Date(
				2026,
				time.July,
				7,
				8,
				0,
				0,
				0,
				time.UTC,
			)

			expiresAt := createdAt.Add(
				30 * time.Minute,
			)

			updatedAt := createdAt

			scanner := stubPasswordResetScanner{
				scanFunc: func(
					dest ...any,
				) error {
					if len(dest) != 9 {
						t.Fatalf(
							"jumlah destination = %d, mengharapkan 9",
							len(dest),
						)
					}

					*(dest[0].(*uuid.UUID)) =
						id

					*(dest[1].(*uuid.UUID)) =
						userID

					*(dest[2].(*[]byte)) =
						append(
							[]byte(nil),
							tokenHash...,
						)

					*(dest[3].(*time.Time)) =
						expiresAt

					*(dest[4].(**time.Time)) =
						nil

					*(dest[5].(**time.Time)) =
						nil

					*(dest[6].(**uuid.UUID)) =
						nil

					*(dest[7].(*time.Time)) =
						createdAt

					*(dest[8].(*time.Time)) =
						updatedAt

					return nil
				},
			}

			// Act
			row, err := scanPasswordResetRow(
				scanner,
			)

			// Assert
			if err != nil {
				t.Fatalf(
					"mengharapkan nil error, mendapatkan %v",
					err,
				)
			}

			if row.ID != id {
				t.Fatalf(
					"ID = %s, mengharapkan %s",
					row.ID,
					id,
				)
			}

			if row.UserID != userID {
				t.Fatalf(
					"UserID = %s, mengharapkan %s",
					row.UserID,
					userID,
				)
			}

			if string(row.TokenHash) !=
				string(tokenHash) {
				t.Fatalf(
					"TokenHash = %v, mengharapkan %v",
					row.TokenHash,
					tokenHash,
				)
			}

			if !row.ExpiresAt.Equal(
				expiresAt,
			) {
				t.Fatalf(
					"ExpiresAt = %s, mengharapkan %s",
					row.ExpiresAt,
					expiresAt,
				)
			}

			if row.UsedAt != nil {
				t.Fatalf(
					"UsedAt = %v, mengharapkan nil",
					row.UsedAt,
				)
			}

			if row.RevokedAt != nil {
				t.Fatalf(
					"RevokedAt = %v, mengharapkan nil",
					row.RevokedAt,
				)
			}

			if row.ReplacedByID != nil {
				t.Fatalf(
					"ReplacedByID = %v, mengharapkan nil",
					row.ReplacedByID,
				)
			}

			if !row.CreatedAt.Equal(
				createdAt,
			) {
				t.Fatalf(
					"CreatedAt = %s, mengharapkan %s",
					row.CreatedAt,
					createdAt,
				)
			}

			if !row.UpdatedAt.Equal(
				updatedAt,
			) {
				t.Fatalf(
					"UpdatedAt = %s, mengharapkan %s",
					row.UpdatedAt,
					updatedAt,
				)
			}
		},
	)

	t.Run(
		"negatif mengembalikan scanner error",
		func(t *testing.T) {
			// Arrange
			expectedErr := errors.New(
				"database scan gagal",
			)

			scanner := stubPasswordResetScanner{
				scanFunc: func(
					dest ...any,
				) error {
					return expectedErr
				},
			}

			// Act
			row, err := scanPasswordResetRow(
				scanner,
			)

			// Assert
			if !errors.Is(
				err,
				expectedErr,
			) {
				t.Fatalf(
					"error = %v, mengharapkan %v",
					err,
					expectedErr,
				)
			}

			if row.ID != uuid.Nil {
				t.Fatalf(
					"ID = %s, mengharapkan uuid.Nil",
					row.ID,
				)
			}

			if row.UserID != uuid.Nil {
				t.Fatalf(
					"UserID = %s, mengharapkan uuid.Nil",
					row.UserID,
				)
			}

			if row.TokenHash != nil {
				t.Fatalf(
					"TokenHash = %v, mengharapkan nil",
					row.TokenHash,
				)
			}

			if !row.ExpiresAt.IsZero() {
				t.Fatalf(
					"ExpiresAt = %s, mengharapkan zero time",
					row.ExpiresAt,
				)
			}

			if row.UsedAt != nil {
				t.Fatalf(
					"UsedAt = %v, mengharapkan nil",
					row.UsedAt,
				)
			}

			if row.RevokedAt != nil {
				t.Fatalf(
					"RevokedAt = %v, mengharapkan nil",
					row.RevokedAt,
				)
			}

			if row.ReplacedByID != nil {
				t.Fatalf(
					"ReplacedByID = %v, mengharapkan nil",
					row.ReplacedByID,
				)
			}

			if !row.CreatedAt.IsZero() {
				t.Fatalf(
					"CreatedAt = %s, mengharapkan zero time",
					row.CreatedAt,
				)
			}

			if !row.UpdatedAt.IsZero() {
				t.Fatalf(
					"UpdatedAt = %s, mengharapkan zero time",
					row.UpdatedAt,
				)
			}
		},
	)

	t.Run(
		"negatif meneruskan error asli tanpa mengganti identitas error",
		func(t *testing.T) {
			// Arrange
			expectedErr := errors.New(
				"connection interrupted",
			)

			scanner := stubPasswordResetScanner{
				scanFunc: func(
					dest ...any,
				) error {
					return expectedErr
				},
			}

			// Act
			_, err := scanPasswordResetRow(
				scanner,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			if !errors.Is(
				err,
				expectedErr,
			) {
				t.Fatalf(
					"error = %v, mengharapkan error asli %v",
					err,
					expectedErr,
				)
			}
		},
	)
}

func TestPasswordResetRepository_Create(
	t *testing.T,
) {
	t.Run(
		"negatif mengembalikan error ketika password reset nil",
		func(t *testing.T) {
			// Arrange
			repository := &PasswordResetRepository{}

			// Act
			err := repository.Create(
				context.Background(),
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "password reset wajib tersedia"

			if err.Error() != expectedMessage {
				t.Fatalf(
					"error = %q, mengharapkan %q",
					err.Error(),
					expectedMessage,
				)
			}
		},
	)
}

func TestPasswordResetRepository_Update(
	t *testing.T,
) {
	t.Run(
		"negatif mengembalikan error ketika password reset nil",
		func(t *testing.T) {
			// Arrange
			repository := &PasswordResetRepository{}

			// Act
			err := repository.Update(
				context.Background(),
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "password reset wajib tersedia"

			if err.Error() != expectedMessage {
				t.Fatalf(
					"error = %q, mengharapkan %q",
					err.Error(),
					expectedMessage,
				)
			}
		},
	)
}
