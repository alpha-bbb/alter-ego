package db_convert

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 会話履歴を変換する関数
func convertToTalkHistories(entityTalkHistory []*entity.TalkHistory, logger *zap.Logger) []db_models.TalkHistory {
	histories := make([]db_models.TalkHistory, len(entityTalkHistory))
	for i, histEntity := range entityTalkHistory {
		histories[i] = db_models.TalkHistory{
			Date: func() time.Time {
				parsedTime, err := time.Parse(time.RFC3339, histEntity.Date)
				if err != nil {
					logger.Error("Failed to parse date", zap.Error(err))
					return time.Time{}
				}
				return parsedTime
			}(),
			User: db_models.User{
				UserID: histEntity.User.UserID,
				Name:   histEntity.User.Name,
				Role:   histEntity.User.Role,
			},
			Message: histEntity.Message,
		}
	}
	return histories
}

// 会話履歴を作成する関数
func ConvertEntityToModel(entityTalkHistory []*entity.TalkHistory, logger *zap.Logger) *db_models.Conversation {
	histories := convertToTalkHistories(entityTalkHistory, logger)
	now := time.Now()

	newConv := &db_models.Conversation{
		ID:          uuid.New().String(), // Generate a unique ID
		Histories:   histories,
		LlmResponse: db_models.LlmResponses{},
		Choice:      "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return newConv
}
