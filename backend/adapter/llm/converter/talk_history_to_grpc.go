package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/entity"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
)

func ConvertTalkHistoryToGRPCTalkRequest(entities []entity.TalkHistory) *llmpb.TalkRequest {
	histories := make([]*llmpb.TalkHistory, len(entities))
	for i := range histories {
		histories[i] = &llmpb.TalkHistory{
			Date: entities[i].Date,
			User: &llmpb.User{
				UserId: entities[i].User.UserID,
				Name:   entities[i].User.Name,
				Role: llmpb.User_UserRole(
					entities[i].User.Role),
			},
			Message: entities[i].Message,
		}
	}
	return &llmpb.TalkRequest{Histories: histories}
}
