package bitflyer

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"  // 買い注文
	OrderSideSell OrderSide = "SELL" // 売り注文
)

type OrderState string

const (
	OrderStateActive    OrderState = "ACTIVE"    // オープンな注文
	OrderStateCompleted OrderState = "COMPLETED" // 全額が取引完了した注文
	OrderStateCanceled  OrderState = "CANCELED"  // キャンセルした注文
	OrderStateExpired   OrderState = "EXPIRED"   // 有効期限に到達したため取り消された注文
	OrderStateRejected  OrderState = "REJECTED"  // 失敗した注文
)

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

type bitflyerOrderRepository struct {
	apiClient *Client
}

func NewBitflyerOrderRepository(apiClient *Client) repository.OrderRepository {
	return &bitflyerOrderRepository{
		apiClient: apiClient,
	}
}

func (bor *bitflyerOrderRepository) Send(order model.Order) (*model.Order, error) {
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	url := "me/sendchildorder"
	resp, err := bor.apiClient.doRequest("POST", url, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	var response ResponseSendChildOrder
	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	childOrderAcceptanceId := response.ChildOrderAcceptanceID
	completedOrder := bor.waitUntilOrderComplete(order.ProductCode, childOrderAcceptanceId)
	if completedOrder == nil {
		return nil, errors.New("order is not completed")
	}

	return completedOrder, nil
}

func (bor *bitflyerOrderRepository) waitUntilOrderComplete(productCode, orderId string) *model.Order {
	// 最長2分待つ
	expire := time.After(2 * time.Minute)
	// 15秒ごとに注文状況をポーリング
	interval := time.Tick(15 * time.Second)

	return func() *model.Order {
		for {
			select {
			case <-expire:
				return nil
			case <-interval:
				orders, err := bor.FetchById(productCode, orderId)
				if err != nil {
					return nil
				}
				if len(orders) == 0 {
					return nil
				}
				order := orders[0]
				if order.ChildOrderState == model.OrderState(OrderStateCompleted) {
					if order.Side == model.OrderSide(OrderSideBuy) {
						return &order
					}
					if order.Side == model.OrderSide(OrderSideSell) {
						return &order
					}
					return nil
				}
			}
		}
	}()
}

func (bor *bitflyerOrderRepository) FetchById(productCode, orderId string) ([]model.Order, error) {
	query := map[string]string{
		"product_code":              productCode,
		"child_order_acceptance_id": orderId,
	}

	resp, err := bor.apiClient.doRequest("GET", "me/getchildorders", query, nil)
	if err != nil {
		return nil, err
	}

	var responseListOrder []model.Order
	if err = json.Unmarshal(resp, &responseListOrder); err != nil {
		return nil, err
	}
	return responseListOrder, nil
}
