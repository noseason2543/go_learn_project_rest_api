package databases

import (
	"go_learn_project_rest_api/config"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnection(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("Failed to connect db: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
