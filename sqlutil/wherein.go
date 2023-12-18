package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func getSelectFn(dbtx any) func(dest any, query string, args ...any) error {
	switch dbtx := dbtx.(type) {
	case *sqlx.DB:
		return dbtx.Select
	case *sqlx.Tx:
		return dbtx.Select
	default:
		panic("invalid type")
	}
}

func getSelectContextFn(dbtx any) func(ctx context.Context, dest any, query string, args ...any) error {
	switch dbtx := dbtx.(type) {
	case *sqlx.DB:
		return dbtx.SelectContext
	case *sqlx.Tx:
		return dbtx.SelectContext
	default:
		panic("invalid type")
	}
}

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

func WhereIn[T any, S any, DBTX sqlx.DB | sqlx.Tx](dbtx *DBTX, query string, args []S) ([]T, error) {
	return whereIn[T](getSelectFn(dbtx), query, args)
}

func WhereInContext[T any, S any, DBTX sqlx.DB | sqlx.Tx](ctx context.Context, dbtx *DBTX, query string, args []S) ([]T, error) {
	selectFn := func(dest any, query string, args ...any) error {
		return getSelectContextFn(dbtx)(ctx, dest, query, args...)
	}

	return whereIn[T](selectFn, query, args)
}
