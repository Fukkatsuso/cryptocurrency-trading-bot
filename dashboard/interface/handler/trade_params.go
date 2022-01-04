package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
	params, err := reqFormToTradeParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func reqFormToTradeParams(r *http.Request) (*model.TradeParams, error) {
	tradeEnable := r.FormValue("trade") == "true"

	productCode := r.FormValue("productCode")

	size, err := strconv.ParseFloat(r.FormValue("size"), 64)
	if err != nil {
		return nil, err
	}

	smaEnable := r.FormValue("sma") == "true"
	smaPeriod1, err := strconv.Atoi(r.FormValue("smaPeriod1"))
	if err != nil {
		return nil, err
	}
	smaPeriod2, err := strconv.Atoi(r.FormValue("smaPeriod2"))
	if err != nil {
		return nil, err
	}
	smaPeriod3, err := strconv.Atoi(r.FormValue("smaPeriod3"))
	if err != nil {
		return nil, err
	}

	emaEnable := r.FormValue("ema") == "true"
	emaPeriod1, err := strconv.Atoi(r.FormValue("emaPeriod1"))
	if err != nil {
		return nil, err
	}
	emaPeriod2, err := strconv.Atoi(r.FormValue("emaPeriod2"))
	if err != nil {
		return nil, err
	}
	emaPeriod3, err := strconv.Atoi(r.FormValue("emaPeriod3"))
	if err != nil {
		return nil, err
	}

	bbandsEnable := r.FormValue("bbands") == "true"
	bbandsN, err := strconv.Atoi(r.FormValue("bbandsN"))
	if err != nil {
		return nil, err
	}
	bbandsK, err := strconv.ParseFloat(r.FormValue("bbandsK"), 64)
	if err != nil {
		return nil, err
	}

	ichimokuEnable := r.FormValue("ichimoku") == "true"

	rsiEnable := r.FormValue("rsi") == "true"
	rsiPeriod, err := strconv.Atoi(r.FormValue("rsiPeriod"))
	if err != nil {
		return nil, err
	}
	rsiBuyThread, err := strconv.ParseFloat(r.FormValue("rsiBuyThread"), 64)
	if err != nil {
		return nil, err
	}
	rsiSellThread, err := strconv.ParseFloat(r.FormValue("rsiSellThread"), 64)
	if err != nil {
		return nil, err
	}

	macdEnable := r.FormValue("macd") == "true"
	macdFastPeriod, err := strconv.Atoi(r.FormValue("macdFastPeriod"))
	if err != nil {
		return nil, err
	}
	macdSlowPeriod, err := strconv.Atoi(r.FormValue("macdSlowPeriod"))
	if err != nil {
		return nil, err
	}
	macdSignalPeriod, err := strconv.Atoi(r.FormValue("macdSignalPeriod"))
	if err != nil {
		return nil, err
	}

	stopLimitPercent, err := strconv.ParseFloat(r.FormValue("stopLimitPercent"), 64)
	if err != nil {
		return nil, err
	}

	params := model.NewTradeParams(
		tradeEnable,
		productCode,
		size,
		smaEnable,
		smaPeriod1,
		smaPeriod2,
		smaPeriod3,
		emaEnable,
		emaPeriod1,
		emaPeriod2,
		emaPeriod3,
		bbandsEnable,
		bbandsN,
		bbandsK,
		ichimokuEnable,
		rsiEnable,
		rsiPeriod,
		rsiBuyThread,
		rsiSellThread,
		macdEnable,
		macdFastPeriod,
		macdSlowPeriod,
		macdSignalPeriod,
		stopLimitPercent,
	)

	if params == nil {
		return nil, errors.New("invalid parameter")
	}
	return params, nil
}
