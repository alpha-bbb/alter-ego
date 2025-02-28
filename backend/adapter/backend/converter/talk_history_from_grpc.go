package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func TalkHistoriesFromGRPC(req *backendpb.TalkRequest) (dto.Talk, error) {
	if req == nil ||
		req.Histories == nil ||
		req.Account == nil {
		return dto.Talk{}, ErrInvalidRequest
	}

	histories := make([]dto.TalkHistory, len(req.Histories))
	for i, history := range req.Histories {
		histories[i] = dto.TalkHistory{
			Date: history.Date,
			User: dto.User{
				UserID: history.User.UserId,
				Name:   history.User.Name,
				Role:   int(history.User.Role),
			},
			Message: history.Message,
		}
	}
	return dto.Talk{
		Histories: histories,
		Account: dto.Account{
			PlatformType: dto.PlatformType(req.Account.PlatformType),
			AccountID:    req.Account.AccountId,
		},
	}, nil
}
