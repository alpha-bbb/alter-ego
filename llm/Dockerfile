# 開発環境用ステージ
FROM ghcr.io/astral-sh/uv:python3.12-bookworm-slim as dev
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy
ENV UV_PYTHON_DOWNLOADS=0

WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen

ADD . /app
RUN --mount=type=cache,target=/root/.cache/uv

# ビルド用ステージ
FROM ghcr.io/astral-sh/uv:python3.12-bookworm-slim AS builder
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy
ENV UV_PYTHON_DOWNLOADS=0

WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project --no-dev

ADD . /app
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --no-dev

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
