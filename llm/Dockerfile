# ベースイメージ（開発用ステージ）
FROM ghcr.io/astral-sh/uv:python3.13-bookworm-slim AS dev

WORKDIR /app

# 環境変数
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy

# 必要なツールとライブラリをインストール
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project --no-dev

# 依存関係のコピーとインストール（開発用を含む）
ADD . /app
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --dev


# ビルド用ステージ
FROM ghcr.io/astral-sh/uv:python3.13-bookworm-slim AS builder
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project --no-dev
ADD . /app
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --no-dev


# 本番用ステージ
FROM python:3.13-slim-bookworm AS prod
WORKDIR /app
# 必要なファイルと依存関係をコピー
COPY --from=builder --chown=app:app /app /app

# 仮想環境のパスを環境変数に設定
ENV PATH="/app/.venv/bin:$PATH"

# gRPCポートを公開
EXPOSE 50052

# エントリポイント
CMD ["python", "./src/server/main.py"]
