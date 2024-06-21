package db

import (
	"database/sql"
	"embed"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

var sl = &stdlogger{zl: &log.Logger}

type stdlogger struct {
	zl *zerolog.Logger
}

func (l *stdlogger) Fatalf(format string, v ...interface{}) {
	l.zl.Fatal().Msgf(format, v...)
}

func (l *stdlogger) Printf(format string, v ...interface{}) {
	format = strings.Trim(format, "\n")
	l.zl.Printf(format, v...)
}

func init() {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Panic().Err(err).Msg("failed to set goose dialect")
	}

	goose.SetLogger(sl)
}

func InitialMigration(db *sql.DB) error {
	return goose.Up(db, "migrations")
}
