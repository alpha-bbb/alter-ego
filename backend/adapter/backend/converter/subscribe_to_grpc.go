package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func ToGRPCSubscribeResponse(res dto.SubscribeInfo) *backendpb.SubscribeResponse {
	return &backendpb.SubscribeResponse{
		Status:      backendpb.SubscribeResponse_SubscribeStatus(res.Status),
		Message:     res.Message,
		RedirectUrl: res.RedirectUrl,
		ExpireAt:    res.ExpiresAt,
	}
}
