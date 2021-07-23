# GCPのプロジェクトに関して

## Caution

### アプリを一時的に止める際のCloud SQLインスタンス

当然インスタンスを停止させるが，パブリックIPアドレスが設定されたままだと何もしなくても1日40円くらい課金されてしまう．
GCPコンソールの`SQL`=>`接続`=>`ネットワーキング`から「プライベートIP」に切り替えておくこと．
その上でインスタンスを停止．

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

- 参考: https://cloud.google.com/sdk/gcloud/reference/projects/create

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
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/storage.admin"
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

- 参考: https://cloud.google.com/sql/docs/mysql/create-instance?hl=ja#gcloud

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

- 参考（ユーザの作成）: https://cloud.google.com/sdk/gcloud/reference/sql/users/create
- 参考（DBの作成）: https://cloud.google.com/sdk/gcloud/reference/sql/databases/create

### GitHubでの作業

#### 環境変数を設定

- doc/env.mdを参照
- `GCP_SA_KEY`は，Cloud Shellの`~/${PROJECT_ID}/${SA_NAME}/key.json`に作成済み

## Job Scheduling

### Cloud Schedulerジョブの作成（Cloud Runにデプロイできてから）

ticker取得や売買を定期的に行うため，ジョブを作成する．

以下のコマンドをCloud Shell上で実行する．
パラメータは適宜変更．
`IAM_ACCOUNT`は以前作成したもの．

```sh
export JOB_NAME=job
export SERVICE_URL=https://xxx/yyy
# http-methodはGET, POST, PUT, DELETE, HEADのうちどれか
gcloud beta scheduler jobs create http ${JOB_NAME} \
  --schedule "5 * * * *" \
  --http-method="POST" \
  --uri="${SERVICE_URL}" \
  --oidc-service-account-email="${IAM_ACCOUNT}" \
  --oidc-token-audience="${SERVICE_URL}" \
  --time-zone="Asia/Tokyo"
```

- 参考（スケジューリング）: https://cloud.google.com/run/docs/triggering/using-scheduler?hl=ja
- 参考（createコマンド）: https://cloud.google.com/sdk/gcloud/reference/beta/scheduler/jobs/create

### crontabの書き方の例

`* * * * *`で左から「分」「時」「日」「月」「曜日」の順

- 5分おき: `*/5 * * * *`
- 毎日9時: `00 9 * * *`

## DBクライアントをCloud SQLに接続

Cloud SQL Proxyを立てる

```sh
# sa_keyファイルにGCPサービスアカウントのキーを保存しておく
docker run \
  -v $PWD/sa_key:/config \
  -p 127.0.0.1:3306:3306 \
  gcr.io/cloudsql-docker/gce-proxy:1.19.1 /cloud_sql_proxy \
  -instances=$GCP_PROJECT:$GCP_REGION:$CLOUDSQL_INSTANCE=tcp:0.0.0.0:3306 \
  -credential_file=/config
```

DBクライアントで以下を設定すれば接続可能になる

```
Ver.: MySQL 5.x
Host: 127.0.0.1
Port: 3306
User: (MYSQL_USERと同じ)
Password: (MYSQL_PASSWORDと同じ)
Database: (MYSQL_DATABASEと同じ)
```

## Cloud Scheduler, Cloud Functinosを用いてCloud SQL データベースをエクスポート

- https://cloud.google.com/solutions/scheduling-cloud-sql-database-exports-using-cloud-scheduler?hl=ja
- https://github.com/Fukkatsuso/cryptocurrency-trading-bot/issues/82

### APIの有効化

```sh
gcloud services enable \
  sqladmin.googleapis.com \
  cloudfunctions.googleapis.com \
  cloudscheduler.googleapis.com \
  appengine.googleapis.com \
  cloudresourcemanager.googleapis.com
```

### カスタムロールを作成

```sh
export PROJECT_ID=trading-xxxxxx
export STORAGE_ROLE="simpleStorageRole"
export SQL_ROLE="sqlExporter"
gcloud iam roles create ${STORAGE_ROLE} --project ${PROJECT_ID} \
    --title "Simple Storage Role" \
    --description "Grant permissions to view and create objects in Cloud Storage" \
    --permissions "storage.objects.create,storage.objects.get"

gcloud iam roles create ${SQL_ROLE} --project ${PROJECT_ID} \
    --title "SQL Exporter Role" \
    --description "Grant permissions to export data from a Cloud SQL instance to a Cloud Storage bucket as a SQL dump or CSV file" \
    --permissions "cloudsql.instances.export"
```

### Cloud Storage バケットを作成する

```sh
export REGION=asia-northeast1
export BUCKET_NAME=xxxxxx-bucket
gsutil mb -l ${REGION} gs://${BUCKET_NAME}
```

### サービスアカウントに権限を付与する

```sh
# gcloud sql instances list
export SQL_INSTANCE=trading-mysql
export SQL_SA=(`gcloud sql instances describe ${SQL_INSTANCE} \
  --project ${PROJECT_ID} \
  --format "value(serviceAccountEmailAddress)"`)
gsutil iam ch serviceAccount:${SQL_SA}:projects/${PROJECT_ID}/roles/${STORAGE_ROLE} gs://${BUCKET_NAME}

# gcloud iam service-accounts list
export GCF_NAME=github-actions
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:${GCF_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="projects/${PROJECT_ID}/roles/${SQL_ROLE}"

# デプロイ時発生した権限周りのエラーの対処
export SA_NAME=github-actions
export IAM_ACCOUNT=${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com
gcloud projects add-iam-policy-binding  ${PROJECT_ID} \
  --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/cloudfunctions.developer"
gcloud iam service-accounts add-iam-policy-binding ${PROJECT_ID}@appspot.gserviceaccount.com \
  --member="serviceAccount:${IAM_ACCOUNT}" \
  --role="roles/iam.serviceAccountUser"
```

### Cloud Scheduler ジョブを作成

```sh
export JOB_NAME=export-db
export FUNCTION_NAME=export-database-to-storage
export SERVICE_URL="https://${REGION}-${PROJECT_ID}.cloudfunctions.net/${FUNCTION_NAME}"
gcloud scheduler jobs create http ${JOB_NAME} \
  --schedule "57 8 * * *" \
  --http-method="POST" \
  --uri="${SERVICE_URL}" \
  --oidc-service-account-email="${IAM_ACCOUNT}" \
  --oidc-token-audience="${SERVICE_URL}" \
  --time-zone="Asia/Tokyo"
```
