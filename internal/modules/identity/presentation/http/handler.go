// =========================================================
// File: internal/modules/identity/presentation/http/handler.go
// =========================================================

package http

import (
	"errors"
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ruangwali/internal/modules/identity/application/usecase"
	"github.com/ruangwali/internal/shared/application/requestcontext"
	"github.com/ruangwali/internal/shared/presentation/httpresponse"
)

type Handler struct {
	login             *usecase.LoginUseCase
	logout            *usecase.LogoutUseCase
	refreshToken      *usecase.RefreshTokenUseCase
	getCurrentUser    *usecase.GetCurrentUserUseCase
	changePassword    *usecase.ChangePasswordUseCase
	forgotPassword    *usecase.ForgotPasswordUseCase
	resetPassword     *usecase.ResetPasswordUseCase
	listSessions      *usecase.ListSessionsUseCase
	revokeSession     *usecase.RevokeSessionUseCase
	revokeAllSessions *usecase.RevokeAllSessionsUseCase
}

func NewHandler(
	login *usecase.LoginUseCase,
	logout *usecase.LogoutUseCase,
	refreshToken *usecase.RefreshTokenUseCase,
	getCurrentUser *usecase.GetCurrentUserUseCase,
	changePassword *usecase.ChangePasswordUseCase,
	forgotPassword *usecase.ForgotPasswordUseCase,
	resetPassword *usecase.ResetPasswordUseCase,
	listSessions *usecase.ListSessionsUseCase,
	revokeSession *usecase.RevokeSessionUseCase,
	revokeAllSessions *usecase.RevokeAllSessionsUseCase,
) *Handler {
	if login == nil {
		panic(
			"identity handler: login use case nil",
		)
	}

	if logout == nil {
		panic(
			"identity handler: logout use case nil",
		)
	}

	if refreshToken == nil {
		panic(
			"identity handler: refresh token use case nil",
		)
	}

	if getCurrentUser == nil {
		panic(
			"identity handler: get current user use case nil",
		)
	}

	if changePassword == nil {
		panic(
			"identity handler: change password use case nil",
		)
	}

	if forgotPassword == nil {
		panic(
			"identity handler: forgot password use case nil",
		)
	}

	if resetPassword == nil {
		panic(
			"identity handler: reset password use case nil",
		)
	}

	if listSessions == nil {
		panic(
			"identity handler: list sessions use case nil",
		)
	}

	if revokeSession == nil {
		panic(
			"identity handler: revoke session use case nil",
		)
	}

	if revokeAllSessions == nil {
		panic(
			"identity handler: revoke all sessions use case nil",
		)
	}

	return &Handler{
		login:             login,
		logout:            logout,
		refreshToken:      refreshToken,
		getCurrentUser:    getCurrentUser,
		changePassword:    changePassword,
		forgotPassword:    forgotPassword,
		resetPassword:     resetPassword,
		listSessions:      listSessions,
		revokeSession:     revokeSession,
		revokeAllSessions: revokeAllSessions,
	}
}

func (h *Handler) Login(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	var body loginRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	response, err := h.login.Execute(
		request.Context(),
		body.toDTO(request),
	)
	if err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.JSON(
		writer,
		stdhttp.StatusOK,
		response,
	)
}

func (h *Handler) Logout(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	var body logoutRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	if strings.TrimSpace(
		body.RefreshToken,
	) == "" {
		writeInvalidRequest(
			writer,
			errors.New(
				"refreshToken wajib diisi",
			),
		)

		return
	}

	if err := h.logout.Execute(
		request.Context(),
		body.toDTO(),
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.NoContent(
		writer,
	)
}

func (h *Handler) RefreshToken(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	var body refreshTokenRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	if strings.TrimSpace(
		body.RefreshToken,
	) == "" {
		writeInvalidRequest(
			writer,
			errors.New(
				"refreshToken wajib diisi",
			),
		)

		return
	}

	response, err := h.refreshToken.Execute(
		request.Context(),
		body.toDTO(request),
	)
	if err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.JSON(
		writer,
		stdhttp.StatusOK,
		response,
	)
}

func (h *Handler) GetCurrentUser(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	userID, ok := requestcontext.UserID(
		request.Context(),
	)
	if !ok {
		writeAuthenticationRequired(
			writer,
		)

		return
	}

	response, err := h.getCurrentUser.Execute(
		request.Context(),
		userID,
	)
	if err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.JSON(
		writer,
		stdhttp.StatusOK,
		response,
	)
}

func (h *Handler) ChangePassword(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	userID, ok := requestcontext.UserID(
		request.Context(),
	)
	if !ok {
		writeAuthenticationRequired(
			writer,
		)

		return
	}

	var body changePasswordRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	input, err := body.toDTO(uuid.MustParse(userID.String()))

	if err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	if err := h.changePassword.Execute(
		request.Context(),
		input,
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.NoContent(
		writer,
	)
}

func (h *Handler) ForgotPassword(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	var body forgotPasswordRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	input, err := body.toDTO()
	if err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	if err := h.forgotPassword.Execute(
		request.Context(),
		input,
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	writeMessage(
		writer,
		stdhttp.StatusAccepted,
		"jika email terdaftar, instruksi reset password akan dikirim",
	)
}

func (h *Handler) ResetPassword(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	var body resetPasswordRequest

	if err := decodeJSON(
		writer,
		request,
		&body,
	); err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	input, err := body.toDTO()
	if err != nil {
		writeInvalidRequest(
			writer,
			err,
		)

		return
	}

	if err := h.resetPassword.Execute(
		request.Context(),
		input,
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.NoContent(
		writer,
	)
}

func (h *Handler) ListSessions(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	userID, ok := requestcontext.UserID(
		request.Context(),
	)
	if !ok {
		writeAuthenticationRequired(
			writer,
		)

		return
	}

	response, err := h.listSessions.Execute(
		request.Context(),
		userID,
		uuid.Nil,
	)
	if err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.JSON(
		writer,
		stdhttp.StatusOK,
		response,
	)
}

func (h *Handler) RevokeSession(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	userID, ok := requestcontext.UserID(
		request.Context(),
	)
	if !ok {
		writeAuthenticationRequired(
			writer,
		)

		return
	}

	sessionID, err := uuid.Parse(
		strings.TrimSpace(
			chi.URLParam(
				request,
				"sessionID",
			),
		),
	)
	if err != nil ||
		sessionID == uuid.Nil {
		writeInvalidRequest(
			writer,
			errors.New(
				"sessionID tidak valid",
			),
		)

		return
	}

	if err := h.revokeSession.Execute(
		request.Context(),
		userID,
		sessionID,
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.NoContent(
		writer,
	)
}

func (h *Handler) RevokeAllSessions(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
) {
	userID, ok := requestcontext.UserID(
		request.Context(),
	)
	if !ok {
		writeAuthenticationRequired(
			writer,
		)

		return
	}

	if err := h.revokeAllSessions.Execute(
		request.Context(),
		userID,
		nil,
	); err != nil {
		writeApplicationError(
			writer,
			err,
		)

		return
	}

	httpresponse.NoContent(
		writer,
	)
}
