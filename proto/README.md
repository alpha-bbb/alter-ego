# Protocol Buffers

## 概要

マイクロサービス間の通信におけるインターフェイスを定義する。  

## ツール

generateなどに `buf` を使用している。  
[Install the Buf CLI - Buf Docs](https://buf.build/docs/installation/)

### Plugins

下記のコマンドを実行して必要なものを導入しておく必要がある。

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Generate

```sh
buf generate
```
