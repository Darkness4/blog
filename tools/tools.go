//go:build tools

package tools

import (
	_ "github.com/jackc/pgx/v5"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
