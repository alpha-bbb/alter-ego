package db_functions

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect は MongoDB に接続し、クライアントと指定した DB、Collection を返します。
func Connect() (*mongo.Client, *mongo.Database, *mongo.Collection, error) {
	// 環境変数から .env ファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// MONGODB_URI の取得
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, nil, nil, ErrMissingMongoURI()
	}

	// MongoDB に接続
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, nil, err
	}

	// 利用する DB、Collection の指定
	db := client.Database("sample_conversations")
	coll := db.Collection("conversations")

	return client, db, coll, nil
}

// ErrMissingMongoURI は MONGODB_URI が設定されていない場合のエラーを返します。
func ErrMissingMongoURI() error {
	return &MissingEnvError{Msg: "You must set your 'MONGODB_URI' environment variable."}
}

// MissingEnvError は環境変数が見つからない場合のエラー型です。
type MissingEnvError struct {
	Msg string
}

func (e *MissingEnvError) Error() string {
	return e.Msg
}
