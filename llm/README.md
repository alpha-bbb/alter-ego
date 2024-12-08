# llm

## ローカル環境での実行方法

### 環境変数の設定

.envの設定
設定値は [Google Docs](https://docs.google.com/document/d/1A8ve0_vrlIVE01P5dR_D02quT08GBfmuvQ8y7J5XjhE/edit?usp=sharing) を参照。

### ツールの導入

```sh
npm install -g pnpm
```

### パッケージのインストール

```sh
pnpm install --frozen-lockfile
```

### 開発環境の立ち上げ

```sh
pnpm dev
```

## gRPCを直接叩くコマンド例

```sh
cd proto
buf curl --protocol grpc --http2-prior-knowledge \
  --schema ./proto \
  --data '{
    "histories": [
        {
            "date": "2024-12-07",
            "user": {
                "userId": "123",
                "name": "太郎"
            },
            "message": "今度ご飯にいきませんか？"
        }
    ],
    "actionKind": 2
  }' \
  http://localhost:8080/llm.v1.LlmService/Talk
```
