# alter-ego Messaging API

## 環境変数

`.env.example` をコピーして `.env` を作成する。

## デプロイ方法

### Draft

```sh
netlify deploy
```

### Product

```sh
netlify deploy --prod
```

### netlifyのwebhook

デプロイ先のURLに `.netlify/functions/webhook` を付け足す必要がある。

例: `https://673f51e364be3a0b8c1c84ad--magenta-horse-85318c.netlify.app/.netlify/functions/webhook`

### ローカルのコードをwebhookで試したい

ngrokがおすすめ。下記のサイトが参考になる。
[Node.js & TypeScriptでLINEBot入門（1）：チャットボット開発の流れと実践方法 | Go-Tech Blog](https://go-tech.blog/nodejs/line-chat-bot/)

## 参考資料

[Messaging APIの概要 | LINE Developers](https://developers.line.biz/ja/docs/messaging-api/overview/)
