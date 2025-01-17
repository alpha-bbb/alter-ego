package db_functions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

// UpdateConversationChoice は、指定された会話IDのドキュメントに LLM のレスポンスを更新します。
func UpdateConversationChoice(ctx context.Context, coll *mongo.Collection, req *backendpb.SubmitUserChoiceRequest) (string, error) {
	// 更新するフィールドの定義: LLM レスポンス（複数メッセージの場合はスライスとして保存）と更新日時
	update := bson.M{
		"$set": bson.M{
			"choice":    req.Choice,
			"updatedAt": time.Now(),
		},
	}

	// convID をキーにしてドキュメントを検索（ここでは "id" が会話IDと仮定）
	filter := bson.M{"_id": req.ConversationId}

	// UpdateOne でドキュメントを更新する
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	// マッチしたドキュメントが存在しない場合はエラーを返す
	if result.MatchedCount == 0 {
		return "", mongo.ErrNoDocuments
	}

	return req.ConversationId, nil
}
