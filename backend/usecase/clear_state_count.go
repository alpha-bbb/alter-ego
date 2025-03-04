package usecase

import (
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/entity"
)

type IClearStateCountUseCase interface {
	Execute() ([]entity.StateCount, error)
}

type ClearStateCountUseCase struct {
	stateCountRepository repository.IStateCountRepository
}

func NewClearStateCountUseCase(stateCountRepository repository.IStateCountRepository) IClearStateCountUseCase {
	return &ClearStateCountUseCase{stateCountRepository: stateCountRepository}
}

func (c ClearStateCountUseCase) Execute() ([]entity.StateCount, error) {
	statusCounts, err := c.stateCountRepository.GetAll()
	if err != nil {
		return []entity.StateCount{}, err
	}
	return statusCounts, c.stateCountRepository.Clear()
}
