package controller

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
)

// 相場を分析して取引実行する
func TradeHandler(w http.ResponseWriter, r *http.Request) {
	// 分析，売買のためのパラメータ
	_ = model.TradeParams{
		ProductCode: config.ProductCode,
		Size:        0.01,
	}

	// 取引実行するクライアント
	_ = bitflyer.NewClient(config.APIKey, config.APISecret)

	// 分析，取引

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade")
}
