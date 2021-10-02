package model

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"  // 買い注文
	OrderSideSell OrderSide = "SELL" // 売り注文
)
