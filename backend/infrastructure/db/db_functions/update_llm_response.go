package db_functions

import (
	"context"
	"time"

	"github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateConversationChoice は、指定された会話IDのドキュメントに LLM のレスポンスを更新します。
func UpdateLlmResponse(ctx context.Context, coll *mongo.Collection, convID string, messages []string) (string, error) {
	update := bson.M{
		"$set": bson.M{
			"llmResponses": db_models.LlmResponses{
				Candidates: messages,
			},
			"updatedAt": time.Now(),
		},
	}

	// convID をキーにしてドキュメントを検索（ここでは "id" が会話IDと仮定）
	filter := bson.M{"_id": convID}

	// UpdateOne でドキュメントを更新する
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	// マッチしたドキュメントが存在しない場合はエラーを返す
	if result.MatchedCount == 0 {
		return "", mongo.ErrNoDocuments
	}

	return convID, nil
}
