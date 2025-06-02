package dbgen

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Queries型にCopyFromメソッドをラッピング
func (q *Queries) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	return q.db.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
