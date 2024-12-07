# ベースイメージとして軽量なNode.jsイメージを使用
FROM node:23-alpine AS base

# グローバルでpnpmをインストール
RUN npm install -g pnpm

# 作業ディレクトリを設定
WORKDIR /usr/src/app

# パッケージ情報をコピー（依存関係をキャッシュするため）
COPY package.json pnpm-lock.yaml ./

# 依存関係をインストール
RUN pnpm install --frozen-lockfile

# ソースコードをコピー
COPY . .

# 必要なポートを公開（gRPCデフォルトポート: 50051）
EXPOSE 50051

# アプリケーションを実行
CMD ["pnpm", "start"]
