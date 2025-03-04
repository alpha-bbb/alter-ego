#!/bin/bash

DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_DIR="$DIR/logs"
mkdir -p "$LOG_DIR"
TARGET_DIR="$DIR/../backend/cmd/clearstatecount"

while true; do
    now=$(date +%s)

    # OS判定して翌日の0時のUNIX時間を取得
    if [[ "$(uname)" == "Darwin" ]]; then
        # macOSの場合
        next_midnight=$(date -v+1d -v0H -v0M -v0S +%s)
    else
        # Linuxの場合
        next_midnight=$(date -d "tomorrow 00:00" +%s)
    fi

    # 次の0時までの秒数を計算
    sleep_seconds=$((next_midnight - now))
    echo "現在: $(date). 次の0時まで ${sleep_seconds} 秒待機します。"

    sleep 5

    cd "$TARGET_DIR" || {
        echo "ディレクトリ移動に失敗しました"
        exit 1
    }
    # 実行対象のコマンドを実行
    /usr/local/go/bin/go run main.go >>"$LOG_DIR/clear_state_count.log" 2>&1
done
