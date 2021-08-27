package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Connect(DSN string) *sqlx.DB {
	db, err := sqlx.Connect("pgx", DSN)
	if err != nil {
		log.Panic().
			AnErr("error", err).
			Msg("failed to connect to db")
	}
	return db
}
