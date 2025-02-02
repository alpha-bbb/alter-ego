from typing import Any
from langchain_core.messages import SystemMessage, HumanMessage
from langchain_openai import ChatOpenAI
from langchain_core.runnables import RunnableConfig
from src.core.model import GraphState, SpeakingStyle

async def analyze_style(
    state: GraphState,
    config: RunnableConfig,
) -> dict[str, Any]:
    """
    話し方を解析するLLMノード

    Args:
        state (GraphState): グラフの現在の状態
        config (RunnableConfig): ランナブルの設定

    Returns:
        dict[str, Any]: 更新された状態
    """
    try:
        # LLMの初期化
        llm = ChatOpenAI(model="gpt-4o-mini").with_structured_output(SpeakingStyle)

        # システムプロンプトの定義
        system_prompt = """
        会話のやり取りの履歴を受け取って、私の話し方の特徴を以下の点から分析してください：

        1. 方言や独自の言葉遣いの特徴
        2. 話の雰囲気や態度
        3. 文の長さや構造の特徴
        4. 絵文字やスタンプの使用頻度と具体例
        5. 冗談や笑いの取り入れ方
        6. 共感や気遣いを示す表現

        分析結果は以下の形式のJSONで出力してください。
        """

        # 自分の会話履歴の取得
        my_histories = state.talk_histories_markdown.my_histories

        if not my_histories:
            return {"speaking_style": None}

        # LLMに分析を依頼
        messages = [
            SystemMessage(content=system_prompt),
            HumanMessage(content=my_histories)
        ]

        speaking_style = await llm.ainvoke(messages, config)

        return {
            "template_val": {
                "speaking_style": speaking_style
            },
        }

    except Exception as e:
        print(f"Error in analyze_style: {str(e)}")
        return {
            "error": str(e)
        }