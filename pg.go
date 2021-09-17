package pgxx

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

var TxReadOnly = pgx.TxOptions{AccessMode: pgx.ReadOnly}

// Rollback rolls back a postgres transaction without a context timeout.
// If there is an error it is logged using the standard library.
// It is a convinience function.
func Rollback(tx pgx.Tx) {
	// ctx,cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	err := tx.Rollback(context.Background())
	if err != nil && err != pgx.ErrTxClosed {
		log.Printf("pgxx.Rollback: ")
		// log.LoggerWithoutCaller.Debug().Caller(1).Err(err).Msg("pg.Rollback")
	}
}
