package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func whereIn[T any](query string, values []any, selectFn func(dest interface{}, query string, args ...interface{}) error) ([]T, error) {
	if len(values) == 0 {
		return []T{}, nil
	}

	uniqueValues := lo.Uniq(values)

	query, args, err := sqlx.In(query, uniqueValues)
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

func WhereIn[T any](db *sqlx.DB, query string, values []any) ([]T, error) {
	return whereIn[T](query, values, db.Select)
}

func TxWhereIn[T any](tx *sqlx.Tx, query string, values []any) ([]T, error) {
	return whereIn[T](query, values, tx.Select)
}

func WhereInContext[T any](ctx context.Context, db *sqlx.DB, query string, values []any) ([]T, error) {
	selectFn := func(dest interface{}, query string, args ...interface{}) error {
		return db.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](query, values, selectFn)
}

func TxWhereInContext[T any](ctx context.Context, tx *sqlx.Tx, query string, values []any) ([]T, error) {
	selectFn := func(dest interface{}, query string, args ...interface{}) error {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return whereIn[T](query, values, selectFn)
}
