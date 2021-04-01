# migrationについて

## CLIツールのインストール (Windows, MacOS, or Linux)

[golang-migrate](https://github.com/golang-migrate/migrate)を使う．

[ここ](https://github.com/golang-migrate/migrate/releases)から必要なバイナリを`curl`で取得し，`tar`で解凍する．

```sh
cd /usr/local/bin
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.darwin-amd64.tar.gz | tar xvz
mv migrate.darwin-amd64 migrate
```

参考: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

## migrationファイル作成の例

```sh
migrate create -ext sql -dir db/migrations -seq create_users_table
```

以下の空ファイルが作られる．
upの方に更新系のSQLを，downの方に「upの操作を取り消す」SQLを書く．

- db/migrations/xxxxxx_create_users_table.up.sql
- db/migrations/xxxxxx_create_users_table.down.sql

## migrationの実行 (MySQL)

上げるとき

```sh
migrate -path db/migrations/ -database 'mysql://user:password@tcp(host:port)/dbname' up <migrationを上げる数>
```

下げるとき

```sh
migrate -path db/migrations/ -database 'mysql://user:password@tcp(host:port)/dbname' down <migrationを下げる数>
```

例)

```sh
migrate -path db/migrations/ -database 'mysql://trading_app:password@tcp(localhost:3306)/trading_db' up 1
```

注意: 数を指定しない場合はmigrationが全て実行される

参考: https://github.com/golang-migrate/migrate/tree/master/database/mysql
