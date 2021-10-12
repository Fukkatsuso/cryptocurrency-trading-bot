package model

type ChildOrderType string

const (
	ChildOrderTypeLimit  ChildOrderType = "LIMIT"  // 指値注文
	ChildOrderTypeMarket ChildOrderType = "MARKET" // 成行注文
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"  // 買い注文
	OrderSideSell OrderSide = "SELL" // 売り注文
)

// 執行数量条件
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // default
	TimeInForceIOC TimeInForce = "IOC"
	TimeInForceFOK TimeInForce = "FOK"
)

type OrderState string

const (
	OrderStateActive    OrderState = "ACTIVE"    // オープンな注文
	OrderStateCompleted OrderState = "COMPLETED" // 全額が取引完了した注文
	OrderStateCanceled  OrderState = "CANCELED"  // キャンセルした注文
	OrderStateExpired   OrderState = "EXPIRED"   // 有効期限に到達したため取り消された注文
	OrderStateRejected  OrderState = "REJECTED"  // 失敗した注文
)

// 新規注文リクエスト
// 注文一覧レスポンス
type Order struct {
	ID                     int            `json:"id"`
	ProductCode            string         `json:"product_code"`
	ChildOrderType         ChildOrderType `json:"child_order_type"`
	Side                   OrderSide      `json:"side"`
	Price                  float64        `json:"price"`
	AveragePrice           float64        `json:"average_price"`
	Size                   float64        `json:"size"`
	MinuteToExpires        int            `json:"minute_to_expire"`
	TimeInForce            TimeInForce    `json:"time_in_force"`
	ChildOrderID           string         `json:"child_order_id"`
	ChildOrderState        OrderState     `json:"child_order_state"`
	ExpireDate             string         `json:"expire_date"`
	ChildOrderDate         string         `json:"child_order_date"`
	ChildOrderAcceptanceID string         `json:"child_order_acceptance_id"`
	OutstandingSize        float64        `json:"outstanding_size"`
	CancelSize             float64        `json:"cancel_size"`
	ExecutedSize           float64        `json:"executed_size"`
	TotalCommission        float64        `json:"total_commission"`
}
