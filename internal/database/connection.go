package database

import (
	"context"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepo struct {
	pool *pgxpool.Pool
}

// func (r *BaseRepo) WithTransaction(ctx context.Context, fn func(pgx.Tx) (any, error)) error {
// 	tx, err := r.pool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)
// 	_, err = fn(tx)
// 	if err != nil {
// 		tx.Rollback(ctx)
// 	} else {
// 		tx.Commit(ctx)
// 	}
// 	return nil
// }

func (r *BaseRepo) execute(statement string, args ...any) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.pool.Exec(ctx, statement, args...)
}

func (r *BaseRepo) query(statement string, args ...any) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

}

type DatabaseHandler struct {
	pool *pgxpool.Pool
	user *UserRepo
}

func Connect(config *pgxpool.Config) *DatabaseHandler {
	poolCon, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Got error when connect database")
		slog.Error(err.Error())
		return nil
	}

	return &DatabaseHandler{
		pool: poolCon,
		user: &UserRepo{
			BaseRepo: &BaseRepo{
				pool: poolCon,
			},
		},
	}
}

func (d *DatabaseHandler) GetUserConnection() *UserRepo {
	return d.user
}

func (d *DatabaseHandler) Execute(statement string, args ...any) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	d.pool.Exec(ctx, statement, args...)
}

func (d *DatabaseHandler) Disconnect() {
	d.pool.Close()
}
