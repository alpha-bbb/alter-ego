package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
)

func ToGRPCTalkRequest(histories []dto.TalkHistory) *llmpb.TalkRequest {
	talkHistories := make([]*llmpb.TalkHistory, len(histories))
	for i, h := range histories {
		talkHistories[i] = &llmpb.TalkHistory{
			Date: h.Date,
			User: &llmpb.User{
				UserId: h.User.UserID,
				Name:   h.User.Name,
				Role:   llmpb.User_UserRole(h.User.Role),
			},
			Message: h.Message,
		}
	}
	return &llmpb.TalkRequest{Histories: talkHistories}
}
