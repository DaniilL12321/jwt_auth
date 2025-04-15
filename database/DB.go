package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDBconnection(ctx context.Context) (dbpool *pgxpool.Pool, err error) {
	url := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return
	}
	cfg.MinConns = int32(1)
	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	return
}
