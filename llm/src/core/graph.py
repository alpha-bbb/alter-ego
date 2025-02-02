from typing import Literal
from langgraph.graph import StateGraph, END

from src.core.model import GraphState
from src.core.nodes.message_converter import convert_history
from src.core.nodes.action_analyzer import analyze_action
from src.core.nodes.knowledge_retriever import retrieve_knowledge
from src.core.nodes.style_analyzer import analyze_style
from src.core.nodes.response_generator import generate_responses
from src.core.nodes.create_prompt_template import create_prompt_template
from src.core.nodes.join_template import join_template
class AlterEgo:
    """Alter-goのメイングラフクラス"""

    def __init__(self):
        workflow = StateGraph(GraphState)

        workflow.add_node("convert_history", convert_history)
        workflow.add_node("analyze_action", analyze_action)
        workflow.add_node("retrieve_knowledge", retrieve_knowledge)
        workflow.add_node("analyze_style", analyze_style)

        workflow.add_node("join_template", join_template)

        workflow.add_node("generate_template", create_prompt_template)
        workflow.add_node("generate_responses", generate_responses)

        workflow.set_entry_point("convert_history")
        workflow.add_edge("convert_history", "analyze_action")
        workflow.add_conditional_edges("convert_history", self._route_after_conversion)
        workflow.add_conditional_edges("analyze_action", self._route_after_analysis)


        workflow.add_edge("analyze_style", "join_template")
        workflow.add_edge("retrieve_knowledge", "join_template")

        workflow.add_conditional_edges("join_template", self._route_after_join)

        workflow.add_edge("generate_template", "generate_responses")
        workflow.add_edge("generate_responses", END)

        self.graph = workflow.compile()

    def _route_after_conversion(self, state: GraphState) -> Literal["analyze_style", "join_template"]:
        """会話履歴変換後のルーティング"""
        if state.talk_histories:
            return "analyze_style"
        return "join_template"

    def _route_after_analysis(
        self,
        state: GraphState
    ) -> Literal["retrieve_knowledge", "join_template"]:
        """行動分析後のルーティング"""
        if state.template_val.action_analysis.predicted_action == "お誘い":
            return "retrieve_knowledge"

        return "join_template"

    def _route_after_join(self, state: GraphState) -> Literal["generate_template", "join_template"]:
        """
        state.template_val が存在すれば generate_template へ、
        存在しなければ join_template を再実行するようにルーティングする。
        """
        if state.template_ready:
            return "generate_template"
        return "join_template"

# グラフのインスタンスを作成
graph = AlterEgo().graph