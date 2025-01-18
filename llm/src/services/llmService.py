import logging
from typing import List

import grpc
from llm.v1 import llm_pb2, llm_pb2_grpc
from src.services.external_api.dify import call_dify_talk

logger = logging.getLogger(__name__)


class LlmServiceServicer(llm_pb2_grpc.LlmServiceServicer):
    """LLMサービスの実装クラス"""

    async def Talk(
        self, request: llm_pb2.TalkRequest, context: grpc.aio.ServicerContext
    ) -> llm_pb2.TalkResponse:
        """
        chatの実装

        Args:
            request: TalkRequest - クライアントからのリクエスト
            context: ServicerContext - gRPCのコンテキスト

        Returns:
            TalkResponse - LLMからの応答

        Raises:
            grpc.RpcError: gRPCエラーが発生した場合
        """
        try:
            # 履歴を300件に制限
            histories = (
                request.histories[-300:]
                if len(request.histories) > 300
                else request.histories
            )

            # Dify APIを呼び出し
            messages = await call_dify_talk(histories)

            # レスポンスの作成
            return llm_pb2.TalkResponse(message=messages["replies"])

        except Exception as e:
            logger.error(f"Error in Talk: {str(e)}", exc_info=True)
            await context.abort(
                code=grpc.StatusCode.INTERNAL, details="Internal server error"
            )
