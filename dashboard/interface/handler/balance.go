package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

type BalanceHandler interface {
	Get() http.HandlerFunc
}

type balanceHandler struct {
	balanceUsecase usecase.BalanceUsecase
}

func NewBalanceHandler(bu usecase.BalanceUsecase) BalanceHandler {
	return &balanceHandler{
		balanceUsecase: bu,
	}
}

func (bh *balanceHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balances, err := bh.balanceUsecase.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		resDto := make([]dto.Balance, 0)
		for i := range balances {
			dto := dto.ConvertBalance(&balances[i])
			if dto != nil {
				resDto = append(resDto, *dto)
			}
		}

		js, err := json.Marshal(resDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
