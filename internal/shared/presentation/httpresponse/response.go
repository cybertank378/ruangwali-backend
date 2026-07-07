package httpresponse

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	shareddomain "github.com/ruangwali/internal/shared/domain"
)

const contentTypeJSON = "application/json; charset=utf-8"

type SuccessResponse struct {
	Success bool `json:"success"`

	Data any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool `json:"success"`

	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	writeJSON(
		w,
		status,
		SuccessResponse{
			Success: true,
			Data:    data,
		},
	)
}

func Error(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    code,
				Message: message,
			},
		},
	)
}

func FromError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(
		err,
		shareddomain.ErrValidation,
	):
		Error(
			w,
			http.StatusUnprocessableEntity,
			"VALIDATION_ERROR",
			"data tidak valid",
		)

	case errors.Is(
		err,
		shareddomain.ErrUnauthorized,
	):
		Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"autentikasi diperlukan",
		)

	case errors.Is(
		err,
		shareddomain.ErrForbidden,
	):
		Error(
			w,
			http.StatusForbidden,
			"FORBIDDEN",
			"akses ditolak",
		)

	case errors.Is(
		err,
		shareddomain.ErrNotFound,
	):
		Error(
			w,
			http.StatusNotFound,
			"NOT_FOUND",
			"resource tidak ditemukan",
		)

	case errors.Is(
		err,
		shareddomain.ErrConflict,
	):
		Error(
			w,
			http.StatusConflict,
			"CONFLICT",
			"resource mengalami konflik",
		)

	default:
		slog.Error(
			"unhandled application error",
			"error",
			err,
		)

		Error(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"terjadi kesalahan internal",
		)
	}
}

func NoContent(
	w http.ResponseWriter,
) {
	w.WriteHeader(
		http.StatusNoContent,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set(
		"Content-Type",
		contentTypeJSON,
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error(
			"failed to encode http response",
			"error",
			err,
			"status",
			status,
		)
	}
}
