package main

import (
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"

	"github.com/alpha-bbb/alter-ego/backend/adapter/database/postgresql"
	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
	"github.com/joho/godotenv"
)

// システムの利用回数を記録しているテーブルをリセットする
func main() {
	logger, _ := log.NewLogger()
	// .envファイルを読み込み、環境変数をセットする
	if err := godotenv.Load("../../.env.cmd"); err != nil {
		logger.Info(fmt.Sprintf("Warning: .env file could not be loaded: %v", err))
	}

	db, err := postgresql.NewPostgreSQL(config.DSN())
	if err != nil {
		panic(err)
	}
	stateCountRepository := repository.NewStateCountRepository(db)
	clearStateCountUseCse := usecase.NewClearStateCountUseCase(stateCountRepository)
	stateCounts, err := clearStateCountUseCse.Execute()
	if err != nil {
		panic(err)
	}
	logger.Info("Cleared state count")
	for _, stateCount := range stateCounts {
		logger.Info(fmt.Sprintf("Cleared state count: %v", stateCount))
	}
}
