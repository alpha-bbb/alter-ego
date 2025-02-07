import operator
from datetime import datetime
from enum import Enum
from typing import Annotated, Optional, Union

from pydantic import BaseModel, Field, field_validator


class Role(str, Enum):
    SELF = "SELF"
    YOU = "YOU"
    OTHER = "OTHER"


class User(BaseModel):
    name: str = Field(..., description="ユーザー名")
    role: Role = Field(..., description="ユーザーの役割 (1: 自分, 2: 相手)")


class TalkHistory(BaseModel):
    date: datetime = Field(..., description="会話の日時")
    message: str = Field(..., description="メッセージ内容")
    user: User = Field(..., description="ユーザー情報")


class MarkdownTalkHistories(BaseModel):
    conversations: str = Field("", description="会話履歴のmarkdown形式")
    my_histories: str = Field("", description="自分の会話履歴のmarkdown形式")
    other_histories: str = Field("", description="相手の会話履歴のmarkdown形式")


class ActionType(str, Enum):
    INVITATION = "お誘い"
    CHAT = "世間話"
    CONSULTATION = "愚痴・悩み相談"
    STATUS_REPORT = "近況報告"
    REQUEST = "お願い"
    OTHER = "その他"


class ActionAnalysis(BaseModel):
    predicted_action: ActionType = Field(..., description="予測される行動タイプ")
    confidence_score: float = Field(..., description="予測の信頼度スコア (0-1)")
    summary: str = Field(..., description="会話の要約")


class Knowledge(BaseModel):
    content: str = Field(..., description="断り方に関する知識コンテンツ")


class ExampleSet(BaseModel):
    type: str = Field(..., description="タイプまたは頻度")
    examples: list[str] = Field(..., description="具体例のリスト")
    description: Optional[str] = Field(None, description="詳細な説明（必要な場合）")


class SpeakingStyle(BaseModel):
    dialect: ExampleSet = Field(..., description="方言や独自の言葉遣いの特徴")
    tone: ExampleSet = Field(..., description="話の雰囲気や態度")
    sentence_structure: ExampleSet = Field(..., description="文の長さや構造の特徴")
    emoji_usage: ExampleSet = Field(
        ..., description="絵文字やスタンプの使用頻度と具体例"
    )
    humor: ExampleSet = Field(..., description="冗談や笑いの取り入れ方")
    expressions_of_care: ExampleSet = Field(..., description="共感や気遣いを示す表現")


class TemplateVal(BaseModel):
    action_analysis: Optional[ActionAnalysis] = Field(None, description="行動分析結果")
    knowledge: Optional[list[Knowledge]] = Field(
        default_factory=list, description="参照された知識"
    )
    speaking_style: Optional[SpeakingStyle] = Field(
        None, description="話し方の分析結果"
    )

    @field_validator("action_analysis")
    def validate_confidence_score(cls, v):
        if v and (v.confidence_score < 0 or v.confidence_score > 1):
            raise ValueError("信頼度スコアは0から1の間である必要があります")
        return v

    def __add__(self, other: Union["TemplateVal", dict]) -> "TemplateVal":
        # other が dict の場合は TemplateVal に変換する
        if isinstance(other, dict):
            try:
                other = TemplateVal.parse_obj(other)
            except Exception as e:
                raise ValueError(
                    "渡された dict を TemplateVal に変換できませんでした"
                ) from e

        if not isinstance(other, TemplateVal):
            return NotImplemented

        # action_analysis は、どちらか片方が None であれば非 None の値を、
        # 両方ともある場合は self の値を採用（必要に応じて別ロジックに変更可能）
        new_action_analysis = (
            self.action_analysis
            if self.action_analysis is not None
            else other.action_analysis
        )

        # knowledge は、両方のリストを連結する
        new_knowledge = (self.knowledge or []) + (other.knowledge or [])

        # speaking_style も、どちらか片方が None であれば非 None の値を、
        # 両方ともある場合は self の値を採用
        new_speaking_style = (
            self.speaking_style
            if self.speaking_style is not None
            else other.speaking_style
        )

        return TemplateVal(
            action_analysis=new_action_analysis,
            knowledge=new_knowledge,
            speaking_style=new_speaking_style,
        )

    # __radd__ を定義することで、dict + TemplateVal も可能にする
    def __radd__(self, other: Union["TemplateVal", dict]) -> "TemplateVal":
        return self.__add__(other)


class GraphState(BaseModel):
    """グラフの状態を管理するメインのStateクラス"""

    talk_histories: list[TalkHistory] = Field(
        default_factory=list, description="会話履歴"
    )
    talk_histories_markdown: MarkdownTalkHistories = Field(
        MarkdownTalkHistories(conversations="", my_histories="", other_histories=""),
        description="会話履歴のmarkdown形式",
    )
    prompt_template: Annotated[str, operator.add] = Field(
        "", description="プロンプトテンプレート"
    )
    template_val: Annotated[TemplateVal, operator.add] = Field(
        TemplateVal(action_analysis=None, knowledge=None, speaking_style=None),
        description="プロンプトテンプレートの変数",
    )
    template_ready: bool = Field(
        False, description="プロンプトテンプレートが完成しているかどうか"
    )
    final_responses: list[str] = Field(
        default_factory=list, description="生成された返答のリスト"
    )

    class Config:
        use_enum_values = True
        json_schema_extra = {
            "example": {
                "messages": ["こんにちは"],
                "talk_histories": [
                    {
                        "date": "2024-02-02T12:00:00",
                        "message": "今度の週末、カラオケに行かない？",
                        "user": {"name": "田中", "role": 2},
                    }
                ],
                "action_analysis": {
                    "predicted_action": "お誘い",
                    "confidence_score": 0.95,
                    "summary": "週末のカラオケへの誘い",
                },
            }
        }
