package controller

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
)

// 相場を分析して取引実行する
func TradeHandler(w http.ResponseWriter, r *http.Request) {
	// 取引履歴
	signalEvents := model.GetSignalEvents(config.DB, config.ProductCode)
	// 見つからなければ終了
	if signalEvents == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get signal_events (productCode=%s)", config.ProductCode)
		return
	}

	// 分析，売買のためのパラメータ
	tradeParams := model.GetTradeParams(config.DB, config.ProductCode)
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
	bot := model.NewTradingBot(config.DB, config.APIKey, config.APISecret, config.ProductCode, config.CandleDuration, 365)
	bot.SignalEvents = signalEvents
	bot.TradeParams = tradeParams

	// 分析，取引
	err := bot.Trade(config.DB, config.CandleTableName, config.TimeFormat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to trade: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade")
}
