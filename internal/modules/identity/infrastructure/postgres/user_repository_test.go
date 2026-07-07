package postgres

import (
	"context"
	"testing"
)

func TestNewUserRepository(
	t *testing.T,
) {
	t.Run(
		"negatif panic ketika database nil",
		func(t *testing.T) {
			// Arrange
			var recovered any

			// Act
			func() {
				defer func() {
					recovered = recover()
				}()

				NewUserRepository(
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

			const expectedMessage = "postgres user repository: db nil"

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

func TestUserRepository_Create(
	t *testing.T,
) {
	t.Run(
		"negatif mengembalikan error ketika user nil",
		func(t *testing.T) {
			// Arrange
			repository := &UserRepository{}

			ctx := context.Background()

			// Act
			err := repository.Create(
				ctx,
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "user tidak boleh nil"

			if err.Error() != expectedMessage {
				t.Fatalf(
					"error = %q, mengharapkan %q",
					err.Error(),
					expectedMessage,
				)
			}
		},
	)

	t.Run(
		"negatif guard clause berjalan sebelum akses database",
		func(t *testing.T) {
			// Arrange
			repository := &UserRepository{
				db: nil,
			}

			ctx := context.Background()

			// Act
			err := repository.Create(
				ctx,
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "user tidak boleh nil"

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

func TestUserRepository_Update(
	t *testing.T,
) {
	t.Run(
		"negatif mengembalikan error ketika user nil",
		func(t *testing.T) {
			// Arrange
			repository := &UserRepository{}

			ctx := context.Background()

			// Act
			err := repository.Update(
				ctx,
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "user tidak boleh nil"

			if err.Error() != expectedMessage {
				t.Fatalf(
					"error = %q, mengharapkan %q",
					err.Error(),
					expectedMessage,
				)
			}
		},
	)

	t.Run(
		"negatif guard clause berjalan sebelum akses database",
		func(t *testing.T) {
			// Arrange
			repository := &UserRepository{
				db: nil,
			}

			ctx := context.Background()

			// Act
			err := repository.Update(
				ctx,
				nil,
			)

			// Assert
			if err == nil {
				t.Fatal(
					"mengharapkan error, mendapatkan nil",
				)
			}

			const expectedMessage = "user tidak boleh nil"

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

func TestUserRepository_GuardClauses(
	t *testing.T,
) {
	tests := []struct {
		name string

		execute func(
			repository *UserRepository,
		) error

		expectedMessage string
	}{
		{
			name: "negatif create menolak user nil",

			execute: func(
				repository *UserRepository,
			) error {
				return repository.Create(
					context.Background(),
					nil,
				)
			},

			expectedMessage: "user tidak boleh nil",
		},
		{
			name: "negatif update menolak user nil",

			execute: func(
				repository *UserRepository,
			) error {
				return repository.Update(
					context.Background(),
					nil,
				)
			},

			expectedMessage: "user tidak boleh nil",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Arrange
				repository := &UserRepository{
					db: nil,
				}

				// Act
				err := tt.execute(
					repository,
				)

				// Assert
				if err == nil {
					t.Fatal(
						"mengharapkan error, mendapatkan nil",
					)
				}

				if err.Error() !=
					tt.expectedMessage {
					t.Fatalf(
						"error = %q, mengharapkan %q",
						err.Error(),
						tt.expectedMessage,
					)
				}
			},
		)
	}
}
