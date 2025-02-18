package usecase

import (
	"fmt"

	converterLlm "github.com/alpha-bbb/alter-ego/backend/adapter/llm/converter"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/client"
)

type ITalkUseCase interface {
	Execute(entities []entity.TalkHistory) (dto.MessageDTO, error)
}

type TalkUseCase struct {
	llmClient client.LLMClient
}

func NewTalkUseCase(llmClient client.LLMClient) *TalkUseCase {
	return &TalkUseCase{llmClient: llmClient}
}

func (u *TalkUseCase) Execute(entities []entity.TalkHistory) (dto.MessageDTO, error) {
	talkRequest := converterLlm.ConvertTalkHistoryToGRPCTalkRequest(entities)
	llmResponse, err := u.llmClient.Talk(talkRequest)
	if err != nil {
		return dto.MessageDTO{}, fmt.Errorf("failed to call LLM service: %w", err)
	}
	return converterLlm.ConvertMessageDTOFromGRPCTalkResponse(llmResponse), nil
}
