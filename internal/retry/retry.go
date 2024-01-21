package retry

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	RetryAttempts    = 3
	FirstRetryDelay  = 1 * time.Second
	SecondRetryDelay = 3 * time.Second
	ThirdRetryDelay  = 5 * time.Second
)

type ExecFunc func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)

func Exec(fn ExecFunc, ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	var res pgconn.CommandTag
	var err error

	for i := 0; i < RetryAttempts; i++ {
		res, err = fn(ctx, sql, args...)
		if err == nil {
			return res, err
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			switch i {
			case 0:
				time.Sleep(FirstRetryDelay)
			case 1:
				time.Sleep(SecondRetryDelay)
			case 2:
				time.Sleep(ThirdRetryDelay)
			}
		}
	}

	return res, err
}

type QueryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

func Query(fn QueryFunc, ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	var res pgx.Rows
	var err error

	for i := 0; i < RetryAttempts; i++ {
		res, err = fn(ctx, sql, args...)
		if err == nil {
			return res, err
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			switch i {
			case 0:
				time.Sleep(FirstRetryDelay)
			case 1:
				time.Sleep(SecondRetryDelay)
			case 2:
				time.Sleep(ThirdRetryDelay)
			}
		}
	}

	return res, err
}
