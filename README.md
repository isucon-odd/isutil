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
