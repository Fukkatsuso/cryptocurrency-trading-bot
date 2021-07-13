# cryptocurrency-trading-bot

[![Deploy](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/deploy.yml/badge.svg)](https://github.com/Fukkatsuso/cryptocurrency-trading-bot/actions/workflows/deploy.yml)

## About

![flow](doc/flow.drawio.png)

- bitflyer APIを使ってイーサリアムを取引するbot
- チャート表示やバックテストができる管理画面付き
- GCPで動かす

## Documents

- [GCPプロジェクトのセットアップ](doc/gcp_project.md)
- [環境変数の設定](doc/env.md)
- [DBのマイグレーション](doc/migration.md)

## Running the app

[ドキュメント](doc/env.md)を参考に環境変数を設定する．

db, trader, dashboard, schedulerを起動．

```bash
$ docker compose up
```

http://localhost:8080 で管理画面を開ける．
