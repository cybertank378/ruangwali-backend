package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Algorithm = "argon2id"

	defaultArgon2Memory      uint32 = 64 * 1024
	defaultArgon2Iterations  uint32 = 3
	defaultArgon2Parallelism uint8  = 2

	defaultArgon2SaltLength uint32 = 16
	defaultArgon2KeyLength  uint32 = 32
)

var (
	ErrPasswordEmpty = errors.New(
		"password tidak boleh kosong",
	)

	ErrInvalidPasswordHash = errors.New(
		"password hash tidak valid",
	)
)

type PasswordHasher struct {
	memory uint32

	iterations uint32

	parallelism uint8

	saltLength uint32
	keyLength  uint32
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		memory: defaultArgon2Memory,

		iterations: defaultArgon2Iterations,

		parallelism: defaultArgon2Parallelism,

		saltLength: defaultArgon2SaltLength,
		keyLength:  defaultArgon2KeyLength,
	}
}

func (h *PasswordHasher) Hash(
	password string,
) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}

	salt := make(
		[]byte,
		h.saltLength,
	)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf(
			"gagal membuat password salt: %w",
			err,
		)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(
		salt,
	)

	encodedHash := base64.RawStdEncoding.EncodeToString(
		hash,
	)

	return fmt.Sprintf(
		"$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2Algorithm,
		argon2.Version,
		h.memory,
		h.iterations,
		h.parallelism,
		encodedSalt,
		encodedHash,
	), nil
}

func (h *PasswordHasher) Verify(
	password string,
	encodedHash string,
) (bool, error) {
	if password == "" {
		return false, nil
	}

	params, salt, expectedHash, err := parseArgon2Hash(
		encodedHash,
	)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		uint32(len(expectedHash)),
	)

	if len(actualHash) != len(expectedHash) {
		return false, nil
	}

	return subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) == 1, nil
}

type argon2Params struct {
	memory uint32

	iterations uint32

	parallelism uint8
}

func parseArgon2Hash(
	encodedHash string,
) (
	argon2Params,
	[]byte,
	[]byte,
	error,
) {
	parts := strings.Split(
		encodedHash,
		"$",
	)

	if len(parts) != 6 {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	if parts[1] != argon2Algorithm {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	versionPart := strings.TrimPrefix(
		parts[2],
		"v=",
	)

	version, err := strconv.Atoi(
		versionPart,
	)
	if err != nil ||
		version != argon2.Version {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	var params argon2Params

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.memory,
		&params.iterations,
		&params.parallelism,
	); err != nil {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	if params.memory == 0 ||
		params.iterations == 0 ||
		params.parallelism == 0 {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(
		parts[4],
	)
	if err != nil ||
		len(salt) == 0 {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(
		parts[5],
	)
	if err != nil ||
		len(hash) == 0 {
		return argon2Params{},
			nil,
			nil,
			ErrInvalidPasswordHash
	}

	return params,
		salt,
		hash,
		nil
}
