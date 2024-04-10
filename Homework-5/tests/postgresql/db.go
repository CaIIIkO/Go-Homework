package postgresql

import (
	"context"
	"fmt"
	"homework-3/internal/pkg/db"
	"testing"
)

type TDB struct {
	DB db.DBops
}

func NewFromEnv() *TDB {
	db, err := db.NewDb(context.Background())
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, tableName string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName)
}

func (d *TDB) TearDown() {
}

func (d *TDB) truncateTable(ctx context.Context, tableName string) {
	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY", tableName)
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
