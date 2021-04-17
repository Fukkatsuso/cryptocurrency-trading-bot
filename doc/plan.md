# plan

## 概要

- イーサリアム（あるいは3〜10万円程度の仮想通貨）の自動取引bot
- 定期的にデータを取ってきて，買い時なら買う，売り時なら売る
- 管理画面付き

## 必要・欲しい機能

- 取引所APIから定期的にデータ取得
- データから取引シグナルを発行
- 取引実行
- 相場と取引情報をDBに保存
- ログイン機能(admin, guest)
- チャート，資産，取引データを表示するダッシュボード
- ゲストユーザには限定的な情報（チャート，資産も?）のみ表示
- 値動き+各種指数のグラフ化
- 取引停止・再開ボタン（管理者のみ利用可能）
- SlackかLINEで取引実行通知

## 構成

- DB: Cloud SQL
- Server: Cloud Run
- Scheduler: Cloud Scheduler
- CI/CD: GitHub Actions

取引に関わる処理（バックエンド）と情報表示（フロントエンド）で分けて開発する．

取引を定期実行するためCloud Schedulerを使用．
バックエンドのCloud Runは認証必須にする．
SchedulerとCloud Runの間にPub/Subを噛ませることが必要だと思っていたが，Schedulerから直接Cloud Runのサービスを認証付きで呼び出せるようなので，Pub/Subは不採用となった．

GitHub Actionsでコードのpushをトリガーにテスト・デプロイする．

![flow](flow.drawio.png)
