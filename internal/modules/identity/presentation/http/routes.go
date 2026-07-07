// =========================================================
// File: internal/modules/identity/presentation/http/routes.go
// =========================================================

package http

import (
	"github.com/go-chi/chi/v5"
)

func MountRoutes(
	router chi.Router,
	handler *Handler,
	authMiddleware *AuthMiddleware,
) {
	if router == nil {
		panic(
			"identity routes: router nil",
		)
	}

	if handler == nil {
		panic(
			"identity routes: handler nil",
		)
	}

	if authMiddleware == nil {
		panic(
			"identity routes: auth middleware nil",
		)
	}

	router.Route(
		"/api/v1/auth",
		func(
			authRouter chi.Router,
		) {
			mountPublicAuthRoutes(
				authRouter,
				handler,
			)

			mountProtectedAuthRoutes(
				authRouter,
				handler,
				authMiddleware,
			)
		},
	)
}

func mountPublicAuthRoutes(
	router chi.Router,
	handler *Handler,
) {
	router.Post(
		"/login",
		handler.Login,
	)

	router.Post(
		"/logout",
		handler.Logout,
	)

	router.Post(
		"/refresh",
		handler.RefreshToken,
	)

	router.Post(
		"/forgot-password",
		handler.ForgotPassword,
	)

	router.Post(
		"/reset-password",
		handler.ResetPassword,
	)
}

func mountProtectedAuthRoutes(
	router chi.Router,
	handler *Handler,
	authMiddleware *AuthMiddleware,
) {
	router.Group(
		func(
			protectedRouter chi.Router,
		) {
			protectedRouter.Use(
				authMiddleware.Authenticate,
			)

			protectedRouter.Get(
				"/me",
				handler.GetCurrentUser,
			)

			protectedRouter.Put(
				"/password",
				handler.ChangePassword,
			)

			protectedRouter.Get(
				"/sessions",
				handler.ListSessions,
			)

			protectedRouter.Delete(
				"/sessions",
				handler.RevokeAllSessions,
			)

			protectedRouter.Delete(
				"/sessions/{sessionID}",
				handler.RevokeSession,
			)
		},
	)
}
