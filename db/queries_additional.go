package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func (q *Queries) CreateOrIncrementPageViewsOnUniqueIP(
	ctx context.Context,
	db *pgxpool.Pool,
	pageID string,
	ip string,
) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Err(err).Msg("failed to rollback transaction")
		}
	}()
	qtx := q.WithTx(tx)

	if _, err := qtx.CreateOrIncrementPageViews(ctx, pageID); err != nil {
		return err
	}

	pvi, err := qtx.CreatePageViewsIPs(ctx, CreatePageViewsIPsParams{
		PageID: pageID,
		Ip:     ip,
	})
	if err != nil {
		return err
	}

	if len(pvi) == 0 { // IP already exists
		return nil
	}

	return tx.Commit(ctx)
}

func (q *Queries) FindPageViewsOrZero(ctx context.Context, pageID string) (PageView, error) {
	ret, err := q.FindPageViews(ctx, pageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PageView{PageID: pageID, Views: 0}, nil
		}
		return PageView{}, err
	}
	return ret, nil
}
