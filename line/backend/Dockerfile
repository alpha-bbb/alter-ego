# 開発用ビルドステージ
FROM node:22-alpine AS dev

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install

COPY . ./

# 本番用ビルドステージ
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --immutable --immutable-cache --check-cache

COPY . ./
RUN yarn run build

# ランタイム用ステージ
FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY package.json ./
COPY ./google_auth.json ./google_auth.json

EXPOSE 3000

CMD ["node", "dist/main.js"]
