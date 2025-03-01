package converter

import (
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
)

func FromGRPCTalkResponse(res *llmpb.TalkResponse) dto.Message {
	return dto.Message{
		Messages: res.Message,
		Status:   1, // TODO: ステータス変換ロジック
	}
}
