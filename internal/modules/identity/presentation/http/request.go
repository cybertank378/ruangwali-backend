// =========================================================
// File: internal/modules/identity/presentation/http/request.go
// =========================================================

package http

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	stdhttp "net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/application/dto"
	"github.com/ruangwali/internal/modules/identity/domain/valueobject"
)

const maxRequestBodyBytes int64 = 1 << 20

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r loginRequest) toDTO(
	request *stdhttp.Request,
) dto.LoginRequest {
	return dto.LoginRequest{
		Email: strings.TrimSpace(
			r.Email,
		),

		Password: r.Password,

		UserAgent: request.UserAgent(),

		IPAddress: clientIPAddress(
			request,
		),
	}
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r logoutRequest) toDTO() dto.LogoutRequest {
	return dto.LogoutRequest{
		RefreshToken: strings.TrimSpace(
			r.RefreshToken,
		),
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r refreshTokenRequest) toDTO(
	request *stdhttp.Request,
) dto.RefreshTokenRequest {
	return dto.RefreshTokenRequest{
		RefreshToken: strings.TrimSpace(
			r.RefreshToken,
		),

		UserAgent: request.UserAgent(),

		IPAddress: clientIPAddress(
			request,
		),
	}
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`

	NewPassword string `json:"newPassword"`
}

func (r changePasswordRequest) toDTO(
	userID uuid.UUID,
) (
	dto.ChangePasswordRequest,
	error,
) {
	if userID == uuid.Nil {
		return dto.ChangePasswordRequest{},
			errors.New(
				"user ID tidak valid",
			)
	}

	password, err := valueobject.NewPassword(
		r.NewPassword,
	)
	if err != nil {
		return dto.ChangePasswordRequest{},
			err
	}

	return dto.ChangePasswordRequest{
		UserID: userID,

		CurrentPassword: r.CurrentPassword,

		NewPassword: password,
	}, nil
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (r forgotPasswordRequest) toDTO() (
	dto.ForgotPasswordRequest,
	error,
) {
	email, err := valueobject.NewEmail(
		r.Email,
	)
	if err != nil {
		return dto.ForgotPasswordRequest{},
			err
	}

	return dto.ForgotPasswordRequest{
		Email: email,
	}, nil
}

type resetPasswordRequest struct {
	Token string `json:"token"`

	NewPassword string `json:"newPassword"`
}

func (r resetPasswordRequest) toDTO() (
	dto.ResetPasswordRequest,
	error,
) {
	password, err := valueobject.NewPassword(
		r.NewPassword,
	)
	if err != nil {
		return dto.ResetPasswordRequest{},
			err
	}

	return dto.ResetPasswordRequest{
		Token: strings.TrimSpace(
			r.Token,
		),

		NewPassword: password,
	}, nil
}

type revokeAllSessionsRequest struct {
	KeepCurrent bool `json:"keepCurrent"`
}

func decodeJSON(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
	target any,
) error {
	if request == nil {
		return errors.New(
			"request tidak tersedia",
		)
	}

	if request.Body == nil {
		return errors.New(
			"request body wajib diisi",
		)
	}

	request.Body = stdhttp.MaxBytesReader(
		writer,
		request.Body,
		maxRequestBodyBytes,
	)

	decoder := json.NewDecoder(
		request.Body,
	)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(
		target,
	); err != nil {
		return err
	}

	var trailing any

	err := decoder.Decode(
		&trailing,
	)

	if !errors.Is(
		err,
		io.EOF,
	) {
		if err == nil {
			return errors.New(
				"request body hanya boleh berisi satu JSON object",
			)
		}

		return err
	}

	return nil
}

func clientIPAddress(
	request *stdhttp.Request,
) string {
	if request == nil {
		return ""
	}

	forwardedFor := strings.TrimSpace(
		request.Header.Get(
			"X-Forwarded-For",
		),
	)

	if forwardedFor != "" {
		parts := strings.Split(
			forwardedFor,
			",",
		)

		if len(parts) > 0 {
			return strings.TrimSpace(
				parts[0],
			)
		}
	}

	realIP := strings.TrimSpace(
		request.Header.Get(
			"X-Real-IP",
		),
	)

	if realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(
		strings.TrimSpace(
			request.RemoteAddr,
		),
	)
	if err == nil {
		return host
	}

	return strings.TrimSpace(
		request.RemoteAddr,
	)
}
