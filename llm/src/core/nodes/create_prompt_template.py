from src.core.model import GraphState


def create_prompt_template(state: GraphState) -> dict[str, str]:
    """動的な断り方のプロンプトテンプレートを作成"""

    # 話し方の特徴や知識はそのまま前回と同様に作成
    style = state.template_val.speaking_style
    style_text = (
        """
### 私の話し方

#### 1. 方言・独自の言葉遣い
- **タイプ**: {dialect_type}
- **具体例**: {dialect_examples}

#### 2. 話し方の雰囲気・態度
- **タイプ**: {tone_type}
- **詳細**: {tone_description}
- **具体例**: {tone_examples}

#### 3. 文の構造
- **特徴**: {sentence_structure_type}
- **具体例**: {sentence_structure_examples}

#### 4. 絵文字・スタンプの使用
- **使用頻度**: {emoji_usage_type}
- **具体例**: {emoji_usage_examples}

#### 5. 冗談や笑いの取り入れ方
- **特徴**: {humor_type}
- **具体例**: {humor_examples}

#### 6. 共感や気遣いの表現
- **特徴**: {expressions_of_care_type}
- **具体例**: {expressions_of_care_examples}
    """.format(
            dialect_type=style.dialect.type if style else "",
            dialect_examples=", ".join(style.dialect.examples) if style else "",
            tone_type=style.tone.type if style else "",
            tone_description=style.tone.description if style else "",
            tone_examples=", ".join(style.tone.examples) if style else "",
            sentence_structure_type=style.sentence_structure.type if style else "",
            sentence_structure_examples=", ".join(style.sentence_structure.examples)
            if style
            else "",
            emoji_usage_type=style.emoji_usage.type if style else "",
            emoji_usage_examples=", ".join(style.emoji_usage.examples) if style else "",
            humor_type=style.humor.type if style else "",
            humor_examples=", ".join(style.humor.examples) if style else "",
            expressions_of_care_type=style.expressions_of_care.type if style else "",
            expressions_of_care_examples=", ".join(style.expressions_of_care.examples)
            if style
            else "",
        )
        if style
        else ""
    )

    knowledge_text = (
        "\n".join([f"- {k.content}" for k in state.template_val.knowledge])
        if state.template_val.knowledge
        else ""
    )

    # 以下のルールにより、断りのトーンを動的に変更する:
    # ・相手の「行動への自信（confidence_score）」が低い場合: 徐々に会話を終了する、柔らかい表現
    # ・相手の「行動への自信」が高い場合: はっきりと断る、強めの意思表示
    template = f"""
あなたはコミュニケーションの達人であり、心理学に基づいた断り方の専門家です。

以下の情報を活用し、相手の心理状態に合わせて会話を終了する一言を生成してください：

1. 相手の心理:
- **行動**: {state.template_val.action_analysis.predicted_action if state.template_val.action_analysis else ""}
- **行動への自信**: {state.template_val.action_analysis.confidence_score if state.template_val.action_analysis else ""}
- **詳細**: {state.template_val.action_analysis.summary if state.template_val.action_analysis else ""}

2. 断り方の知識:
{knowledge_text}

3. 話し方の特徴:
{style_text}

### 動的な対応ルール
- **行動への自信が低い場合**（例: 自信度が低い・曖昧な行動の場合）は、会話を徐々に終了させるよう、柔らかく穏やかな表現で「そろそろ失礼します」などといったフェードアウトを意識してください。
- **行動への自信が高い場合**（例: 明確な行動や強い主張が見られる場合）は、はっきりと断る意思を示す強い表現で「申し訳ありませんが、これ以上は難しいです」といった断固たる返答を意識してください。

また、以下の断りの技法も適切なものを内部で活用してください：
- **サンドイッチ話法**: 感謝や敬意を先に示し、断りのメッセージを挟み、最後に前向きな締めの言葉を添える。
- **ポライトネス理論**: 相手のフェイス（面子）に配慮した、丁寧で敬意ある表現を取り入れる。
- **ダブルバインドの応用**: 断りながらも、代替案や前向きな選択肢をほのめかす表現を検討する。

### 回答ルール
- 上記情報（会話状況や話し方の特徴）に基づき、一言で私がする回答を生成してください。
- 必ず、相手の最後の一言をみて、その内容に合わせた自然な返信を行ってください。

### 出力形式
- 続く私の回答を一言のみを出力してください。
"""

    return {"prompt_template": template}
