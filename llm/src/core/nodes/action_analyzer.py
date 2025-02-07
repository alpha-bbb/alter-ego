from typing import Any, Dict

from langchain_core.messages import HumanMessage, SystemMessage
from langchain_core.runnables import RunnableConfig
from langchain_openai import ChatOpenAI

from src.core.model import ActionAnalysis, GraphState


async def analyze_action(
    state: GraphState,
    config: RunnableConfig,
) -> Dict[str, Any]:
    """
    相手の言動を分析するLLMノード

    Args:
        state (GraphState): グラフの現在の状態
        config (RunnableConfig): ランナブルの設定

    Returns:
        Dict[str, Any]: 更新された状態（action_analysis を含む）
    """
    try:
        # LLMの初期化
        llm = ChatOpenAI(model="gpt-4o-mini").with_structured_output(ActionAnalysis)

        # システムプロンプトの定義
        system_prompt = """
        1. 与えられた会話を注意深く読み、相手が取ろうとしている行動や意図を分析してください。
        2. roleが1の場合は自分を表します。
        3. 行動の予測には、以下の選択肢から最も適切なものを選んでください：
           - お誘い
           - 世間話
           - 愚痴・悩み相談
           - 近況報告
           - お願い
           - その他
        4. 予測した行動に対する自信の度合いを0から1の範囲でスコアとして評価してください。
        5. 相手の発言の主なポイントを分析して下さい。
        """

        # LLMに分析を依頼
        messages = [
            SystemMessage(content=system_prompt),
            HumanMessage(content=state.talk_histories_markdown.conversations),
        ]

        analysis = await llm.ainvoke(messages, config)

        return {"template_val": {"action_analysis": analysis}}

    except Exception as e:
        print(f"Error in analyze_action: {str(e)}")
        return {"error": str(e)}
