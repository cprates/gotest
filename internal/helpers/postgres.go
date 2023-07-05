package helpers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	testing "github.com/cprates/gotest/runner"
)

func NewPostgres(t *testing.T, user, password, table string) *sqlx.DB {
	ds := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, table)
	db, err := sqlx.Connect("postgres", ds)
	require.NoError(t, err)

	return db
}
