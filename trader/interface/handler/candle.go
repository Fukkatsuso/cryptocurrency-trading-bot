package handler

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

type CandleHandler interface {
	UpdateCandle(productCode string) http.HandlerFunc
}

type candleHandler struct {
	candleUsecase usecase.CandleUsecase
}

func NewCandleHandler(cu usecase.CandleUsecase) CandleHandler {
	return &candleHandler{
		candleUsecase: cu,
	}
}

func (ch *candleHandler) UpdateCandle(productCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ch.candleUsecase.UpdateCandle(productCode)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to update candle")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success")
	}
}
