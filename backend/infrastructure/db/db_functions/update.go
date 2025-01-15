package db_functions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateConversationChoice は指定した会話ドキュメントに、ユーザの選択肢情報を更新する
func UpdateConversationChoice(ctx context.Context, coll *mongo.Collection, convID interface{}, choice string) (int64, error) {
	choiceDate := time.Now()

	filter := bson.D{{Key: "_id", Value: convID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "llmResponses.choice", Value: choice},
		{Key: "llmResponses.choiceDate", Value: choiceDate},
		{Key: "updatedAt", Value: time.Now()},
	}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
