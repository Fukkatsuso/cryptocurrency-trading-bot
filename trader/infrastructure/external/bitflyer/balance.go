package bitflyer

import (
	"encoding/json"
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func (b *Balance) toDomainModelBalance() *model.Balance {
	return model.NewBalance(b.CurrencyCode, b.Amount, b.Available)
}

type bitflyerBalanceRepository struct {
	apiClient *Client
}

func NewBitFlyerBalanceRepository(apiClient *Client) repository.BalanceRepository {
	return &bitflyerBalanceRepository{
		apiClient: apiClient,
	}
}

func (bbr *bitflyerBalanceRepository) FetchAll() ([]model.Balance, error) {
	path := "me/getbalance"
	resp, err := bbr.apiClient.doRequest("GET", path, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	var balances []Balance
	err = json.Unmarshal(resp, &balances)
	if err != nil {
		return nil, err
	}

	// ドメインモデルに移し替える
	domainModelBalances := make([]model.Balance, len(balances))
	for i, balance := range balances {
		domainModelBalances[i] = *balance.toDomainModelBalance()
	}

	return domainModelBalances, nil
}

func (bbr *bitflyerBalanceRepository) FetchByCurrencyCode(currencyCode string) (*model.Balance, error) {
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
