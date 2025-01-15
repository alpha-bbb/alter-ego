package db_functions

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Conversation struct {
	// define the fields of the Conversation struct
	ID interface{} `bson:"_id"`
	// add other fields as needed
}

// FindConversationByID は指定した会話IDのドキュメントを取得する
func FindConversationByID(ctx context.Context, coll *mongo.Collection, convID interface{}) (Conversation, error) {
	var conv Conversation

	filter := bson.D{{Key: "_id", Value: convID}}
	err := coll.FindOne(ctx, filter).Decode(&conv)
	return conv, err
}
