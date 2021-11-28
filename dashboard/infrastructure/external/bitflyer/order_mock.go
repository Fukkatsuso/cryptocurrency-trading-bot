package bitflyer

import (
	"math/rand"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type bitflyerOrderMockRepository struct {
}

func NewBitflyerOrderMockRepository() repository.OrderRepository {
	return &bitflyerOrderMockRepository{}
}

func (bor *bitflyerOrderMockRepository) Send(order model.Order) (*model.Order, error) {
	rand.Seed(time.Now().UnixNano())
	price := 200000 + float64(rand.Intn(300000))

	completedOrder := &model.Order{
		ProductCode:     order.ProductCode,
		ChildOrderType:  order.ChildOrderType,
		Side:            order.Side,
		AveragePrice:    price,
		Size:            order.Size,
		MinuteToExpires: order.MinuteToExpires,
		TimeInForce:     order.TimeInForce,
		ChildOrderState: model.OrderState(OrderStateCompleted),
		ChildOrderDate:  time.Now().Format(TimestampFormat),
		TotalCommission: order.Size * 0.0015,
	}

	return completedOrder, nil
}
