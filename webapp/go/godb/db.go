package godb

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/isucon/isucon14/webapp/go/sqlcgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Queries(ctx context.Context) *sqlcgen.Queries
	Close(ctx context.Context)
	Base() *pgxpool.Pool
	Conn(ctx context.Context) sqlcgen.DBTX
}

type database struct {
	DB *pgxpool.Pool
}

func NewDB(dbConfig *mysql.Config) (DB, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}
	return &database{DB: dbpool}, nil
}

func (db database) Queries(ctx context.Context) *sqlcgen.Queries {
	q := sqlcgen.New(db.DB)
	if tx, ok := GetTx(ctx); ok {
		return q.WithTx(tx)
	}
	return q
}

func (db database) Close(_ context.Context) {
	db.DB.Close()
}

func (db database) Base() *pgxpool.Pool {
	return db.DB
}

func (db database) Conn(ctx context.Context) sqlcgen.DBTX {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return db.DB
}
