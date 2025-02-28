package main

import (
	"github.com/alpha-bbb/alter-ego/backend/adapter/database/postgresql"
	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
)

// システムの利用回数を記録しているテーブルをリセットする
func main() {
	db, err := postgresql.NewPostgreSQL(config.DSN())
	if err != nil {
		panic(err)
	}
	stateCountRepository := repository.NewStateCountRepository(db)
	clearStateCountUseCse := usecase.NewClearStateCountUseCase(stateCountRepository)
	err = clearStateCountUseCse.Execute()
	if err != nil {
		panic(err)
	}
}
