# ベースイメージ（開発用ステージ）
FROM ghcr.io/astral-sh/uv:python3.13-bookworm-slim AS dev

WORKDIR /app

# 環境変数
ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1 UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy

# 必要なツールとライブラリをインストール
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 依存関係のコピーとインストール（開発用を含む）
COPY pyproject.toml uv.lock /app/
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --dev

# ソースコードをコピー
COPY . /app

# 本番用ステージ
FROM python:3.13-slim-bookworm AS prod

WORKDIR /app

# 必要なファイルと依存関係をコピー
COPY --from=dev /app /app
COPY --from=dev /app/.venv /app/.venv

# 仮想環境のパスを環境変数に設定
ENV PATH="/app/.venv/bin:$PATH"

# gRPCポートを公開
EXPOSE 50052

# エントリポイント
CMD ["python", "src/server/main.py"]
