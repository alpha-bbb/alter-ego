# backend

## ローカル環境立ち上げ

```sh
make run
```

## gRPCを直接叩くコマンド例

```sh
buf curl --protocol grpc --http2-prior-knowledge \
  --schema ./proto \
  --data '{
    "histories": [
        {
            "date": "2024-12-07",
            "user": {
                "userId": "123",
                "name": "太郎",
                "role": 2
            },
            "message": "今度ご飯にいきませんか？"
        }
    ],
    "actionKind": 2
  }' \
  http://localhost:50051/backend.v1.BackendService/Talk
```

## Stripe

### 開発用のローカルwebhookを立てる

```sh
stripe listen --forward-to localhost:50080/stripe/webhook
```

## cmd

```sh
docker exec -it alter-ego-backend-1 sh -c "cd /app/backend/cmd/clearstatecount && go run main.go"
```
