FROM python:3.12-slim-bookworm AS dev
RUN apt-get update && apt-get install -y --no-install-recommends build-essential

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ADD . /app
WORKDIR /app

RUN uv sync --frozen

# ビルド用ステージ
FROM python:3.12-slim-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends build-essential

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ADD . /app

WORKDIR /app
RUN uv sync --frozen

# 本番用ステージ
FROM python:3.12-slim-bookworm AS prod
WORKDIR /app

COPY --from=builder --chown=app:app /app /app

# 仮想環境のパスを環境変数に設定
ENV PATH="/app/.venv/bin:$PATH"

# gRPCポートを公開
EXPOSE 50052

# エントリポイント
CMD ["python", "./src/server/main.py"]
