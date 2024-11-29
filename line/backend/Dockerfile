FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install --immutable --immutable-cache --check-cache

COPY . ./

RUN yarn run build

FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 80

CMD ["node", "dist/main.js"]
