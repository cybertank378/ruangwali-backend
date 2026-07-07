package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/domain/entity"
)

type stubSessionScanner struct {
	scanFunc func(
		dest ...any,
	) error
}

func (s stubSessionScanner) Scan(
	dest ...any,
) error {
	if s.scanFunc == nil {
		panic(
			"stubSessionScanner.Scan tidak dikonfigurasi",
		)
	}

	return s.scanFunc(
		dest...,
	)
}

func TestScanSessionRow(
	t *testing.T,
) {
	t.Run(
		"berhasil memindai seluruh kolom session",
		func(t *testing.T) {
			sessionID := uuid.New()
			userID := uuid.New()
			familyID := uuid.New()
			replacedByID := uuid.New()

			tokenHash := []byte{
				1,
				2,
				3,
				4,
			}

			userAgent := "Mozilla/5.0"
			ipAddress := "127.0.0.1"

			createdAt := time.Date(
				2026,
				time.July,
				7,
				10,
				0,
				0,
				0,
				time.UTC,
			)

			updatedAt := createdAt.Add(
				30 * time.Minute,
			)

			expiresAt := createdAt.Add(
				24 * time.Hour,
			)

			lastUsedAt := createdAt.Add(
				10 * time.Minute,
			)

			revokedAt := createdAt.Add(
				20 * time.Minute,
			)

			revokeReason := "USER_REQUESTED"

			scanner := stubSessionScanner{
				scanFunc: func(
					dest ...any,
				) error {
					if len(dest) != 13 {
						t.Fatalf(
							"jumlah destination = %d, mengharapkan 13",
							len(dest),
						)
					}

					*(dest[0].(*uuid.UUID)) =
						sessionID

					*(dest[1].(*uuid.UUID)) =
						userID

					*(dest[2].(*uuid.UUID)) =
						familyID

					*(dest[3].(*[]byte)) =
						append(
							[]byte(nil),
							tokenHash...,
						)

					*(dest[4].(**string)) =
						&userAgent

					*(dest[5].(**string)) =
						&ipAddress

					*(dest[6].(*time.Time)) =
						expiresAt

					*(dest[7].(**time.Time)) =
						&lastUsedAt

					*(dest[8].(**time.Time)) =
						&revokedAt

					*(dest[9].(**string)) =
						&revokeReason

					*(dest[10].(**uuid.UUID)) =
						&replacedByID

					*(dest[11].(*time.Time)) =
						createdAt

					*(dest[12].(*time.Time)) =
						updatedAt

					return nil
				},
			}

			row, err := scanSessionRow(
				scanner,
			)
			if err != nil {
				t.Fatalf(
					"mengharapkan nil error, mendapatkan %v",
					err,
				)
			}

			if row.ID != sessionID {
				t.Fatalf(
					"ID = %s, mengharapkan %s",
					row.ID,
					sessionID,
				)
			}

			if row.UserID != userID {
				t.Fatalf(
					"UserID = %s, mengharapkan %s",
					row.UserID,
					userID,
				)
			}

			if row.FamilyID != familyID {
				t.Fatalf(
					"FamilyID = %s, mengharapkan %s",
					row.FamilyID,
					familyID,
				)
			}

			if string(row.TokenHash) !=
				string(tokenHash) {
				t.Fatal(
					"TokenHash tidak sesuai",
				)
			}

			if row.UserAgent == nil {
				t.Fatal(
					"UserAgent tidak boleh nil",
				)
			}

			if *row.UserAgent != userAgent {
				t.Fatalf(
					"UserAgent = %q, mengharapkan %q",
					*row.UserAgent,
					userAgent,
				)
			}

			if row.IPAddress == nil {
				t.Fatal(
					"IPAddress tidak boleh nil",
				)
			}

			if *row.IPAddress != ipAddress {
				t.Fatalf(
					"IPAddress = %q, mengharapkan %q",
					*row.IPAddress,
					ipAddress,
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

			if row.LastUsedAt == nil {
				t.Fatal(
					"LastUsedAt tidak boleh nil",
				)
			}

			if !row.LastUsedAt.Equal(
				lastUsedAt,
			) {
				t.Fatalf(
					"LastUsedAt = %s, mengharapkan %s",
					*row.LastUsedAt,
					lastUsedAt,
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

			if row.RevokeReason == nil {
				t.Fatal(
					"RevokeReason tidak boleh nil",
				)
			}

			if *row.RevokeReason !=
				revokeReason {
				t.Fatalf(
					"RevokeReason = %q, mengharapkan %q",
					*row.RevokeReason,
					revokeReason,
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
		"mengembalikan scanner error",
		func(t *testing.T) {
			expectedErr := errors.New(
				"scan gagal",
			)

			scanner := stubSessionScanner{
				scanFunc: func(
					dest ...any,
				) error {
					return expectedErr
				},
			}

			row, err := scanSessionRow(
				scanner,
			)

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
					"row ID = %s, mengharapkan uuid.Nil",
					row.ID,
				)
			}
		},
	)
}

func TestSessionRow_ToEntity(
	t *testing.T,
) {
	t.Run(
		"berhasil rehydrate session lengkap",
		func(t *testing.T) {
			sessionID := uuid.New()
			userID := uuid.New()
			familyID := uuid.New()
			replacedByID := uuid.New()

			tokenHash := []byte{
				10,
				20,
				30,
				40,
			}

			userAgent := "Mozilla/5.0"
			ipAddress := "192.168.1.10"

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

			lastUsedAt := createdAt.Add(
				15 * time.Minute,
			)

			revokedAt := createdAt.Add(
				30 * time.Minute,
			)

			updatedAt := revokedAt

			expiresAt := createdAt.Add(
				24 * time.Hour,
			)

			revokeReason := "ROTATED"

			row := sessionRow{
				ID: sessionID,

				UserID: userID,

				FamilyID: familyID,

				TokenHash: tokenHash,

				UserAgent: &userAgent,
				IPAddress: &ipAddress,

				ExpiresAt: expiresAt,

				LastUsedAt: &lastUsedAt,

				RevokedAt: &revokedAt,

				RevokeReason: &revokeReason,

				ReplacedByID: &replacedByID,

				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}

			session, err := row.toEntity()
			if err != nil {
				t.Fatalf(
					"mengharapkan nil error, mendapatkan %v",
					err,
				)
			}

			if session == nil {
				t.Fatal(
					"session tidak boleh nil",
				)
			}

			if session.ID() != sessionID {
				t.Fatalf(
					"ID = %s, mengharapkan %s",
					session.ID(),
					sessionID,
				)
			}

			if session.UserID() != userID {
				t.Fatalf(
					"UserID = %s, mengharapkan %s",
					session.UserID(),
					userID,
				)
			}

			if session.FamilyID() !=
				familyID {
				t.Fatalf(
					"FamilyID = %s, mengharapkan %s",
					session.FamilyID(),
					familyID,
				)
			}

			if string(session.TokenHash()) !=
				string(tokenHash) {
				t.Fatal(
					"TokenHash tidak sesuai",
				)
			}

			if session.Fingerprint().
				UserAgent() != userAgent {
				t.Fatalf(
					"UserAgent = %q, mengharapkan %q",
					session.Fingerprint().
						UserAgent(),
					userAgent,
				)
			}

			if session.Fingerprint().
				IPAddress() != ipAddress {
				t.Fatalf(
					"IPAddress = %q, mengharapkan %q",
					session.Fingerprint().
						IPAddress(),
					ipAddress,
				)
			}

			if !session.ExpiresAt().Equal(
				expiresAt,
			) {
				t.Fatalf(
					"ExpiresAt = %s, mengharapkan %s",
					session.ExpiresAt(),
					expiresAt,
				)
			}

			if session.LastUsedAt() == nil {
				t.Fatal(
					"LastUsedAt tidak boleh nil",
				)
			}

			if !session.LastUsedAt().Equal(
				lastUsedAt,
			) {
				t.Fatalf(
					"LastUsedAt = %s, mengharapkan %s",
					*session.LastUsedAt(),
					lastUsedAt,
				)
			}

			if session.RevokedAt() == nil {
				t.Fatal(
					"RevokedAt tidak boleh nil",
				)
			}

			if !session.RevokedAt().Equal(
				revokedAt,
			) {
				t.Fatalf(
					"RevokedAt = %s, mengharapkan %s",
					*session.RevokedAt(),
					revokedAt,
				)
			}

			if session.RevokeReason() == nil {
				t.Fatal(
					"RevokeReason tidak boleh nil",
				)
			}

			if *session.RevokeReason() !=
				revokeReason {
				t.Fatalf(
					"RevokeReason = %q, mengharapkan %q",
					*session.RevokeReason(),
					revokeReason,
				)
			}

			if session.ReplacedByID() == nil {
				t.Fatal(
					"ReplacedByID tidak boleh nil",
				)
			}

			if *session.ReplacedByID() !=
				replacedByID {
				t.Fatalf(
					"ReplacedByID = %s, mengharapkan %s",
					*session.ReplacedByID(),
					replacedByID,
				)
			}

			if !session.CreatedAt().Equal(
				createdAt,
			) {
				t.Fatalf(
					"CreatedAt = %s, mengharapkan %s",
					session.CreatedAt(),
					createdAt,
				)
			}

			if !session.UpdatedAt().Equal(
				updatedAt,
			) {
				t.Fatalf(
					"UpdatedAt = %s, mengharapkan %s",
					session.UpdatedAt(),
					updatedAt,
				)
			}
		},
	)

	t.Run(
		"berhasil rehydrate session dengan nullable field kosong",
		func(t *testing.T) {
			sessionID := uuid.New()
			userID := uuid.New()
			familyID := uuid.New()

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

			row := sessionRow{
				ID: sessionID,

				UserID: userID,

				FamilyID: familyID,

				TokenHash: []byte{
					1,
					2,
					3,
				},

				UserAgent: nil,
				IPAddress: nil,

				ExpiresAt: createdAt.Add(
					24 * time.Hour,
				),

				LastUsedAt: nil,

				RevokedAt: nil,

				RevokeReason: nil,

				ReplacedByID: nil,

				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}

			session, err := row.toEntity()
			if err != nil {
				t.Fatalf(
					"mengharapkan nil error, mendapatkan %v",
					err,
				)
			}

			if session == nil {
				t.Fatal(
					"session tidak boleh nil",
				)
			}

			if !session.Fingerprint().
				IsZero() {
				t.Fatal(
					"fingerprint harus zero",
				)
			}

			if session.LastUsedAt() != nil {
				t.Fatal(
					"LastUsedAt harus nil",
				)
			}

			if session.RevokedAt() != nil {
				t.Fatal(
					"RevokedAt harus nil",
				)
			}

			if session.RevokeReason() != nil {
				t.Fatal(
					"RevokeReason harus nil",
				)
			}

			if session.ReplacedByID() != nil {
				t.Fatal(
					"ReplacedByID harus nil",
				)
			}
		},
	)

	t.Run(
		"mengembalikan error ketika IP address tidak valid",
		func(t *testing.T) {
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

			row := sessionRow{
				ID: uuid.New(),

				UserID: uuid.New(),

				FamilyID: uuid.New(),

				TokenHash: []byte{
					1,
					2,
					3,
				},

				IPAddress: new("alamat-ip-tidak-valid"),

				ExpiresAt: createdAt.Add(
					24 * time.Hour,
				),

				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}

			session, err := row.toEntity()

			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			if session != nil {
				t.Fatal(
					"session harus nil",
				)
			}
		},
	)

	t.Run(
		"mengembalikan invalid session ID",
		func(t *testing.T) {
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

			row := sessionRow{
				ID: uuid.Nil,

				UserID: uuid.New(),

				FamilyID: uuid.New(),

				TokenHash: []byte{
					1,
					2,
					3,
				},

				ExpiresAt: createdAt.Add(
					24 * time.Hour,
				),

				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}

			session, err := row.toEntity()

			if session != nil {
				t.Fatal(
					"session harus nil",
				)
			}

			if !errors.Is(
				err,
				entity.ErrInvalidSessionID,
			) {
				t.Fatalf(
					"error = %v, mengharapkan ErrInvalidSessionID",
					err,
				)
			}
		},
	)

	t.Run(
		"mengembalikan invalid token hash",
		func(t *testing.T) {
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

			row := sessionRow{
				ID: uuid.New(),

				UserID: uuid.New(),

				FamilyID: uuid.New(),

				TokenHash: nil,

				ExpiresAt: createdAt.Add(
					24 * time.Hour,
				),

				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}

			session, err := row.toEntity()

			if session != nil {
				t.Fatal(
					"session harus nil",
				)
			}

			if !errors.Is(
				err,
				entity.ErrInvalidTokenHash,
			) {
				t.Fatalf(
					"error = %v, mengharapkan ErrInvalidTokenHash",
					err,
				)
			}
		},
	)
}
