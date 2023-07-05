package main

import (
	"context"
	"os"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cprates/gotest/internal/helpers"
	testing "github.com/cprates/gotest/runner"
)

type SQLiteTests struct {
	*testing.Runner
}

func (s *SQLiteTests) TestSQLiteQueryAndInsert(t *testing.T) {
	tableName := "testsdb"
	db := helpers.NewSQLite(t, tableName)
	defer db.Close()
	defer os.Remove(tableName + ".db")

	ctx := context.Background()
	sql := `CREATE TABLE tests (
		id bigint,
		created_at TIMESTAMP NOT NULL
	);`
	_, err := helpers.Exec(ctx, db, sql)
	require.NoError(t, err)

	id := time.Now().UnixNano()
	n, err := helpers.Exec(ctx, db, `INSERT INTO tests (id, created_at) VALUES ($1, $2)`, id, time.Now().Format(time.RFC3339))
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	rows, err := helpers.Query[*TableModel](ctx, db, "SELECT * FROM tests WHERE id = $1", id)
	require.NoError(t, err)
	require.Equal(t, 1, len(rows))
}
