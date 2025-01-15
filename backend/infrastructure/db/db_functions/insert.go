package db_functions

import (
	"context"
	"log"
	"time"

	"github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_models"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertConversation は新規会話データ（候補レスポンス付き）を挿入する
func InsertConversation(ctx context.Context, coll *mongo.Collection, conv db_models.Conversation) (interface{}, error) {
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = conv.CreatedAt

	result, err := coll.InsertOne(ctx, conv)
	if err != nil {
		log.Printf("InsertConversation error: %v", err)
		return nil, err
	}
	return result.InsertedID, nil
}
