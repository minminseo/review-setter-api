package transaction

import (
	"context"
)

type ITransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
