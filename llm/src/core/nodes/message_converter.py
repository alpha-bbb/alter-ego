from typing import Any
from langchain_core.runnables import RunnableConfig
from src.core.model import GraphState, Role, TalkHistory, MarkdownTalkHistories

async def convert_history(
    state: GraphState,
    config: RunnableConfig,
) -> dict[str, Any]:
    """
    トーク履歴をmarkdown形式に変換するノード

    Args:
        state (GraphState): グラフの現在の状態
        config (RunnableConfig): ランナブルの設定

    Returns:
        dict[str, Any]: 更新された状態
    """
    try:
        # 自分のトーク履歴とその他のトーク履歴を分離
        my_histories = [
            history for history in state.talk_histories
            if history.user.role == Role.SELF
        ]
        other_histories = [
            history for history in state.talk_histories
            if history.user.role == Role.YOU
        ]

        # markdown形式に変換
        markdown_histories = {
            "conversations": _to_markdown(state.talk_histories),
            "my_histories": _to_markdown(my_histories),
            "other_histories": _to_markdown(other_histories)
        }

        return {
            "talk_histories_markdown" :MarkdownTalkHistories(**markdown_histories)
        }

    except Exception as e:
        # エラーハンドリング
        print(f"Error in convert_history: {str(e)}")
        return {
            "error": str(e)
        }

def _to_markdown(histories: list[TalkHistory]) -> str:
    """
    履歴をmarkdown形式に変換する補助関数

    Args:
        histories (list): 変換する履歴のリスト

    Returns:
        str: markdown形式の文字列
    """
    markdown_list = []
    for history in histories:
        markdown = f"""
- **Date**: {history.date}
- **Message**: {history.message}
- **User**:
  - **Name**: {history.user.name}
  - **Role**: {history.user.role}
        """.strip()
        markdown_list.append(markdown)

    return "\n\n".join(markdown_list)