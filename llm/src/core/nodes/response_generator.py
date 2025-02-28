from typing import Any

from langchain_google_genai import ChatGoogleGenerativeAI

from src.core.model import GraphState, LLMResponse


async def generate_responses(
    state: GraphState,
) -> dict[str, Any]:
    """
    複数のLLMを使用して返答を生成するノード

    Args:
        state (GraphState): グラフの現在の状態
        config (RunnableConfig): ランナブルの設定

    Returns:
        dict[str, Any]: 更新された状態
    """
    try:
        # LLMsの初期化
        gemini = ChatGoogleGenerativeAI(model="gemini-1.5-pro").with_structured_output(
            LLMResponse
        )

        # プロンプトテンプレートの作成
        template = state.prompt_template

        response = gemini.invoke(template)

        return {"final_responses": response.model_dump()}

    except Exception as e:
        print(f"Error in generate_responses: {str(e)}")
        return {"error": str(e)}
