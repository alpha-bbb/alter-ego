# PostgreSQL

## ツールの導入

マイグレーションに [golang-migrate](https://github.com/golang-migrate/migrate) を使用している。  
導入の仕方は [migrate CLI](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md) を参照。

## マイグレーションファイルの作成

```sh
migrate create -ext sql -dir postgresql/migrations -seq "ファイル名"
```

## マイグレーション

```sh
migrate --path postgresql/migrations --database 'postgresql://root:password@localhost:5432/alterego?sslmode=disable' -verbose up
```

## データベースを直接見る

```sh
docker exec -it alter-ego-postgres-1 bash
psql -d alterego
```
