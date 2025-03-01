# nginx

## Let's Encryptのsslを生成

```sh
sudo certbot certonly --webroot -w html -d alter-ego.fuu.jp -d api.line.alter-ego.fuu.jp -d dify.alter-ego.fuu.jp -d api.backend.alter-ego.fuu.jp
```

## Let's Encryptのsslにドメインを追加

```sh
sudo certbot certonly --webroot -w html --force-renew --cert-name fuu.jp -d alter-ego.fuu.jp -d api.line.alter-ego.fuu.jp -d dify.alter-ego.fuu.jp -d api.backend.alter-ego.fuu.jp
```

## Dify用basic認証ファイルの生成

```sh
htpasswd -c -b dify/docker/nginx/conf.d/.htpasswd {user} {password}
```

```diff
--- a/docker/nginx/conf.d/default.conf.template
+++ b/docker/nginx/conf.d/default.conf.template
@@ -25,6 +25,8 @@ server {
     }

     location / {
+      auth_basic "Restricted";
+      auth_basic_user_file /etc/nginx/conf.d/.htpasswd;
       proxy_pass http://web:3000;
       include proxy.conf;
     }
```
