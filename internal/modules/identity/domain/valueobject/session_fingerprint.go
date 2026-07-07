package valueobject

import (
	"fmt"
	"net"
	"strings"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

const maxUserAgentLength = 2048

var ErrInvalidSessionFingerprint = fmt.Errorf(
	"session fingerprint tidak valid: %w",
	shareddomain.ErrValidation,
)

type SessionFingerprint struct {
	userAgent string
	ipAddress string
}

func NewSessionFingerprint(
	userAgent string,
	ipAddress string,
) (SessionFingerprint, error) {
	normalizedUserAgent := strings.TrimSpace(
		userAgent,
	)

	normalizedIPAddress := strings.TrimSpace(
		ipAddress,
	)

	if len(normalizedUserAgent) > maxUserAgentLength {
		return SessionFingerprint{}, fmt.Errorf(
			"user agent melebihi %d karakter: %w",
			maxUserAgentLength,
			ErrInvalidSessionFingerprint,
		)
	}

	if normalizedIPAddress != "" {
		parsedIP := net.ParseIP(
			normalizedIPAddress,
		)

		if parsedIP == nil {
			return SessionFingerprint{}, fmt.Errorf(
				"IP address tidak valid: %w",
				ErrInvalidSessionFingerprint,
			)
		}

		normalizedIPAddress = parsedIP.String()
	}

	return SessionFingerprint{
		userAgent: normalizedUserAgent,
		ipAddress: normalizedIPAddress,
	}, nil
}

func (f SessionFingerprint) UserAgent() string {
	return f.userAgent
}

func (f SessionFingerprint) IPAddress() string {
	return f.ipAddress
}

func (f SessionFingerprint) IsZero() bool {
	return f.userAgent == "" &&
		f.ipAddress == ""
}
