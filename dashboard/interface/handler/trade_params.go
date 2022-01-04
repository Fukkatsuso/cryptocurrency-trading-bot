package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

type TradeParamsHandler interface {
	HandlerFunc() http.HandlerFunc
}

type tradeParamsHandler struct {
	tradeParamsUsecase usecase.TradeParamsUsecase
}

func NewTradeParamsHandler(tu usecase.TradeParamsUsecase) TradeParamsHandler {
	return &tradeParamsHandler{
		tradeParamsUsecase: tu,
	}
}

func (th *tradeParamsHandler) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			th.Get(w, r)
		case http.MethodPost:
			th.Post(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (th *tradeParamsHandler) Get(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("productCode")

	params, err := th.tradeParamsUsecase.Get(productCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dto := dto.ConvertTradeParams(params)

	js, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (th *tradeParamsHandler) Post(w http.ResponseWriter, r *http.Request) {
	params, err := reqJsonToTradeParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = th.tradeParamsUsecase.Save(*params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func reqJsonToTradeParams(r *http.Request) (*model.TradeParams, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var dto dto.TradeParams
	if err := json.Unmarshal(body, &dto); err != nil {
		return nil, err
	}

	params := model.NewTradeParams(
		dto.TradeEnable,
		dto.ProductCode,
		dto.Size,
		dto.SMAEnable,
		dto.SMAPeriod1,
		dto.SMAPeriod2,
		dto.SMAPeriod3,
		dto.EMAEnable,
		dto.EMAPeriod1,
		dto.EMAPeriod2,
		dto.EMAPeriod3,
		dto.BBandsEnable,
		dto.BBandsN,
		dto.BBandsK,
		dto.IchimokuEnable,
		dto.RSIEnable,
		dto.RSIPeriod,
		dto.RSIBuyThread,
		dto.RSISellThread,
		dto.MACDEnable,
		dto.MACDFastPeriod,
		dto.MACDSlowPeriod,
		dto.MACDSignalPeriod,
		dto.StopLimitPercent,
	)

	if params == nil {
		return nil, errors.New("invalid parameter")
	}
	return params, nil
}
