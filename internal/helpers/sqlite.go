package helpers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	testing "github.com/cprates/gotest/runner"
)

func NewSQLite(t *testing.T, table string) *sqlx.DB {
	ds := fmt.Sprintf("file:%s.db", table)
	db, err := sqlx.Connect("sqlite3", ds)
	require.NoError(t, err)

	return db
}
