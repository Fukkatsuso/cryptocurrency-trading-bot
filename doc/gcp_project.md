# GCPのプロジェクトに関して

## Setup

### Cloud Shell上での作業

#### 1. プロジェクトの作成

```sh
export PROJECT_ID=trading-xxxxxx
export REGION=asia-northeast1
gcloud projects create ${PROJECT_ID} --name=trading
gcloud config set project ${PROJECT_ID}
gcloud config set run/region ${REGION}
```

参考: https://cloud.google.com/sdk/gcloud/reference/projects/create

#### 2. 課金の有効化

```sh
gcloud alpha billing accounts list
```

以下のようなアカウントのリストが出てくる．

```
ACCOUNT_ID            NAME              OPEN  MASTER_ACCOUNT_ID
XXXXXX-YYYYYY-ZZZZZZ  請求先アカウント     True
```

ACCOUNT_IDを使って課金を有効化する

```sh
gcloud alpha billing projects link ${PROJECT_ID} --billing-account <ACCOUNT_ID>
gcloud services enable cloudbilling.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

#### 3. APIの有効化

```sh
gcloud services enable run.googleapis.com
gcloud services enable sql-component.googleapis.com sqladmin.googleapis.com
```

#### 4. サービスアカウント, サービスアカウントキーの作成

```sh
export SA_NAME=github-actions
gcloud iam service-accounts create ${SA_NAME} \
  --description="used by GitHub Actions" \
  --display-name="${SA_NAME}"
gcloud iam service-accounts list # サービスアカウントが作られたことを確認

export IAM_ACCOUNT=${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com
gcloud iam service-accounts keys create ~/${PROJECT_ID}/${SA_NAME}/key.json \
  --iam-account ${IAM_ACCOUNT}
```

#### 5. サービスアカウントにroleを付与

```sh
export PROJECT_NUMBER=XXXXXXXXXXXX # GCPコンソールの「プロジェクト情報」に「プロジェクト番号」として表示されている数字

gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/run.admin"
gcloud iam service-accounts add-iam-policy-binding ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/cloudsql.admin"
```

#### 6. Cloud SQL for MySQLインスタンスの作成とDB初期化

インスタンスの作成

- MySQL 5.7
- db-f1-micro

```sh
export CLOUDSQL_INSTANCE=trading-mysql
gcloud sql instances create ${CLOUDSQL_INSTANCE} --database-version=MYSQL_5_7 --region=${REGION} --tier=db-f1-micro
```

参考: https://cloud.google.com/sql/docs/mysql/create-instance?hl=ja#gcloud

DB初期化

- ユーザの作成
- DBの作成

```sh
export MYSQL_USER=hoge
export MYSQL_PASSWORD=fuga
gcloud sql users create ${MYSQL_USER} --instance=${CLOUDSQL_INSTANCE} --password=${MYSQL_PASSWORD}
gcloud sql users list -i ${CLOUDSQL_INSTANCE} # ユーザが作成されたかチェック

export MYSQL_DATABASE=trading_db
gcloud sql databases create ${MYSQL_DATABASE} --instance=${CLOUDSQL_INSTANCE} --charset=utf8
gcloud sql databases list --instance=${CLOUDSQL_INSTANCE} # DBが作成されたかチェック
```

参考（ユーザの作成）: https://cloud.google.com/sdk/gcloud/reference/sql/users/create
参考（DBの作成）: https://cloud.google.com/sdk/gcloud/reference/sql/databases/create

### GitHubでの作業

#### 環境変数を設定

```
GCP_PROJECT: プロジェクトID
GCP_REGION: リージョン
GCP_SA_KEY_JSON: サービスアカウントのJSON鍵
GCP_SA_KEY: サービスアカウントのJSON鍵をBase64エンコード
CLOUDSQL_INSTANCE: Cloud SQLのインスタンス名
MYSQL_USER: 作成したユーザ
MYSQL_PASSWORD: ユーザのパスワード
MYSQL_HOST: Cloud SQLインスタンスのパブリックIP
MYSQL_PORT: ポート番号
MYSQL_DATABASE: 作成したデータベース
```

`GCP_SA_KEY`の取得方法

```sh
# Cloud Shell上で
openssl base64 -in ~/${PROJECT_ID}/${SA_NAME}/key.json
```
