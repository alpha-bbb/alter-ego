import asyncio
import logging
from concurrent import futures

from grpc import aio

from llm.v1.llm_pb2_grpc import add_LlmServiceServicer_to_server
from src.services.llmService import LlmServiceServicer

# ロギングの設定
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

async def serve() -> None:
    # gRPCサーバーの設定
    server = aio.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[
            ('grpc.max_send_message_length', 100 * 1024 * 1024),
            ('grpc.max_receive_message_length', 100 * 1024 * 1024),
        ]
    )

    # LLMサービスの実装をサーバーに登録
    llm_service = LlmServiceServicer()
    add_LlmServiceServicer_to_server(llm_service, server)

    # サーバーのポート設定
    listen_addr = '[::]:50052'
    server.add_insecure_port(listen_addr)

    logger.info(f"Starting server on {listen_addr}")

    # サーバーの起動
    await server.start()

    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Server stopping...")
        await server.stop(0)
        logger.info("Server stopped.")

def main() -> None:
    """
    メイン関数：非同期サーバーを起動
    """
    try:
        asyncio.run(serve())
    except Exception as e:
        logger.error(f"Server failed: {e}")
        raise

if __name__ == "__main__":
    main()
