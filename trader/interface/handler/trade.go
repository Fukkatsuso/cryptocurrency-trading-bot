package handler

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

type TradeHandler interface {
	Trade(productCode string, pastPeriod int) http.HandlerFunc
}

type tradeHandler struct {
	tradeUsecase usecase.TradeUsecase
}

func NewTradeHandler(tu usecase.TradeUsecase) TradeHandler {
	return &tradeHandler{
		tradeUsecase: tu,
	}
}

func (th *tradeHandler) Trade(productCode string, pastPeriod int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := th.tradeUsecase.Trade(productCode, pastPeriod)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to trade")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success")
	}
}
