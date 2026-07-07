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
