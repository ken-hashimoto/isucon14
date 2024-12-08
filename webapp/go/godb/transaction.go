package godb

import (
	"context"
	"fmt"

	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey struct {
	Name string
}

//nolint:gochecknoglobals
var transactionCtxKey = &contextKey{Name: "TransactionDB"}

//go:generate mockgen -destination=../db/mocks/mock_transaction.go -package=mocks github.com/hrbrain/hrbrain/apps/tama/app/internal/repositories/db Transaction
type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transaction struct {
	db *pgxpool.Pool
}

func NewTransaction(db DB) Transaction {
	return &transaction{
		db: db.Base(),
	}
}

func (t *transaction) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, hasTx := GetTx(ctx); hasTx {
		return fn(ctx)
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	log.Printf("TRANSACTION: BEGIN")

	ctx = SaveTx(ctx, tx)

	err = fn(ctx)
	if err != nil {
		err2 := tx.Rollback(ctx)
		log.Printf("TRANSACTION: ROLLBACK")
		return fmt.Errorf("(rollback err=%v): %w", err2, err)
	}

	if err := tx.Commit(ctx); err != nil {
		err2 := tx.Rollback(ctx)
		log.Printf("TRANSACTION: ROLLBACK")
		return fmt.Errorf("failed transaction commit (rollback err=%v): %w", err2, err)
	}

	log.Printf("TRANSACTION: COMMIT")

	return nil
}

func SaveTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, transactionCtxKey, tx)
}

func GetTx(ctx context.Context) (pgx.Tx, bool) {
	// contextからsessionを取得
	v := ctx.Value(&transactionCtxKey)
	if v == nil {
		return nil, false
	}

	session, ok := v.(pgx.Tx)

	return session, ok
}
