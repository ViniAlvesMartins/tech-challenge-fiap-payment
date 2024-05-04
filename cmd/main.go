package main

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/infra"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/external/database/postgres"
	"log/slog"
	"os"
)

func main() {
	var err error
	var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

	cfg, err := infra.NewConfig()
	if err != nil {
		logger.Error("error loading config", err)
		panic(err)
	}

	db, err := postgres.NewConnection(cfg)
	if err != nil {
		logger.Error("error connecting tdo database", err)
		panic(err)
	}

	migration := postgres.NewMigration(db, cfg, logger)
	migration.CreateSchema()
	migration.Migrate()
}
