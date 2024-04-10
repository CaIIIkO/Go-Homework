package tests

import "homework-3/tests/postgresql"

var (
	db *postgresql.TDB
)

func init() {
	db = postgresql.NewFromEnv()
}
