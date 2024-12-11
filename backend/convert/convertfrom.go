package convert

import (
	"github.com/alpha-bbb/alter-ego/backend/entity"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func ConvertTalkHistoryFromGRPCTalkRequest(req *backendpb.TalkRequest) []*entity.TalkHistory {
    if req == nil || req.Histories == nil {
        return nil
    }

    result := make([]*entity.TalkHistory, len(req.Histories))
    for index := range req.Histories {
        result[index] = &entity.TalkHistory{
            Date: req.Histories[index].Date,
            User: entity.User{
                UserID: req.Histories[index].User.UserId,
                Name:   req.Histories[index].User.Name,
            },
            Message: req.Histories[index].Message,
        }
    }
    println("🟥", result)
    return result
}