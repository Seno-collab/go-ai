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
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				fmt.Println("✅ Connected to Postgres!")
				cancel()
				return pool, nil
			}
			fmt.Println("Ping error: %v\n", pingErr)
			pool.Close()
		}
		cancel()
		fmt.Printf("❌ Attempt %d/%d failed: %v — retrying in %v…\n", attempt, maxRetries, err, backoff)

		time.Sleep(backoff)
		backoff *= 2

	}

	return nil, fmt.Errorf("cannot connect to postgres after %d attempts: %v", maxRetries, err)
}
