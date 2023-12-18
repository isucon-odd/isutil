package sqlutil

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var schema = `
DROP TABLE IF EXISTS users;
CREATE TABLE users (
	user_id    INTEGER PRIMARY KEY,
    first_name VARCHAR(80)  DEFAULT '',
    last_name  VARCHAR(80)  DEFAULT '',
	email      VARCHAR(250) DEFAULT ''
);
`

type User struct {
	UserId    int    `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func setupDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Smith", "john.smith@example.com")
	tx.NamedExec("INSERT INTO users (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &User{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"})
	tx.Commit()

	return db
}

func TestWhereInDB(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	testCases := []struct {
		name     string
		args     []int
		expected []User
	}{
		{
			name:     "Empty args",
			args:     []int{},
			expected: []User{},
		},
		{
			name: "Single arg",
			args: []int{1},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
			},
		},
		{
			name: "Multiple args",
			args: []int{1, 2},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
				{
					UserId:    2,
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jane.smith@example.com",
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			users, err := WhereIn[User](db, "SELECT * FROM users WHERE user_id IN (?)", tc.args)
			require.NoError(t, err)

			require.Equal(t, tc.expected, users)
		})
	}
}

func TestWhereInTx(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	tx := db.MustBegin()
	defer tx.Commit()

	testCases := []struct {
		name     string
		args     []int
		expected []User
	}{
		{
			name:     "Empty args",
			args:     []int{},
			expected: []User{},
		},
		{
			name: "Single arg",
			args: []int{1},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
			},
		},
		{
			name: "Multiple args",
			args: []int{1, 2},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
				{
					UserId:    2,
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jane.smith@example.com",
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			users, err := WhereIn[User](tx, "SELECT * FROM users WHERE user_id IN (?)", tc.args)
			require.NoError(t, err)

			require.Equal(t, tc.expected, users)
		})
	}
}

func TestWhereInContextDB(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)
	defer db.Close()

	testCases := []struct {
		name     string
		args     []int
		expected []User
	}{
		{
			name:     "Empty args",
			args:     []int{},
			expected: []User{},
		},
		{
			name: "Single arg",
			args: []int{1},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
			},
		},
		{
			name: "Multiple args",
			args: []int{1, 2},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
				{
					UserId:    2,
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jane.smith@example.com",
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			users, err := WhereInContext[User](ctx, db, "SELECT * FROM users WHERE user_id IN (?)", tc.args)
			require.NoError(t, err)

			require.Equal(t, tc.expected, users)
		})
	}
}

func TestWhereInContextTx(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)
	defer db.Close()

	tx := db.MustBegin()
	defer tx.Commit()

	testCases := []struct {
		name     string
		args     []int
		expected []User
	}{
		{
			name:     "Empty args",
			args:     []int{},
			expected: []User{},
		},
		{
			name: "Single arg",
			args: []int{1},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
			},
		},
		{
			name: "Multiple args",
			args: []int{1, 2},
			expected: []User{
				{
					UserId:    1,
					FirstName: "John",
					LastName:  "Smith",
					Email:     "john.smith@example.com",
				},
				{
					UserId:    2,
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jane.smith@example.com",
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			users, err := WhereInContext[User](ctx, tx, "SELECT * FROM users WHERE user_id IN (?)", tc.args)
			require.NoError(t, err)

			require.Equal(t, tc.expected, users)
		})
	}
}
