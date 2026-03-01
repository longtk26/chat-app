package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/longtk26/chat-app/configs"
)

func NewPostgresPool(cfg configs.DBConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	fmt.Printf("Connecting to PostgreSQL with DSN: %s\n", cfg.DSN())
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	fmt.Println("PostgreSQL connection pool created")
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return pool, nil
}
