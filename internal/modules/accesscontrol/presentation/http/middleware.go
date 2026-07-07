// Package http FilePath: /internal/modules/accesscontrol/presentation/http
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	accessservice "github.com/ruangwali/internal/modules/accesscontrol/application/service"
	accessdomain "github.com/ruangwali/internal/modules/accesscontrol/domain"
)

type UserIDResolver interface {
	ResolveUserID(
		ctx context.Context,
	) (uuid.UUID, bool)
}

type Middleware struct {
	authorizationService *accessservice.AuthorizationService
	userIDResolver       UserIDResolver
}

func NewMiddleware(
	authorizationService *accessservice.AuthorizationService,
	userIDResolver UserIDResolver,
) *Middleware {
	if authorizationService == nil {
		panic(
			"access control middleware: authorization service nil",
		)
	}

	if userIDResolver == nil {
		panic(
			"access control middleware: user id resolver nil",
		)
	}

	return &Middleware{
		authorizationService: authorizationService,
		userIDResolver:       userIDResolver,
	}
}

func (m *Middleware) RequirePermission(
	permission string,
) func(http.Handler) http.Handler {
	return func(
		next http.Handler,
	) http.Handler {
		if next == nil {
			panic(
				"access control middleware: next handler nil",
			)
		}

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				userID, ok := m.userIDResolver.ResolveUserID(
					r.Context(),
				)
				if !ok || userID == uuid.Nil {
					writeError(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"autentikasi diperlukan",
					)

					return
				}

				err := m.authorizationService.Authorize(
					r.Context(),
					userID,
					permission,
				)
				if err != nil {
					m.writeAuthorizationError(
						w,
						err,
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}

func (m *Middleware) RequireAnyPermission(
	permissions ...string,
) func(http.Handler) http.Handler {
	return func(
		next http.Handler,
	) http.Handler {
		if next == nil {
			panic(
				"access control middleware: next handler nil",
			)
		}

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				userID, ok := m.userIDResolver.ResolveUserID(
					r.Context(),
				)
				if !ok || userID == uuid.Nil {
					writeError(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"autentikasi diperlukan",
					)

					return
				}

				err := m.authorizationService.AuthorizeAny(
					r.Context(),
					userID,
					permissions...,
				)
				if err != nil {
					m.writeAuthorizationError(
						w,
						err,
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}

func (m *Middleware) RequireAllPermissions(
	permissions ...string,
) func(http.Handler) http.Handler {
	return func(
		next http.Handler,
	) http.Handler {
		if next == nil {
			panic(
				"access control middleware: next handler nil",
			)
		}

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				userID, ok := m.userIDResolver.ResolveUserID(
					r.Context(),
				)
				if !ok || userID == uuid.Nil {
					writeError(
						w,
						http.StatusUnauthorized,
						"unauthorized",
						"autentikasi diperlukan",
					)

					return
				}

				err := m.authorizationService.AuthorizeAll(
					r.Context(),
					userID,
					permissions...,
				)
				if err != nil {
					m.writeAuthorizationError(
						w,
						err,
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}

func (m *Middleware) writeAuthorizationError(
	w http.ResponseWriter,
	err error,
) {
	if errors.Is(
		err,
		accessdomain.ErrAccessDenied,
	) {
		writeError(
			w,
			http.StatusForbidden,
			"forbidden",
			"akses ditolak",
		)

		return
	}

	writeError(
		w,
		http.StatusInternalServerError,
		"internal_error",
		"terjadi kesalahan internal",
	)
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		status,
	)

	_ = json.NewEncoder(
		w,
	).Encode(
		errorResponse{
			Error: errorDetail{
				Code:    code,
				Message: message,
			},
		},
	)
}
