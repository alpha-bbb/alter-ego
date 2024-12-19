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
    for i := range req.Histories {
        result[i] = &entity.TalkHistory{
            Date: req.Histories[i].Date,
            User: entity.User{
                UserID: req.Histories[i].User.UserId,
                Name:   req.Histories[i].User.Name,
            },
            Message: req.Histories[i].Message,
        }
    }
    return result
}
