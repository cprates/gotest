package main

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cprates/gotest/internal/helpers"
	testing "github.com/cprates/gotest/runner"
)

type TableModel struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type PostgresTests struct {
	*testing.Runner
}

func (g *PostgresTests) TestPostgresQueryAndInsert(t *testing.T) {
	db := helpers.NewPostgres(t, "postgres", "postgres", "testsdb")
	defer db.Close()

	ctx := context.Background()
	id := time.Now().UnixNano()
	n, err := helpers.Exec(ctx, db, `INSERT INTO tests (id, created_at) VALUES ($1, $2)`, id, time.Now().Format(time.RFC3339))
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	rows, err := helpers.Query[*TableModel](ctx, db, "SELECT * FROM tests WHERE id = $1", id)
	require.NoError(t, err)
	require.Equal(t, 1, len(rows))
}

func (g *PostgresTests) TestPostgresQueryAndInsertAsSubtests(t *testing.T) {
	db := helpers.NewPostgres(t, "postgres", "postgres", "testsdb")
	defer db.Close()
	ctx := context.Background()
	id := time.Now().UnixNano()

	testsSet := []struct {
		description string
		action      func(t *testing.T)
	}{
		{
			description: "insert row",
			action: func(t *testing.T) {
				n, err := helpers.Exec(ctx, db, `INSERT INTO tests (id, created_at) VALUES ($1, $2)`, id, time.Now().Format(time.RFC3339))
				require.NoError(t, err)
				require.Equal(t, int64(1), n)
			},
		},
		{
			description: "query rows",
			action: func(t *testing.T) {
				rows, err := helpers.Query[*TableModel](ctx, db, "SELECT * FROM tests WHERE id = $1", id)
				require.NoError(t, err)
				require.Equal(t, 1, len(rows))
			},
		},
	}

	for _, test := range testsSet {
		t.Run(test.description, func(t *testing.T) {
			test.action(t)
		})
	}
}
