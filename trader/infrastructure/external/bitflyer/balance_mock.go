package bitflyer

import (
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type bitflyerBalanceMockRepository struct {
	apiClient *Client
}

func NewBitFlyerBalanceMockRepository(apiClient *Client) repository.BalanceRepository {
	return &bitflyerBalanceRepository{
		apiClient: apiClient,
	}
}

func (bbr *bitflyerBalanceMockRepository) FetchAll() ([]model.Balance, error) {
	balances := []Balance{
		{
			CurrencyCode: "JPY",
			Amount:       10000,
			Available:    10000,
		},
		{
			CurrencyCode: "BTC",
			Amount:       0,
			Available:    0,
		},
		{
			CurrencyCode: "ETH",
			Amount:       0.01,
			Available:    0.01,
		},
		{
			CurrencyCode: "XRP",
			Amount:       0,
			Available:    0,
		},
	}

	// ドメインモデルに移し替える
	domainModelBalances := make([]model.Balance, len(balances))
	for i, balance := range balances {
		domainModelBalances[i] = *balance.toDomainModelBalance()
	}

	return domainModelBalances, nil
}

func (bbr *bitflyerBalanceMockRepository) FetchByCurrencyCode(currencyCode string) (*model.Balance, error) {
	balances, err := bbr.FetchAll()
	if err != nil {
		return nil, errors.New("cannot fetch balance")
	}

	for _, balance := range balances {
		if balance.CurrencyCode() == currencyCode {
			return &balance, nil
		}
	}

	return nil, errors.New("invalid currencyCode")
}
