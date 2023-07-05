package helpers

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// for queries
func Query[R any](ctx context.Context, db *sqlx.DB, query string, args ...any) ([]R, error) {
	res := []R{}
	err := db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// for inserts/updates
func Exec(ctx context.Context, db *sqlx.DB, query string, args ...any) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
