package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const migrationSource = "file://migrations"

func main() {
	loadEnvironment()

	if len(os.Args) < 2 {
		log.Fatal(
			"perintah migration wajib diisi: up, down, version, atau force",
		)
	}

	databaseURL := strings.TrimSpace(
		os.Getenv("DATABASE_URL"),
	)

	if databaseURL == "" {
		log.Fatal(
			"DATABASE_URL wajib diisi",
		)
	}

	migrator, err := migrate.New(
		migrationSource,
		databaseURL,
	)
	if err != nil {
		log.Fatalf(
			"gagal membuat migrator: %v",
			err,
		)
	}

	defer closeMigrator(migrator)

	command := strings.ToLower(
		strings.TrimSpace(os.Args[1]),
	)

	switch command {
	case "up":
		runUp(migrator)

	case "down":
		runDown(migrator)

	case "version":
		showVersion(migrator)

	case "force":
		runForce(migrator)

	default:
		log.Fatalf(
			"perintah migration tidak dikenal: %s",
			command,
		)
	}
}

func loadEnvironment() {
	err := godotenv.Load()

	if err == nil {
		return
	}

	if errors.Is(err, os.ErrNotExist) {
		return
	}

	log.Fatalf(
		"gagal memuat file .env: %v",
		err,
	)
}

func closeMigrator(
	migrator *migrate.Migrate,
) {
	sourceErr, databaseErr := migrator.Close()

	if sourceErr != nil {
		log.Printf(
			"gagal menutup migration source: %v",
			sourceErr,
		)
	}

	if databaseErr != nil {
		log.Printf(
			"gagal menutup migration database: %v",
			databaseErr,
		)
	}
}

func runUp(
	migrator *migrate.Migrate,
) {
	err := migrator.Up()

	if err == nil {
		log.Println(
			"migration berhasil dijalankan",
		)
		return
	}

	if errors.Is(
		err,
		migrate.ErrNoChange,
	) {
		log.Println(
			"database sudah menggunakan migration terbaru",
		)
		return
	}

	log.Fatalf(
		"migration gagal: %v",
		err,
	)
}

func runDown(
	migrator *migrate.Migrate,
) {
	err := migrator.Steps(-1)

	if err == nil {
		log.Println(
			"rollback 1 migration berhasil",
		)
		return
	}

	if errors.Is(
		err,
		migrate.ErrNoChange,
	) {
		log.Println(
			"tidak ada migration yang dapat di-rollback",
		)
		return
	}

	log.Fatalf(
		"rollback migration gagal: %v",
		err,
	)
}

func showVersion(
	migrator *migrate.Migrate,
) {
	version, dirty, err := migrator.Version()
	if err != nil {
		if errors.Is(
			err,
			migrate.ErrNilVersion,
		) {
			log.Println(
				"database belum memiliki migration version",
			)
			return
		}

		log.Fatalf(
			"gagal membaca migration version: %v",
			err,
		)
	}

	fmt.Printf(
		"version=%d dirty=%t\n",
		version,
		dirty,
	)
}

func runForce(
	migrator *migrate.Migrate,
) {
	if len(os.Args) < 3 {
		log.Fatal(
			"version wajib diisi untuk perintah force",
		)
	}

	version, err := strconv.Atoi(
		strings.TrimSpace(os.Args[2]),
	)
	if err != nil {
		log.Fatalf(
			"migration version tidak valid: %v",
			err,
		)
	}

	if version < 0 {
		log.Fatal(
			"migration version tidak boleh negatif",
		)
	}

	if err := migrator.Force(version); err != nil {
		log.Fatalf(
			"gagal force migration version: %v",
			err,
		)
	}

	log.Printf(
		"migration version berhasil dipaksa ke %d",
		version,
	)
}
