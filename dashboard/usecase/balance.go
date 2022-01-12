package usecase

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type BalanceUsecase interface {
	Get() ([]model.Balance, error)
}

type balanceUsecase struct {
	balanceRepository repository.BalanceRepository
}

func NewBalanceUsecase(br repository.BalanceRepository) BalanceUsecase {
	return &balanceUsecase{
		balanceRepository: br,
	}
}

func (bu *balanceUsecase) Get() ([]model.Balance, error) {
	return bu.balanceRepository.FetchAll()
}
