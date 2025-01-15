package db_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User は会話に参加するユーザ（または LLM）の情報
type User struct {
	Name   string `bson:"name"`
	UserID string `bson:"userId"`
	Role   string `bson:"role"`
}

// TalkHistory は会話の1エントリを表現する
type TalkHistory struct {
	Date    time.Time `bson:"date"`
	User    User      `bson:"user"`
	Message string    `bson:"message"`
}

// LlmResponseCandidate は LLM 層が返す候補レスポンス
type LlmResponseCandidate struct {
	Option  string `bson:"option"`  // 候補番号（例："1", "2", "3"）
	Message string `bson:"message"` // レスポンス内容
}

// LlmResponses は候補レスポンスと、後からユーザが選択した choice を保持する
type LlmResponses struct {
	Candidates []LlmResponseCandidate `bson:"candidates,omitempty"`
	Choice     *string                `bson:"choice,omitempty"`     // ユーザが選択した候補番号。初回は nil
	ChoiceDate *time.Time             `bson:"choiceDate,omitempty"` // 選択日時（任意）
}

// Conversation は1つの会話セッションを表現するドキュメント
type Conversation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userId"`
	ActionKind  int                `bson:"actionKind"`
	Histories   []TalkHistory      `bson:"histories"`
	LlmResponse LlmResponses       `bson:"llmResponses"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}
