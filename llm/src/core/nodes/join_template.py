from src.core.model import GraphState

def join_template(state: GraphState) -> GraphState:
    """
    analyze_style と retrieve_knowledge の結果が揃っているかを確認し、
    両方が揃っていれば TemplateVal を生成して state.template_val にセットし、
    state.template_ready を True にする。
    """

    if state.template_val.knowledge and state.template_val.speaking_style and state.template_val.action_analysis:
        state.template_ready = True
    else:
        state.template_ready = False

    return state