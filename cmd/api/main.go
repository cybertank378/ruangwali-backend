package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/ruangwali/internal/composition"
	"github.com/ruangwali/internal/platform/buildinfo"
	"github.com/ruangwali/internal/platform/config"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(
			"gagal memuat konfigurasi: ",
			err,
		)
	}

	app, err := composition.Build(
		context.Background(),
		cfg,
	)
	if err != nil {
		log.Fatal(
			"gagal membangun aplikasi: ",
			err,
		)
	}

	runErr := make(
		chan error,
		1,
	)

	go func() {
		log.Printf(
			"RuangWali API listening on %s version=%s",
			cfg.HTTP.Addr,
			buildinfo.Version,
		)

		runErr <- app.Run()
	}()

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-signalCtx.Done():
		log.Print(
			"shutdown signal received",
		)

	case err := <-runErr:
		if err != nil {
			log.Printf(
				"application stopped unexpectedly: %v",
				err,
			)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := app.Shutdown(
		shutdownCtx,
	); err != nil {
		if errors.Is(
			err,
			context.DeadlineExceeded,
		) {
			log.Printf(
				"graceful shutdown timeout: %v",
				err,
			)

			return
		}

		log.Printf(
			"gagal melakukan graceful shutdown: %v",
			err,
		)

		return
	}

	log.Print(
		"application stopped gracefully",
	)
}
