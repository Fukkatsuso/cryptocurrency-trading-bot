package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
)

// 相場を分析して取引実行する
func TradeHandler(w http.ResponseWriter, r *http.Request) {
	// 分析，売買のためのパラメータ
	tradeParams := model.GetTradeParams(config.DB, config.TradeParamTableName, config.ProductCode)
	fmt.Println("params:", tradeParams)
	// パラメータが見つからなければ終了
	if tradeParams == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "trade_params has no param record (productCode=%s)", config.ProductCode)
		return
	}
	// 取引無効になっていたら終了
	if !tradeParams.TradeEnable {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "trade is not enabled (productCode=%s)", config.ProductCode)
		return
	}

	// 取引bot
	_ = model.NewTradingBot(config.ProductCode, 24*time.Hour, 365)

	// 分析，取引

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade")
}
