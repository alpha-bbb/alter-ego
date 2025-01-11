# nginx

## Let's Encryptのsslを生成

```sh
sudo certbot certonly --webroot -w html -d alter-ego.jtj.jp -d api.line.alter-ego.jtj.jp -d dify.alter-ego.jtj.jp
```

## Let's Encryptのsslにドメインを追加

```sh
sudo certbot certonly --webroot -w html --force-renew --cert-name fuu.jp -d alter-ego.jtj.jp -d api.line.alter-ego.jtj.jp -d dify.alter-ego.jtj.jp
```

## basic認証ファイルの生成

```sh
htpasswd -c -b .htpasswd {user} {password}
```
