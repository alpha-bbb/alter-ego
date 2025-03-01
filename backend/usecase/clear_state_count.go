package usecase

import "github.com/alpha-bbb/alter-ego/backend/database/repository"

type IClearStateCountUseCase interface {
	Execute() error
}

type ClearStateCountUseCase struct {
	stateCountRepository repository.IStateCountRepository
}

func NewClearStateCountUseCase(stateCountRepository repository.IStateCountRepository) IClearStateCountUseCase {
	return &ClearStateCountUseCase{stateCountRepository: stateCountRepository}
}

func (c ClearStateCountUseCase) Execute() error {
	return c.stateCountRepository.Clear()
}
