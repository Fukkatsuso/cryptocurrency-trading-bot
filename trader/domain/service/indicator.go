package service

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type IndicatorService interface {
	BuySignalOfEMA(emaFast, emaSlow *model.EMA, at int) bool
	SellSignalOfEMA(emaFast, emaSlow *model.EMA, at int) bool
	BuySignalOfBBands(bbands *model.BBands, candles []model.Candle, at int) bool
	SellSignalOfBBands(bbands *model.BBands, candles []model.Candle, at int) bool
	BuySignalOfIchimoku(ichimoku *model.IchimokuCloud, candles []model.Candle, at int) bool
	SellSignalOfIchimoku(ichimoku *model.IchimokuCloud, candles []model.Candle, at int) bool
	BuySignalOfRSI(rsi *model.RSI, buyThread float64, at int) bool
	SellSignalOfRSI(rsi *model.RSI, sellThread float64, at int) bool
	BuySignalOfMACD(macd *model.MACD, at int) bool
	SellSignalOfMACD(macd *model.MACD, at int) bool
}

type indicatorService struct{}

func NewIndicatorService() IndicatorService {
	return &indicatorService{}
}

func (is *indicatorService) BuySignalOfEMA(emaFast, emaSlow *model.EMA, at int) bool {
	if at < emaFast.Period() || at < emaSlow.Period() {
		return false
	}

	// ゴールデンクロス
	goldenCross := emaFast.Values()[at-1] < emaSlow.Values()[at-1] &&
		emaFast.Values()[at] >= emaSlow.Values()[at]

	return goldenCross
}

func (is *indicatorService) SellSignalOfEMA(emaFast, emaSlow *model.EMA, at int) bool {
	if at < emaFast.Period() || at < emaSlow.Period() {
		return false
	}

	// デッドクロス
	deadCross := emaFast.Values()[at-1] > emaSlow.Values()[at-1] &&
		emaFast.Values()[at] <= emaSlow.Values()[at]

	return deadCross
}

func (is *indicatorService) BuySignalOfBBands(bbands *model.BBands, candles []model.Candle, at int) bool {
	if at < bbands.N() {
		return false
	}

	return bbands.Down()[at-1] > candles[at-1].Close() &&
		bbands.Down()[at] <= candles[at].Close()
}

func (is *indicatorService) SellSignalOfBBands(bbands *model.BBands, candles []model.Candle, at int) bool {
	if at < bbands.N() {
		return false
	}

	return bbands.Up()[at-1] < candles[at-1].Close() &&
		bbands.Up()[at] >= candles[at].Close()
}

func (is *indicatorService) BuySignalOfIchimoku(ichimoku *model.IchimokuCloud, candles []model.Candle, at int) bool {
	if at < 1 {
		return false
	}

	// 三役好転
	return ichimoku.Chikou()[at-1] < candles[at-1].High() &&
		ichimoku.Chikou()[at] >= candles[at].High() &&
		ichimoku.SenkouA()[at] < candles[at].Low() &&
		ichimoku.SenkouB()[at] < candles[at].Low() &&
		ichimoku.Tenkan()[at] > ichimoku.Kijun()[at]
}

func (is *indicatorService) SellSignalOfIchimoku(ichimoku *model.IchimokuCloud, candles []model.Candle, at int) bool {
	if at < 1 {
		return false
	}

	// 三役逆転
	return ichimoku.Chikou()[at-1] > candles[at-1].Low() &&
		ichimoku.Chikou()[at] <= candles[at].Low() &&
		ichimoku.SenkouA()[at] > candles[at].High() &&
		ichimoku.SenkouB()[at] > candles[at].High() &&
		ichimoku.Tenkan()[at] < ichimoku.Kijun()[at]
}

func (is *indicatorService) BuySignalOfRSI(rsi *model.RSI, buyThread float64, at int) bool {
	if at < 1 {
		return false
	}

	return rsi.Values()[at-1] < buyThread &&
		rsi.Values()[at] >= buyThread
}

func (is *indicatorService) SellSignalOfRSI(rsi *model.RSI, sellThread float64, at int) bool {
	if at < 1 {
		return false
	}

	return rsi.Values()[at-1] > sellThread &&
		rsi.Values()[at] <= sellThread
}

func (is *indicatorService) BuySignalOfMACD(macd *model.MACD, at int) bool {
	if at < 1 {
		return false
	}

	return macd.Macd()[at] < 0 &&
		macd.MacdSignal()[at] < 0 &&
		macd.Macd()[at-1] < macd.MacdSignal()[at-1] &&
		macd.Macd()[at] >= macd.MacdSignal()[at]
}

func (is *indicatorService) SellSignalOfMACD(macd *model.MACD, at int) bool {
	if at < 1 {
		return false
	}

	return macd.Macd()[at] > 0 &&
		macd.MacdSignal()[at] > 0 &&
		macd.Macd()[at-1] > macd.MacdSignal()[at-1] &&
		macd.Macd()[at] <= macd.MacdSignal()[at]
}
