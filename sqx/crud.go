package sqx

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/dbscan"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

// QueryRowScan queries and scans single row responses.
func QueryRowScan(ctx context.Context, tx pgx.Tx, sb sq.Sqlizer, v ...interface{}) error {
	q, vs, err := sb.ToSql()
	if err != nil {
		return err
	}
	row := tx.QueryRow(ctx, q, vs...)
	if err := row.Scan(v...); err != nil {
		return err
	}
	return nil
}

func Query(ctx context.Context, tx pgx.Tx, sb sq.SelectBuilder, v interface{}) (pgx.Rows, error) {
	q, vs, err := sb.ToSql()
	if err != nil {
		return nil, err
	}
	return tx.Query(ctx, q, vs...)
}

func ScanAll(ctx context.Context, tx pgx.Tx, sb sq.SelectBuilder, v interface{}) error {
	q, vs, err := sb.ToSql()
	if err != nil {
		return err
	}
	rows, err := tx.Query(ctx, q, vs...)
	if err != nil {
		return err
	}
	err = pgxscan.ScanAll(v, rows)
	if err != nil {
		return err
	}
	return nil
}

func ScanOne(ctx context.Context, tx pgx.Tx, sb sq.SelectBuilder, v interface{}) error {
	q, vs, err := sb.ToSql()
	if err != nil {
		return err
	}
	rows, err := tx.Query(ctx, q, vs...)
	if err != nil {
		return err
	}
	err = pgxscan.ScanOne(v, rows)
	if err != nil {
		if dbscan.NotFound(err) {
			return pgx.ErrNoRows
		}
		return err
	}
	return nil
}

// Update does a sql Update and returns the number of rows updated.
func Update(ctx context.Context, tx pgx.Tx, sb sq.UpdateBuilder) (int64, error) {
	q, vs, err := sb.ToSql()
	if err != nil {
		return -1, err
	}
	tag, err := tx.Exec(ctx, q, vs...)
	if err != nil {
		return -1, err
	}
	return tag.RowsAffected(), nil
}

// UpdateOne does a sql update and checks if exactly one row if affected by the update.
//
// Returns ErrNoRowsAffected if 0 rows were affaced by the update.
func UpdateOne(ctx context.Context, tx pgx.Tx, sb sq.UpdateBuilder) error {
	q, vs, err := sb.ToSql()
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, q, vs...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrNoRowsAffected
	}
	return nil
}

// DeleteOne does a sql update and checks if exactly one row if affected by the update.
//
// Returns ErrNoRowsAffected if 0 rows were affaced by the update.
func DeleteOne(ctx context.Context, tx pgx.Tx, sb sq.DeleteBuilder) error {
	q, vs, err := sb.ToSql()
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, q, vs...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrNoRowsAffected
	}
	return nil
}
