// =========================================================
// File: internal/platform/openapi/routes.go
// =========================================================

package openapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(
	router chi.Router,
	handler *Handler,
) {
	if router == nil {
		panic(
			"openapi routes: router nil",
		)
	}

	if handler == nil {
		panic(
			"openapi routes: handler nil",
		)
	}

	router.Get(
		"/openapi.yaml",
		handler.Specification,
	)

	router.Get(
		"/docs",
		handler.Documentation,
	)

	router.Get(
		"/docs/",
		handler.Documentation,
	)

	router.Get(
		"/docs/index.html",
		handler.Documentation,
	)

	router.Method(
		http.MethodGet,
		"/swagger",
		http.RedirectHandler(
			"/docs",
			http.StatusTemporaryRedirect,
		),
	)
}
