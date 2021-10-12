package repository

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type OrderRepository interface {
	Send(order model.Order) (string, error)
	FetchById(productCode, orderId string) ([]model.Order, error)
}
