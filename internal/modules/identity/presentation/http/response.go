// =========================================================
// File: internal/modules/identity/presentation/http/response.go.go
// =========================================================

package http

import (
	"errors"
	stdhttp "net/http"

	"github.com/ruangwali/internal/shared/presentation/httpresponse"
)

var (
	errInvalidRequestBody = errors.New(
		"request body tidak valid",
	)

	errAuthenticationRequired = errors.New(
		"autentikasi diperlukan",
	)
)

type messageResponse struct {
	Message string `json:"message"`
}

func writeInvalidRequest(
	writer stdhttp.ResponseWriter,
	err error,
) {
	message := errInvalidRequestBody.Error()

	if err != nil {
		message = err.Error()
	}

	httpresponse.Error(
		writer,
		stdhttp.StatusBadRequest,
		"INVALID_REQUEST",
		message,
	)
}

func writeAuthenticationRequired(
	writer stdhttp.ResponseWriter,
) {
	httpresponse.Error(
		writer,
		stdhttp.StatusUnauthorized,
		"UNAUTHORIZED",
		errAuthenticationRequired.Error(),
	)
}

func writeApplicationError(
	writer stdhttp.ResponseWriter,
	err error,
) {
	if err == nil {
		return
	}

	httpresponse.FromError(
		writer,
		err,
	)
}

func writeMessage(
	writer stdhttp.ResponseWriter,
	status int,
	message string,
) {
	httpresponse.JSON(
		writer,
		status,
		messageResponse{
			Message: message,
		},
	)
}
