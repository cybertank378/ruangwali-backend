package composition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruangwali/internal/platform/buildinfo"
	"github.com/ruangwali/internal/platform/config"
	"github.com/ruangwali/internal/platform/database"
)

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 15 * time.Second
	defaultWriteTimeout      = 30 * time.Second
	defaultIdleTimeout       = 60 * time.Second
)

type App struct {
	server *http.Server
	db     *pgxpool.Pool
}

type healthResponse struct {
	Status string `json:"status"`
}

type readinessResponse struct {
	Status   string         `json:"status"`
	Database string         `json:"database"`
	Build    buildinfo.Info `json:"build"`
}

func Build(
	ctx context.Context,
	cfg config.Config,
) (*App, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	db, err := database.OpenPostgres(
		ctx,
		cfg.Database,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membangun PostgreSQL dependency: %w",
			err,
		)
	}

	router := buildRouter(
		cfg,
		db,
	)

	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}

	return &App{
		server: server,
		db:     db,
	}, nil
}

func (a *App) Run() error {
	if a == nil {
		return errors.New(
			"application belum diinisialisasi",
		)
	}

	if a.server == nil {
		return errors.New(
			"HTTP server belum diinisialisasi",
		)
	}

	slog.Info(
		"HTTP server started",
		"addr",
		a.server.Addr,
		"version",
		buildinfo.Version,
	)

	err := a.server.ListenAndServe()
	if err != nil &&
		!errors.Is(
			err,
			http.ErrServerClosed,
		) {
		return fmt.Errorf(
			"HTTP server gagal berjalan: %w",
			err,
		)
	}

	return nil
}

func (a *App) Shutdown(
	ctx context.Context,
) error {
	if a == nil {
		return nil
	}

	var shutdownErr error

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf(
				"gagal menghentikan HTTP server: %w",
				err,
			)
		}
	}

	if a.db != nil {
		a.db.Close()
	}

	return shutdownErr
}

func buildRouter(
	cfg config.Config,
	db *pgxpool.Pool,
) http.Handler {
	router := chi.NewRouter()

	router.Get(
		"/health/live",
		handleLiveness,
	)

	router.Get(
		"/health/ready",
		handleReadiness(
			cfg.Database.HealthTimeout,
			db,
		),
	)

	return router
}

func handleLiveness(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		healthResponse{
			Status: "ok",
		},
	)
}

func handleReadiness(
	timeout time.Duration,
	db *pgxpool.Pool,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			timeout,
		)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			slog.Error(
				"readiness database check failed",
				"error",
				err,
			)

			writeJSON(
				w,
				http.StatusServiceUnavailable,
				readinessResponse{
					Status:   "not_ready",
					Database: "unavailable",
					Build:    buildinfo.Current(),
				},
			)

			return
		}

		writeJSON(
			w,
			http.StatusOK,
			readinessResponse{
				Status:   "ready",
				Database: "available",
				Build:    buildinfo.Current(),
			},
		)
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error(
			"gagal encode HTTP response",
			"error",
			err,
			"status",
			status,
		)
	}
}
