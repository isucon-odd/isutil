package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func whereIn[T any, S any](selectFn func(dest any, query string, args ...any) error, query string, args []S) ([]T, error) {
	if len(args) == 0 {
		return []T{}, nil
	}

	_query, _args, err := sqlx.In(query, []any{args}...)
	if err != nil {
		return []T{}, err
	}

	var result []T

	err = selectFn(&result, _query, _args...)
	if err != nil {
		return []T{}, err
	}

	return result, nil
}

func WhereIn[T any, S any](db *sqlx.DB, query string, args []S) ([]T, error) {
	return whereIn[T](db.Select, query, args)
}

func TxWhereIn[T any, S any](tx *sqlx.Tx, query string, args []S) ([]T, error) {
	return whereIn[T](tx.Select, query, args)
}

func WhereInContext[T any, S any](ctx context.Context, db *sqlx.DB, query string, args []S) ([]T, error) {
	selectFn := func(dest any, query string, args ...any) error {
		return db.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](selectFn, query, args)
}

func TxWhereInContext[T any, S any](ctx context.Context, tx *sqlx.Tx, query string, args []S) ([]T, error) {
	selectFn := func(dest any, query string, args ...any) error {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](selectFn, query, args)
}
