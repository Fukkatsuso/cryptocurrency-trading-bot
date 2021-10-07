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
		if ds.indicatorService.BuySignalOfRSI(rsi, buyThread, sellThread, i) {
			signal := model.NewSignalEvent(candle.Time().Time(), df.ProductCode(), model.OrderSideBuy, candle.Close(), size)
			if signal != nil {
				signalEvents.AddBuySignal(*signal)
			}
		}

		if ds.indicatorService.SellSignalOfRSI(rsi, buyThread, sellThread, i) {
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
