from typing import Any
from langchain_core.messages import SystemMessage, HumanMessage
from langchain_core.runnables import RunnableConfig
from langchain_openai import ChatOpenAI
from langchain_anthropic import ChatAnthropic
from langchain_google_genai import ChatGoogleGenerativeAI

from src.core.model import GraphState

async def generate_responses(
    state: GraphState,
    config: RunnableConfig,
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
        gpt4 = ChatOpenAI(model="gpt-4o-mini")
        claude = ChatAnthropic(model_name="claude-3-5-sonnet-20240620", timeout=60, stop=["Human:"])
        gemini = ChatGoogleGenerativeAI(model="gemini-1.5-pro")

        # プロンプトテンプレートの作成
        template = state.prompt_template

        # 各LLMで並行して生成
        tasks = [
            _generate_single_response(gpt4, template, config),
            _generate_single_response(claude, template, config),
            _generate_single_response(gemini, template, config)
        ]

        # 並行実行
        import asyncio
        response = await asyncio.gather(*tasks)

        return {
            "final_responses" : response
        }

    except Exception as e:
        print(f"Error in generate_responses: {str(e)}")
        return {
            "error": str(e)
        }

async def _generate_single_response(
    llm: Any,
    template: str,
    config: RunnableConfig
) -> str:
    """単一のLLMで応答を生成"""

    messages = [
        SystemMessage(content=template),
        HumanMessage(content="会話を終了させる一言を生成してください。")
    ]

    response = await llm.ainvoke(messages, config)
    return response.content