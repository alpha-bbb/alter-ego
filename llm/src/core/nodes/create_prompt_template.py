from src.core.model import GraphState

def create_prompt_template(state: GraphState) -> dict[str, str]:
    """プロンプトテンプレートを作成"""

    # 話し方の特徴を文字列化
    style = state.template_val.speaking_style
    style_text = """
### 私の話し方

#### 1. Dialect
- **type**: {dialect_type}
- **examples**: {dialect_examples}

#### 2. Tone
- **type**: {tone_type}
- **description**: {tone_description}
- **examples**: {tone_examples}

...
    """.format(
        dialect_type=style.dialect.type if style else "",
        dialect_examples=", ".join(style.dialect.examples) if style else "",
        tone_type=style.tone.type if style else "",
        tone_description=style.tone.description if style else "",
        tone_examples=", ".join(style.tone.examples) if style else "",
    ) if style else ""

    # 知識ベースの文字列化
    knowledge_text = "\n".join([
        f"- {k.content}" for k in state.template_val.knowledge
    ]) if state.template_val.knowledge else ""

    # 最終的なプロンプト
    template = f"""
あなたはコミュニケーションの達人であり、心理学に基づいた、気まずさを最小限に抑える断り方の専門家です。

以下の情報を活用し、ユーザーの会話相手が自然に感じる形で会話を終了する一言を考えてください：

1. 相手の心理:
- 行動: {state.template_val.action_analysis.predicted_action if state.template_val.action_analysis else ""}
- その行動への自信: {state.template_val.action_analysis.confidence_score if state.template_val.action_analysis else ""}
- 詳細: {state.template_val.action_analysis.summary if state.template_val.action_analysis else ""}

2. 断り方のナレッジ:
{knowledge_text}

3. 話し方の特徴:
{style_text}

### 回答ルール
- 会話状況や話し方の特徴を反映した自然な表現を使い、一言で会話を終了してください。
- 上記情報（会話状況や話し方の特徴）については直接言及しないでください。
- ユーザーの話し方を取り入れて、より自然で親しみやすい言葉遣いを意識してください。
- 断り方の知識を参考にしつつ、相手の心理を考慮して総合的に判断してください。

### 出力形式
- 会話を終わらせる一言のみを出力してください。
"""

    return {
        "prompt_template": template
    }