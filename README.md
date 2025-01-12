# alter-ego

## ディレクトリ構成

```plaintext
.
├── backend         # 中間層
├── cdk             # cdk
├── line            # LINE関連
│   ├── backend     # Messaging API and LIFF用backend
│   └── liff        # LIFFアプリ
├── llm             # 言語処理
└── proto           # schema
```

## Docker立ち上げ

```sh
docker compose up --build
docker compose -f nginx/docker-compose.yml up --build
```

各README.mdを参照。

