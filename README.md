# isucon-odd/isutil

ISUCON用のGo言語ユーティリティです。

## インストール

```sh
go get github.com/isucon-odd/isutil
```

## モジュールの内容

### cache

ジェネリクスを使ったキャッシュです。有効期限の設定も可能です。

```go
package main

import (
	"time"

	"github.com/isucon-odd/isutil/cache"
)

func main() {
	cache := cache.NewCache[string]()

	cache.Set("key", "value")
	cache.Get("key")
	cache.Delete("key")

	cache.SetWithExpiration("key", "value", time.Minute*10)
	cache.GetAndDeleteExpired("key")
}
```

### sqlutil

SQLクエリの実行を簡単にするためのユーティリティです。

#### WHERE IN

```go
package main

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/isucon-odd/isutil/sqlutil"
)

func main(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) {
	var users []User
	var err error

	query := "SELECT * FROM users WHERE id IN (?)"
	userIDs := []int{1, 2, 3}

	users, err = WhereIn[User](db, query, userIDs)
	users, err = WhereInContext[User](ctx, db, query, userIDs)

	users, err = TxWhereIn[User](tx, query, userIDs)
	users, err = TxWhereInContext[User](ctx, tx, query, userIDs)
}
```
