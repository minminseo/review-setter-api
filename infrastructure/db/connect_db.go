package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

func NewDB(ctx context.Context) (*pgxpool.Pool, *dbgen.Queries, error) {

	// 開発環境用
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	// クエリ用オブジェクトを生成
	queries := dbgen.New(pool)
	return pool, queries, nil
}
