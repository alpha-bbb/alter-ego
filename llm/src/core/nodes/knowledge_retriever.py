from typing import Any
from langchain_core.runnables import RunnableConfig
from src.core.model import GraphState, Knowledge
from src.services.external_api.dify import call_dify_retrieve

async def retrieve_knowledge(
    state: GraphState,
    config: RunnableConfig,
) -> dict[str, Any]:
    """
    断り方のナレッジを取得するノード

    Args:
        state (GraphState): グラフの現在の状態
        config (RunnableConfig): ランナブルの設定
        dataset_ids (list[str], optional): 検索対象のデータセットID

    Returns:
        dict[str, Any]: 更新された状態
    """
    try:
        # 行動分析結果の取得
        if not state.template_val.action_analysis:
            return {"knowledge": []}

        # 検索クエリの作成
        query = f"""
        行動タイプ: {state.template_val.action_analysis.predicted_action}
        要約: {state.template_val.action_analysis.summary}
        """

        # ダミーの断り方ナレッジ（実際にはデータベースから取得）
        # documents = [
        #     Document(
        #         page_content="丁寧に感謝を示しながら断る: まずは誘ってくれたことへの感謝を示し、その後に断る理由を簡潔に説明します。",
        #         metadata={"type": "お誘い", "score": 0.9}
        #     ),
        #     Document(
        #         page_content="代替案を提示する: 今回は参加できないが、別の機会や方法を提案することで、相手への配慮を示します。",
        #         metadata={"type": "お誘い", "score": 0.8}
        #     ),
        #     Document(
        #         page_content="明確に断る: あいまいな表現を避け、はっきりと意思を伝えることで、誤解を防ぎます。",
        #         metadata={"type": "お誘い", "score": 0.7}
        #     )
        # ]

        # 関連ドキュメントの検索
        retrieved_docs = await call_dify_retrieve(dataset_id="ea577e2f-7fc8-408a-bd8c-2e8cc254ea37", query=query, retrieval_model={
            "search_method": "hybrid_search",
            "reranking_enable": True,
            "reranking_mode":"reranking_model",
            "reranking_model": {
                "reranking_provider_name": "cohere",
                "reranking_model_name": "rerank-multilingual-v3.0"
            },
            "weights": None,
            "top_k": 3,
            "score_threshold_enabled": False,
            "score_threshold": None
        })

        return {
            "template_val": {
                "knowledge": [
                    Knowledge(content=doc)
                    for doc in retrieved_docs
                ]
            }
        }

    except Exception as e:
        print(f"Error in retrieve_knowledge: {str(e)}")
        return {
            "error": str(e)
        }