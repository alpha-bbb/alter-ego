import json
import logging
import os

import httpx
from dotenv import load_dotenv
from google.protobuf.json_format import MessageToDict

from llm.v1.llm_pb2 import TalkHistory

# 環境変数の読み込み
load_dotenv()

logger = logging.getLogger(__name__)

DIFY_API_URL = os.getenv("DIFY_API_URL", "https://api.dify.ai/v1")
DIFY_API_KEY = os.getenv("DIFY_API_KEY", "")
DIFY_KNOWLEDGE_BASE_API_KEY = os.getenv("DIFY_KNOWLEDGE_BASE_API_KEY", "")

if not DIFY_API_URL or not DIFY_API_KEY:
    raise ValueError(
        "DIFY_API_URL and DIFY_API_KEY must be set in environment variables"
    )


class DifyClient:
    """Dify APIクライアント"""

    def __init__(self) -> None:
        raise NotImplementedError(
            "DifyClient is a singleton class. Use the instance dify_client instead."
        )

    async def close(self) -> None:
        """クライアントをクローズ"""
        raise NotImplementedError(
            "DifyClient is a singleton class. Use the instance dify_client instead."
        )


class DifyChatClient(DifyClient):
    """Dify Chat APIクライアント"""

    def __init__(self) -> None:
        self.client = httpx.AsyncClient(
            base_url=DIFY_API_URL,
            headers={
                "Authorization": f"Bearer {DIFY_API_KEY}",
                "Content-Type": "application/json",
            },
        )

    async def close(self) -> None:
        """クライアントをクローズ"""
        await self.client.aclose()


class DifyKnowledgeBaseClient(DifyClient):
    def __init__(self) -> None:
        self.client = httpx.AsyncClient(
            base_url=DIFY_API_URL,
            headers={
                "Authorization": f"Bearer {DIFY_KNOWLEDGE_BASE_API_KEY}",
                "Content-Type": "application/json",
            },
        )

    async def close(self) -> None:
        """クライアントをクローズ"""
        await self.client.aclose()


# シングルトンクライアントのインスタンス
dify_chat_client = DifyChatClient()
dify_knowledge_base_client = DifyKnowledgeBaseClient()


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
        response = await dify_chat_client.client.post(
            "/chat-messages",
            json={
                "query": json.dumps(
                    [MessageToDict(history) for history in input_histories]
                ),
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


async def call_dify_retrieve(
    dataset_id: str, query: str, retrieval_model: dict
) -> dict:
    """
    Dify APIの /datasets/{dataset_id}/retrieve エンドポイントを呼び出して
    Knowledge Base からチャンクを取得する。

    Args:
        dataset_id (str): 取得対象の Knowledge ID。
        query (str): 検索キーワード。
        retrieval_model (dict): 検索に使用する retrieval_model の設定
            （例:
            {
                "search_method": "keyword_search",
                "reranking_enable": False,
                "reranking_mode": None,
                "reranking_model": {
                    "reranking_provider_name": "",
                    "reranking_model_name": ""
                },
                "weights": None,
                "top_k": 1,
                "score_threshold_enabled": False,
                "score_threshold": None
            }）

    Returns:
        dict: APIから返されたレスポンス（例:
            {
              "query": { "content": "test" },
              "records": [ ... ]
            }
        )

    Raises:
        Exception: API呼び出しに失敗した場合
    """
    try:
        response = await dify_knowledge_base_client.client.post(
            f"/datasets/{dataset_id}/retrieve",
            json={
                "query": query,
                "retrieval_model": retrieval_model,
            },
        )
        response.raise_for_status()
        return response.json()

    except httpx.HTTPError as e:
        logger.error(f"HTTP error occurred in retrieve endpoint: {e}", exc_info=True)
        raise Exception("Dify API call failed for retrieve endpoint")
    except json.JSONDecodeError as e:
        logger.error(f"JSON decode error in retrieve endpoint: {e}", exc_info=True)
        raise Exception("Invalid response from Dify API in retrieve endpoint")
    except Exception as e:
        logger.error(f"Unexpected error in retrieve endpoint: {e}", exc_info=True)
        raise
