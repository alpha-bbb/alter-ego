package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func FromGRPCSubscribeRequest(req *backendpb.SubscribeRequest) (dto.Subscribe, error) {
	if req == nil ||
		req.Account == nil {
		return dto.Subscribe{}, ErrInvalidRequest
	}
	return dto.Subscribe{
		Account: dto.Account{
			PlatformType: dto.PlatformType(req.Account.PlatformType),
			AccountID:    req.Account.AccountId,
		},
		Action: dto.SubscribeAction(req.Action),
	}, nil
}
