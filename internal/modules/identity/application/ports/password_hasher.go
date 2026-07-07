// Package ports FilePath: /internal/modules/identity/application/ports/password_hasher.go
package ports

type PasswordHasher interface {
	Hash(
		password string,
	) (string, error)

	Verify(
		password string,
		encodedHash string,
	) (bool, error)
}
