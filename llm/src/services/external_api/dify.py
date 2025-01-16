import json
import logging
import os
from typing import Dict, List

import httpx
from dotenv import load_dotenv
from google.protobuf.json_format import MessageToDict

from llm.v1.llm_pb2 import TalkHistory

# 環境変数の読み込み
load_dotenv()

logger = logging.getLogger(__name__)

DIFY_API_URL = os.getenv("DIFY_API_URL")
DIFY_API_KEY = os.getenv("DIFY_API_KEY")

if not DIFY_API_URL or not DIFY_API_KEY:
    raise ValueError(
        "DIFY_API_URL and DIFY_API_KEY must be set in environment variables"
    )


class DifyClient:
    """Dify APIクライアント"""

    def __init__(self) -> None:
        self.client = httpx.AsyncClient(
            base_url=DIFY_API_URL,
            headers={
                "Authorization": f"Bearer {DIFY_API_KEY}",
                "Content-Type": "application/json",
            },
            timeout=30.0,  # タイムアウトを30秒に設定
        )

    async def close(self) -> None:
        """クライアントをクローズ"""
        await self.client.aclose()

    async def __aenter__(self) -> "DifyClient":
        """非同期コンテキストマネージャーのエントリー"""
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:
        """非同期コンテキストマネージャーのイグジット"""
        await self.close()


# シングルトンクライアントのインスタンス
dify_client = DifyClient()


async def call_dify_talk(input_histories: list[TalkHistory]) -> dict[str, list[str]]:
    """
    Dify APIを呼び出してチャット応答を取得

    Args:
        input_histories: List[TalkHistory] - チャット履歴

    Returns:
        Dict[str, List[str]] - Difyからの応答 {"replies": [str, ...]}

    Raises:
        Exception: API呼び出しに失敗した場合
    """
    try:
        response = await dify_client.client.post(
            "/chat-messages",
            json={
                "query": json.dumps([MessageToDict(history) for history in input_histories]),
                "inputs": {},
                "response_mode": "blocking",
                "user": "abc-123",
                "conversation_id": "",
                "files": [],
                "auto_generate_name": True,
            },
        )

        response.raise_for_status()
        data = response.json()

        # 応答をパース
        return json.loads(data["answer"])

    except httpx.HTTPError as e:
        logger.error(f"HTTP error occurred: {e}", exc_info=True)
        raise Exception("Dify API call failed")
    except json.JSONDecodeError as e:
        logger.error(f"JSON decode error: {e}", exc_info=True)
        raise Exception("Invalid response from Dify API")
    except Exception as e:
        logger.error(f"Unexpected error: {e}", exc_info=True)
        raise
