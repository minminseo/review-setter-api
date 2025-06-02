// infrastructure/repository/transaction_manager.go
package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"

	transactionUsecase "github.com/minminseo/recall-setter/usecase/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minminseo/recall-setter/infrastructure/db"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) transactionUsecase.ITransactionManager {
	return &TransactionManager{pool: pool}
}

func (tm *TransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 渡されたプールを使ってトランザクション開始
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// トランザクション用のQueriesを生成
	txQ := dbgen.New(tx)

	// トランザクション用のQueriesをcontextに詰め込むメソッドを実行。これでctxがトランザクション内か外かを判別できるようになる。
	ctxWithTx := db.WithQueries(ctx, txQ)

	err = fn(ctxWithTx)
	if err != nil {
		log.Printf("DBロールバック: %v\n", err)
		return err
	}

	if err := tx.Commit(ctxWithTx); err != nil {
		return fmt.Errorf("commit失敗: %w", err)
	}
	return nil
}
