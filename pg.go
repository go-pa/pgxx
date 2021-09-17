package pgxx

import (
	"context"
	"fmt"

	"github.com/go-pa/pgxx/logging"
	"github.com/jackc/pgx/v4"
)

var TxReadOnly = pgx.TxOptions{AccessMode: pgx.ReadOnly}

// Rollback rolls back a postgres transaction without a context timeout.
// If there is an error it is logged using the standard library.
// It is a convinience function.
func Rollback(ctx context.Context, tx pgx.Tx) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := tx.Rollback(ctx)
	if err != nil && err != pgx.ErrTxClosed {
		if logger := logging.Get(ctx); logger != nil {
			fmt.Fprint(logger, "pgxx.Rollback error: ", err)
		}
	}
}
