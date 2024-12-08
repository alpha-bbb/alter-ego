# backend

## とりあえずの立ち上げ

```sh
go run ./main.go
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
                "name": "太郎"
            },
            "message": "今度ご飯にいきませんか？"
        }
    ],
    "actionKind": 2
  }' \
  http://localhost:50051/backend.v1.BackendService/Talk
```
