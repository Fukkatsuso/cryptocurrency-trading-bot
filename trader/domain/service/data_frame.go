package service

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
)

type DataFrameService interface {
	BacktestEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) *model.SignalEvents
	BacktestBBands(df *model.DataFrame, n int, k float64, size float64) *model.SignalEvents
	BacktestIchimoku(df *model.DataFrame, size float64) *model.SignalEvents
	BacktestRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) *model.SignalEvents
	BacktestMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) *model.SignalEvents

	Backtest(df *model.DataFrame, tp *model.TradeParams)
	Analyze(df *model.DataFrame, at int, params *model.TradeParams) (bool, bool)
}

type dataFrameService struct {
	indicatorService IndicatorService
}

func NewDataFrameService(is IndicatorService) DataFrameService {
	return &dataFrameService{
		indicatorService: is,
	}
}

func (ds *dataFrameService) BacktestEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) *model.SignalEvents {
	emaFast := model.NewEMA(df.Closes(), fastPeriod)
	if emaFast == nil {
		return nil
	}
	emaSlow := model.NewEMA(df.Closes(), slowPeriod)
	if emaSlow == nil {
		return nil
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		if ds.indicatorService.BuySignalOfEMA(emaFast, emaSlow, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfEMA(emaFast, emaSlow, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), size)
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	return signalEvents
}

func (ds *dataFrameService) BacktestBBands(df *model.DataFrame, n int, k float64, size float64) *model.SignalEvents {
	bbands := model.NewBBands(df.Closes(), n, k)
	if bbands == nil {
		return nil
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		if ds.indicatorService.BuySignalOfBBands(bbands, df.Candles(), i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfBBands(bbands, df.Candles(), i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), size)
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	return signalEvents
}

func (ds *dataFrameService) BacktestIchimoku(df *model.DataFrame, size float64) *model.SignalEvents {
	ichimoku := model.NewIchimokuCloud(df.Closes())
	if ichimoku == nil {
		return nil
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		if ds.indicatorService.BuySignalOfIchimoku(ichimoku, df.Candles(), i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfIchimoku(ichimoku, df.Candles(), i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), size)
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	return signalEvents
}

func (ds *dataFrameService) BacktestRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) *model.SignalEvents {
	rsi := model.NewRSI(df.Closes(), period)
	if rsi == nil {
		return nil
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		if ds.indicatorService.BuySignalOfRSI(rsi, buyThread, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfRSI(rsi, sellThread, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), size)
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	return signalEvents
}

func (ds *dataFrameService) BacktestMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) *model.SignalEvents {
	macd := model.NewMACD(df.Closes(), fastPeriod, slowPeriod, signalPeriod)
	if macd == nil {
		return nil
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		if ds.indicatorService.BuySignalOfMACD(macd, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfMACD(macd, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), size)
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	return signalEvents
}

func (ds *dataFrameService) Backtest(df *model.DataFrame, params *model.TradeParams) {
	if df == nil || params == nil {
		return
	}

	signals := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(signals)
	for i, candle := range df.Candles() {
		buy, sell := ds.Analyze(df, i, params)

		if buy {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), params.Size())
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if sell ||
			signalEvents.ShouldCutLoss(candle.Close(), params.StopLimitPercent()) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideSell, candle.Close(), params.Size())
			if signal != nil {
				signalEvents.AddSellSignal(*signal)
			}
		}
	}

	signalEvents.EstimateProfit()

	df.AddBacktestEvents(signalEvents)
}

// 各指標の時点"at"で分析する
// buyPoint, sellPointを返す
func (ds *dataFrameService) Analyze(df *model.DataFrame, at int, params *model.TradeParams) (bool, bool) {
	buyPoint, sellPoint := 0, 0

	if at <= 0 {
		return false, false
	}

	if params.EMAEnable() &&
		len(df.EMAs()) >= 2 {
		emaFast := df.EMAs()[0]
		emaSlow := df.EMAs()[1]
		if ds.indicatorService.BuySignalOfEMA(&emaFast, &emaSlow, at) {
			buyPoint++
		}
		if ds.indicatorService.SellSignalOfEMA(&emaFast, &emaSlow, at) {
			sellPoint++
		}
	}

	if params.BBandsEnable() {
		bbands := df.BBands()
		if ds.indicatorService.BuySignalOfBBands(bbands, df.Candles(), at) {
			buyPoint++
		}
		if ds.indicatorService.SellSignalOfBBands(bbands, df.Candles(), at) {
			sellPoint++
		}
	}

	if params.IchimokuEnable() {
		ichomoku := df.IchimokuCloud()
		if ds.indicatorService.BuySignalOfIchimoku(ichomoku, df.Candles(), at) {
			buyPoint++
		}
		if ds.indicatorService.SellSignalOfIchimoku(ichomoku, df.Candles(), at) {
			sellPoint++
		}
	}

	if params.RSIEnable() {
		rsi := df.RSI()
		if ds.indicatorService.BuySignalOfRSI(rsi, params.RSIBuyThread(), at) {
			buyPoint++
		}
		if ds.indicatorService.SellSignalOfRSI(rsi, params.RSISellThread(), at) {
			sellPoint++
		}
	}

	if params.MACDEnable() {
		macd := df.MACD()
		if ds.indicatorService.BuySignalOfMACD(macd, at) {
			buyPoint++
		}
		if ds.indicatorService.SellSignalOfMACD(macd, at) {
			sellPoint++
		}
	}

	return buyPoint > 1, sellPoint > 1
}
