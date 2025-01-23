package db_functions

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect は MongoDB に接続し、クライアントと指定した DB、Collection を返します。
func Connect() (*mongo.Client, *mongo.Database, *mongo.Collection, error) {
	// 環境変数の取得
	host, err := getEnv("MONGO_HOST")
	if err != nil {
		return nil, nil, nil, err
	}

	port, err := getEnv("MONGO_PORT")
	if err != nil {
		return nil, nil, nil, err
	}

	dbName, err := getEnv("MONGO_DB")
	if err != nil {
		return nil, nil, nil, err
	}

	user, err := getEnv("MONGO_USER")
	if err != nil {
		return nil, nil, nil, err
	}

	password, err := getEnv("MONGO_PASSWORD")
	if err != nil {
		return nil, nil, nil, err
	}

	// MongoDB URI の構築
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", user, password, host, port, dbName)

	// MongoDB に接続
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 利用する DB、Collection の指定
	db := client.Database(dbName)
	coll := db.Collection("conversations")

	return client, db, coll, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable '%s' is missing or empty", key)
	}
	return value, nil
}
