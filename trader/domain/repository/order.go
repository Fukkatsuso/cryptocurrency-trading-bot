package repository

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type OrderRepository interface {
	Send(order model.Order) (*model.Order, error)
	FetchById(productCode, orderId string) ([]model.Order, error)
}
