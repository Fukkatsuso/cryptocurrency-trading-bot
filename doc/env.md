# 環境変数について

## ローカル環境

プロジェクトのルートディレクトリに以下の内容の`.env`ファイルを作成する

```
MYSQL_USER=<MySQLに作成したユーザ>
MYSQL_PASSWORD=<ユーザのパスワード>
MYSQL_HOST=db
MYSQL_PORT=3306
MYSQL_DATABASE=<データベース名>
BITFLYER_API_KEY=<bitflyerのAPIキー>
BITFLYER_API_SECRET=<bitflyerのAPIシークレット>
PRODUCT_CODE=ETH_JPY
```

## 本番環境(GCP)

GitHubのリポジトリのSettings=>Secretsで各変数を設定する

```
GCP_PROJECT: <GCPのプロジェクトID>
GCP_REGION: <Cloud RunとCloud SQLのリージョン>
GCP_SA_KEY: <サービスアカウントのJSON鍵>
CLOUDSQL_INSTANCE: <Cloud SQLのインスタンス名>
MYSQL_USER: <MySQLに作成したユーザ>
MYSQL_PASSWORD: <ユーザのパスワード>
MYSQL_DATABASE: <データベース名>
BITFLYER_API_KEY: <bitflyerのAPIキー>
BITFLYER_API_SECRET: <bitflyerのAPIシークレット>
PRODUCT_CODE: ETH_JPY
```
