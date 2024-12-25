# cdk

## 必要なアプリ

```sh
npm install -g aws-cdk
```

## 認証情報を設定

環境変数はGoogle Driveを参照。
コマンドが無い場合は [AWS CLI install and update instructions](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/getting-started-install.html) を参照。

```sh
aws configure --profile alter-ego
```

## Deploy方法 (コマンド)

```sh
# AWSに初回デプロイする場合のみ。リージョンやアカウントレベルで変更がない場合は不要
cdk bootstrap --profile alter-ego
```
