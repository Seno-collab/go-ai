package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(dsn string) (*pgxpool.Pool, error) {

	var pool *pgxpool.Pool
	var err error

	maxRetries := 8
	backoff := time.Second // start with 1s
	for attempt := 1; attempt < maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				fmt.Println("✅ Connected to Postgres!")
				cancel()
				return pool, nil
			}
		}
		cancel()
		fmt.Printf("❌ Failed: %v — retrying in %v…\n", err, backoff)

		time.Sleep(backoff)
		backoff *= 2

	}

	return nil, fmt.Errorf("cannot connect to postgres after %d attempts: %v", maxRetries, err)
}
