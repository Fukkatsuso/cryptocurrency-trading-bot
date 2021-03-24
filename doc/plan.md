# plan

## 概要

- イーサリアム（あるいは3〜10万円程度の仮想通貨）の自動取引bot
- 定期的にデータを取ってきて，買い時なら買う，売り時なら売る
- 管理画面付き

## 必要・欲しい機能

- 取引所APIから定期的にデータ取得
- データから取引シグナルを発行
- 取引実行
- 板情報と取引情報をDBに保存
- ログイン機能(admin, guest)
- 板情報，資産，取引データを表示するダッシュボード
- ゲストユーザには限定的な情報（板情報，資産も?）のみ表示
- 値動き+各種指数のグラフ化
- 取引停止・再開ボタン（管理者のみ利用可能）
- SlackかLINEで取引実行通知

## 構成

- DB: Cloud SQL
- Server: Cloud Run
- Scheduler: Cloud Scheduler, Pub/Sub
- CI/CD: GitHub Actions

取引に関わる処理（バックエンド）と情報表示（フロントエンド）で分けて開発する．

取引を定期実行するためCloud SchedulerとPub/Subを使用．
バックエンドのCloud Runは認証必須にする．

GitHub Actionsでコードのpushをトリガーにテスト・デプロイする．

![flow](flow.drawio.png)
