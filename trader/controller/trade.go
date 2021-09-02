package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/slack"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
)

// 相場を分析して取引実行する
func TradeHandler(w http.ResponseWriter, r *http.Request) {
	// 取引履歴
	signalEvents := model.GetSignalEvents(config.DB, config.ProductCode)
	// 見つからなければ終了
	if signalEvents == nil {
		slackMsg := slack.BuildTextMessage(
			fmt.Sprintf("%s（%s）", slack.SlackEmojiDizzyFace, config.ProductCode),
			"SignalEventsが取得できません",
		)
		err := slack.PostTextMessage(config.SlackBotToken, config.SlackChannelID, slackMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get signal_events (productCode=%s)", config.ProductCode)
		return
	}

	// 分析，売買のためのパラメータ
	tradeParams := model.GetTradeParams(config.DB, config.ProductCode)
	fmt.Println("params:", tradeParams)
	// パラメータが見つからなければ終了
	if tradeParams == nil {
		slackMsg := slack.BuildTextMessage(
			fmt.Sprintf("%s（%s）", slack.SlackEmojiDizzyFace, config.ProductCode),
			"TradeParamsが取得できません",
		)
		err := slack.PostTextMessage(config.SlackBotToken, config.SlackChannelID, slackMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "trade_params has no param record (productCode=%s)", config.ProductCode)
		return
	}
	// 取引無効になっていたら終了
	if !tradeParams.TradeEnable {
		slackMsg := slack.BuildTextMessage(
			fmt.Sprintf("%s（%s）", slack.SlackEmojiDizzyFace, config.ProductCode),
			"取引が無効に設定されています",
		)
		err := slack.PostTextMessage(config.SlackBotToken, config.SlackChannelID, slackMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
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
		slackMsg := slack.BuildTextMessage(
			fmt.Sprintf("%s（%s）", slack.SlackEmojiDizzyFace, config.ProductCode),
			"取引時にエラーが生じました",
			"```",
			err.Error(),
			"```",
		)
		err := slack.PostTextMessage(config.SlackBotToken, config.SlackChannelID, slackMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to trade: %s", err.Error())
		return
	}

	// slack通知
	for _, signal := range bot.SignalEvents.Signals {
		if signal.Time.After(beforeTradeTime) {
			slackMsg := SignalEventToSlackTextMessage(&signal)
			err := slack.PostTextMessage(config.SlackBotToken, config.SlackChannelID, slackMsg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade")
}

func SignalEventToSlackTextMessage(signal *model.SignalEvent) string {
	timeString := signal.Time.In(config.LocalTime).Format(config.TimeFormat)
	msg := slack.BuildTextMessage(
		fmt.Sprintf("%s *%s*: %s", slack.SlackEmojiCoin, signal.Side, signal.ProductCode),
		fmt.Sprintf("At: %s", timeString),
		fmt.Sprintf("Price: %.3f", signal.Price),
		fmt.Sprintf("Size: %.3f", signal.Size),
	)
	return msg
}
