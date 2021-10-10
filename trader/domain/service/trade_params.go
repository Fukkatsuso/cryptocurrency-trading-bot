package service

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type TradeParamsService interface {
	OptimizeEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) (float64, int, int)
	OptimizeBBands(df *model.DataFrame, n int, k float64, size float64) (float64, int, float64)
	OptimizeIchimoku(df *model.DataFrame, size float64) float64
	OptimizeRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) (float64, int, float64, float64)
	OptimizeMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) (float64, int, int, int)

	OptimizeAll(df *model.DataFrame, params *model.TradeParams) *model.TradeParams
}

type tradeParamsService struct {
	dataFrameService DataFrameService
}

func NewTradeParamsService(ds DataFrameService) TradeParamsService {
	return &tradeParamsService{
		dataFrameService: ds,
	}
}

func (ts *tradeParamsService) OptimizeEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) (float64, int, int) {
	performance := float64(0)
	bestFastPeriod := fastPeriod
	bestSlowPeriod := slowPeriod

	for fastPeriod = 5; fastPeriod < 11; fastPeriod++ {
		for slowPeriod = 12; slowPeriod < 20; slowPeriod++ {
			signalEvents := ts.dataFrameService.BacktestEMA(df, fastPeriod, slowPeriod, size)
			if signalEvents == nil {
				continue
			}
			profit := signalEvents.EstimateProfit()
			if performance < profit {
				performance = profit
				bestFastPeriod = fastPeriod
				bestSlowPeriod = slowPeriod
			}
		}
	}

	return performance, bestFastPeriod, bestSlowPeriod
}

func (ts *tradeParamsService) OptimizeBBands(df *model.DataFrame, n int, k float64, size float64) (float64, int, float64) {
	performance := float64(0)
	bestN := n
	bestK := k

	for n := 10; n <= 30; n++ {
		for k := 1.8; k <= 2.2; k += 0.1 {
			signalEvents := ts.dataFrameService.BacktestBBands(df, n, k, size)
			if signalEvents == nil {
				continue
			}
			profit := signalEvents.EstimateProfit()
			if performance < profit {
				performance = profit
				bestN = n
				bestK = k
			}
		}
	}

	return performance, bestN, bestK
}

func (ts *tradeParamsService) OptimizeIchimoku(df *model.DataFrame, size float64) float64 {
	signalEvents := ts.dataFrameService.BacktestIchimoku(df, size)
	if signalEvents == nil {
		return 0
	}
	performance := signalEvents.EstimateProfit()

	return performance
}

func (ts *tradeParamsService) OptimizeRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) (float64, int, float64, float64) {
	performance := float64(0)
	bestPeriod := period
	bestBuyThread, bestSellThread := buyThread, sellThread

	for period := 3; period < 30; period++ {
		for buyThread := float64(20); buyThread <= 40; buyThread++ {
			for sellThread := float64(60); sellThread <= 80; sellThread++ {
				signalEvents := ts.dataFrameService.BacktestRSI(df, period, buyThread, sellThread, size)
				if signalEvents == nil {
					continue
				}
				profit := signalEvents.EstimateProfit()
				if performance < profit {
					performance = profit
					bestPeriod = period
					bestBuyThread = buyThread
					bestSellThread = sellThread
				}
			}
		}
	}

	return performance, bestPeriod, bestBuyThread, bestSellThread
}

func (ts *tradeParamsService) OptimizeMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) (float64, int, int, int) {
	performance := float64(0)
	bestFastPeriod := fastPeriod
	bestSlowPeriod := slowPeriod
	bestSignalPeriod := signalPeriod

	for fastPeriod := 5; fastPeriod < 20; fastPeriod++ {
		for slowPeriod := 20; slowPeriod < 40; slowPeriod++ {
			for signalPeriod := 5; signalPeriod < 15; signalPeriod++ {
				signalEvents := ts.dataFrameService.BacktestMACD(df, fastPeriod, slowPeriod, signalPeriod, size)
				if signalEvents == nil {
					continue
				}
				profit := signalEvents.EstimateProfit()
				if performance < profit {
					performance = profit
					bestFastPeriod = fastPeriod
					bestSlowPeriod = slowPeriod
					bestSignalPeriod = signalPeriod
				}
			}
		}
	}

	return performance, bestFastPeriod, bestSlowPeriod, bestSignalPeriod
}

func (ts *tradeParamsService) OptimizeAll(df *model.DataFrame, params *model.TradeParams) *model.TradeParams {
	_, emaPeriod1, emaPeriod2 := ts.OptimizeEMA(df, params.EMAPeriod1(), params.EMAPeriod2(), params.Size())
	_, bbandsN, bbandsK := ts.OptimizeBBands(df, params.BBandsN(), params.BBandsK(), params.Size())
	_, rsiPeriod, rsiBuyThread, rsiSellThread := ts.OptimizeRSI(df, params.RSIPeriod(), params.RSIBuyThread(), params.RSISellThread(), params.Size())
	_, macdFastPeriod, macdSlowPeriod, macdSignalPeriod := ts.OptimizeMACD(df, params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSignalPeriod(), params.Size())

	newParams := model.NewTradeParams(
		params.TradeEnable(),
		params.ProductCode(),
		params.Size(),
		params.SMAEnable(),
		params.SMAPeriod1(),
		params.SMAPeriod2(),
		params.SMAPeriod3(),
		params.EMAEnable(),
		emaPeriod1,
		emaPeriod2,
		params.EMAPeriod3(),
		params.BBandsEnable(),
		bbandsN,
		bbandsK,
		params.IchimokuEnable(),
		params.RSIEnable(),
		rsiPeriod,
		rsiBuyThread,
		rsiSellThread,
		params.MACDEnable(),
		macdFastPeriod,
		macdSlowPeriod,
		macdSignalPeriod,
		params.StopLimitPercent(),
	)
	return newParams
}
