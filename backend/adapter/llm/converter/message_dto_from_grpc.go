package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
)

func ConvertMessageDTOFromGRPCTalkResponse(res *llmpb.TalkResponse) dto.MessageDTO {
	return dto.MessageDTO{
		Messages: res.Message,
		Status:   1, // XXX:暫定
	}
}
