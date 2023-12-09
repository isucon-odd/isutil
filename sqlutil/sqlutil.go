package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func whereIn[T any](selectFn func(dest any, query string, args ...any) error, query string, args ...any) ([]T, error) {
	if len(args) == 0 {
		return []T{}, nil
	}

	uniqArgs := lo.Uniq(args)

	query, args, err := sqlx.In(query, uniqArgs)
	if err != nil {
		return []T{}, err
	}

	result := []T{}

	err = selectFn(&result, query, args...)
	if err != nil {
		return []T{}, err
	}

	return result, nil
}

func WhereIn[T any](db *sqlx.DB, query string, args ...any) ([]T, error) {
	return whereIn[T](db.Select, query, args)
}

func TxWhereIn[T any](tx *sqlx.Tx, query string, args ...any) ([]T, error) {
	return whereIn[T](tx.Select, query, args)
}

func WhereInContext[T any](ctx context.Context, db *sqlx.DB, query string, args ...any) ([]T, error) {
	selectFn := func(dest any, query string, args ...any) error {
		return db.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](selectFn, query, args)
}

func TxWhereInContext[T any](ctx context.Context, tx *sqlx.Tx, query string, args ...any) ([]T, error) {
	selectFn := func(dest any, query string, args ...any) error {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](selectFn, query, args)
}
