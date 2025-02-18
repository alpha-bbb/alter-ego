package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func ConvertMessageDTOToGRPCTalkResponse(messageDTO dto.MessageDTO) *backendpb.TalkResponse {
	return &backendpb.TalkResponse{
		Message: messageDTO.Messages,
		Status:  backendpb.TalkResponse_TalkResponseStatus(messageDTO.Status),
	}
}
