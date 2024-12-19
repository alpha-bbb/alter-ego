package convert

import (
	"github.com/alpha-bbb/alter-ego/backend/entity"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
)

func ConvertTalkHistoryToGRPCTalkResponse(histories []*entity.TalkHistory) []*llmpb.TalkHistory {
    if histories == nil {
        return nil
    }

    result := make([]*llmpb.TalkHistory, len(histories))
    for i := range histories {
        result[i] = &llmpb.TalkHistory{
            Date: histories[i].Date,
            User: &llmpb.User{
                UserId: histories[i].User.UserID,
                Name:   histories[i].User.Name,
            },
            Message: histories[i].Message,
        }
    }
    return result
}
