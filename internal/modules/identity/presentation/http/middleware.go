//Package http Files: internal/modules/identity/presentation/http/middleware.go

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ruangwali/internal/modules/identity/application/ports"
	"github.com/ruangwali/internal/shared/application/requestcontext"
)

const (
	authorizationHeader = "Authorization"
	bearerScheme        = "Bearer"
)

var (
	errMissingAuthorizationHeader = errors.New(
		"authorization header tidak tersedia",
	)

	errInvalidAuthorizationHeader = errors.New(
		"authorization header tidak valid",
	)

	errInvalidBearerScheme = errors.New(
		"authorization scheme harus Bearer",
	)
)

type AuthMiddleware struct {
	accessTokens ports.AccessTokenService
}

func NewAuthMiddleware(
	accessTokens ports.AccessTokenService,
) *AuthMiddleware {
	if accessTokens == nil {
		panic(
			"auth middleware: access token service nil",
		)
	}

	return &AuthMiddleware{
		accessTokens: accessTokens,
	}
}

func (m *AuthMiddleware) Authenticate(
	next http.Handler,
) http.Handler {
	if next == nil {
		panic(
			"auth middleware: next handler nil",
		)
	}

	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			rawToken, err := extractBearerToken(
				request,
			)
			if err != nil {
				writeUnauthorized(
					writer,
					err,
				)

				return
			}

			claims, err := m.accessTokens.Parse(
				request.Context(),
				rawToken,
			)
			if err != nil {
				writeUnauthorized(
					writer,
					errors.New(
						"access token tidak valid",
					),
				)

				return
			}

			ctx := requestcontext.WithUserID(
				request.Context(),
				claims.UserID,
			)

			next.ServeHTTP(
				writer,
				request.WithContext(ctx),
			)
		},
	)
}

func extractBearerToken(
	request *http.Request,
) (string, error) {
	if request == nil {
		return "",
			errMissingAuthorizationHeader
	}

	header := strings.TrimSpace(
		request.Header.Get(
			authorizationHeader,
		),
	)

	if header == "" {
		return "",
			errMissingAuthorizationHeader
	}

	parts := strings.Fields(
		header,
	)

	if len(parts) != 2 {
		return "",
			errInvalidAuthorizationHeader
	}

	if !strings.EqualFold(
		parts[0],
		bearerScheme,
	) {
		return "",
			errInvalidBearerScheme
	}

	rawToken := strings.TrimSpace(
		parts[1],
	)

	if rawToken == "" {
		return "",
			errInvalidAuthorizationHeader
	}

	return rawToken, nil
}

func writeUnauthorized(
	writer http.ResponseWriter,
	err error,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	writer.Header().Set(
		"Cache-Control",
		"no-store",
	)

	writer.WriteHeader(
		http.StatusUnauthorized,
	)

	response := unauthorizedResponse{
		Error: unauthorizedError{
			Code: "UNAUTHORIZED",

			Message: err.Error(),
		},
	}

	_ = json.NewEncoder(
		writer,
	).Encode(response)
}

type unauthorizedResponse struct {
	Error unauthorizedError `json:"error"`
}

type unauthorizedError struct {
	Code string `json:"code"`

	Message string `json:"message"`
}
