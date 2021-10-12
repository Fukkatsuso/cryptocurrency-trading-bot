package bitflyer

import (
	"encoding/json"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
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

func (bor *bitflyerOrderRepository) Send(order model.Order) (string, error) {
	data, err := json.Marshal(order)
	if err != nil {
		return "", err
	}

	url := "me/sendchildorder"
	resp, err := bor.apiClient.doRequest("POST", url, map[string]string{}, data)
	if err != nil {
		return "", err
	}

	var response ResponseSendChildOrder
	if err = json.Unmarshal(resp, &response); err != nil {
		return "", err
	}

	return response.ChildOrderAcceptanceID, nil
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
