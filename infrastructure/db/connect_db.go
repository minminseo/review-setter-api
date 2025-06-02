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

var (
	globalQueries *dbgen.Queries
)

func NewDB(ctx context.Context) (*pgxpool.Pool, error) {

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
		return nil, err
	}

	queries := dbgen.New(pool)

	globalQueries = queries

	return pool, nil
}

// context.WithValue用のキーの型
type ctxKey string

const queriesKey ctxKey = "queries"

// contextにトランザクション用の*dbgen.Queriesがあるならそれを返す。
// なければグローバルに初期化された*dbgen.Queriesを返す
func GetQuery(ctx context.Context) *dbgen.Queries {
	txq := getQueriesWithContext(ctx)
	if txq != nil {
		return txq
	}
	return globalQueries
}

func getQueriesWithContext(ctx context.Context) *dbgen.Queries {
	queries, ok := ctx.Value(queriesKey).(*dbgen.Queries)
	if !ok {
		return nil
	}
	return queries
}

// ctxにトランザクション用の*dbgen.Queriesをセットして返す。
func WithQueries(ctx context.Context, q *dbgen.Queries) context.Context {
	return context.WithValue(ctx, queriesKey, q)
}
