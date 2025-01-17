package db_models

import (
	"time"
)

// User は会話に参加するユーザ（または LLM）の情報
type User struct {
	Name   string `bson:"name"`
	UserID string `bson:"userId"`
	Role   int    `bson:"role"`
}

// TalkHistory は会話の1エントリを表現する
type TalkHistory struct {
	Date    time.Time `bson:"date"`
	User    User      `bson:"user"`
	Message string    `bson:"message"`
}

// LlmResponses は候補レスポンスと、後からユーザが選択した choice を保持する
type LlmResponses struct {
	Candidates []string `bson:"candidates,omitempty"`
}

// Conversation は1つの会話セッションを表現するドキュメント
type Conversation struct {
	ID          string        `bson:"_id"`
	Histories   []TalkHistory `bson:"histories"`
	LlmResponse LlmResponses  `bson:"llmResponses"`
	Choice      string        `bson:"choice"`
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
}
