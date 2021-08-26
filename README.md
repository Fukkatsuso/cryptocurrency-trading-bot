# cryptocurrency-trading-bot

[![Test](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/test.yml/badge.svg)](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/test.yml)
[![Deploy](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/deploy.yml/badge.svg)](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/deploy.yml)

## About

![flow](doc/flow.drawio.png)

- bitflyer APIを使ってイーサリアムを取引するbot
- チャート表示やバックテストができる管理画面付き
  - https://dashboard-qboqrrfuaq-an.a.run.app
- GCPで動かす

## Documents

- [GCPプロジェクトのセットアップ](doc/gcp_project.md)
- [環境変数の設定](doc/env.md)
- [DBのマイグレーション](doc/migration.md)

## Running the app

[ドキュメント](doc/env.md)を参考に環境変数を設定する．

テストを実行する場合，GCPのサービスアカウントのキーを`sa_key`というファイルでプロジェクトルートに配置しておくこと．

db, trader, dashboard, schedulerを起動．

```bash
$ docker compose up
```

データベースのマイグレーション（詳しくは[ドキュメント](doc/migration.md)参照）

```sh
# example
migrate -path db/migrations/ -database 'mysql://trading_app:password@tcp(localhost:3306)/trading_db' up
```

http://localhost:8080 で管理画面を開ける．
