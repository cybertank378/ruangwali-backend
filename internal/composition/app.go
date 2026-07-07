package composition

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	acapp "github.com/ruangwali/backend/internal/modules/accesscontrol/application"
	acdomain "github.com/ruangwali/backend/internal/modules/accesscontrol/domain"
	acpg "github.com/ruangwali/backend/internal/modules/accesscontrol/infrastructure/postgres"
	"github.com/ruangwali/backend/internal/modules/identity/infrastructure/security"
	httpidentity "github.com/ruangwali/backend/internal/modules/identity/presentation/http"
	studentapp "github.com/ruangwali/backend/internal/modules/student/application"
	studentpg "github.com/ruangwali/backend/internal/modules/student/infrastructure/postgres"
	httpstudent "github.com/ruangwali/backend/internal/modules/student/presentation/http"
	"github.com/ruangwali/backend/internal/platform/config"
)

type App struct {
	Router http.Handler
	db     *pgxpool.Pool
}

func (a *App) Close() {
	a.db.Close()
}

func Build(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}

	permissionResolver := acpg.NewPermissionResolver(db)
	authorizer := acapp.NewAuthorizer(permissionResolver)
	tokens := security.NewTokenService(
		cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTSecret, cfg.AccessTTL,
	)

	studentRepository := studentpg.NewRepository(db)
	createStudent := studentapp.NewCreateStudentUseCase(studentRepository)
	studentHandler := httpstudent.NewHandler(createStudent)

	router := chi.NewRouter()
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"ruangwali-api"}`))
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(protected chi.Router) {
			protected.Use(httpidentity.Authenticate(tokens, authorizer))

			protected.With(
				httpidentity.RequirePermission(acdomain.StudentCreate),
			).Post("/students", studentHandler.Create)
		})
	})

	return &App{Router: router, db: db}, nil
}
