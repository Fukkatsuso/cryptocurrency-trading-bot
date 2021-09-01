package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
	"github.com/slack-go/slack"
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

	// 取引前の時刻
	// 通知発生基準にする
	beforeTradeTime := time.Now().UTC()

	// 分析，取引
	err := bot.Trade(config.DB, config.CandleTableName, config.TimeFormat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to trade: %s", err.Error())
		return
	}

	// slack通知
	for _, signal := range bot.SignalEvents.Signals {
		if signal.Time.After(beforeTradeTime) {
			slackMsg := SignalEventToSlackTextMessage(&signal)
			err := PostSlackTextMessage(slackMsg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade")
}

func PostSlackTextMessage(msg string) error {
	slackBot := slack.New(config.SlackBotToken)
	_, _, err := slackBot.PostMessage(config.SlackChannelID, slack.MsgOptionText(msg, true))
	return err
}

func SignalEventToSlackTextMessage(signal *model.SignalEvent) string {
	msg := fmt.Sprintf(":coin: *%s*: %s\nAt: %s\nPrice: %f\nSize: %f",
		signal.Side,
		signal.ProductCode,
		signal.Time.In(config.LocalTime).Format(config.TimeFormat),
		signal.Price,
		signal.Size,
	)
	return msg
}
