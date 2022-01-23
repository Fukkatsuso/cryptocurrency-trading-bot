package service

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type TradeParamsService interface {
	Save(params model.TradeParams) error
	Find(productCode string) (*model.TradeParams, error)

	OptimizeEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) (float64, int, int, bool)
	OptimizeBBands(df *model.DataFrame, n int, k float64, size float64) (float64, int, float64, bool)
	OptimizeIchimoku(df *model.DataFrame, size float64) (float64, bool)
	OptimizeRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) (float64, int, float64, float64, bool)
	OptimizeMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) (float64, int, int, int, bool)

	OptimizeAll(df *model.DataFrame, params *model.TradeParams) (*model.TradeParams, bool)
}

type tradeParamsService struct {
	tradeParamsRepository repository.TradeParamsRepository
	dataFrameService      DataFrameService
}

func NewTradeParamsService(ts repository.TradeParamsRepository, ds DataFrameService) TradeParamsService {
	return &tradeParamsService{
		tradeParamsRepository: ts,
		dataFrameService:      ds,
	}
}

func (ts *tradeParamsService) Save(params model.TradeParams) error {
	return ts.tradeParamsRepository.Save(params)
}

func (ts *tradeParamsService) Find(productCode string) (*model.TradeParams, error) {
	return ts.tradeParamsRepository.Find(productCode)
}

func (ts *tradeParamsService) OptimizeEMA(df *model.DataFrame, fastPeriod, slowPeriod int, size float64) (float64, int, int, bool) {
	performance := float64(0)
	bestFastPeriod := fastPeriod
	bestSlowPeriod := slowPeriod

	for fastPeriod := 7; fastPeriod <= 10; fastPeriod++ {
		for slowPeriod := 20; slowPeriod <= 25; slowPeriod++ {
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

	changed := fastPeriod != bestFastPeriod ||
		slowPeriod != bestSlowPeriod

	return performance, bestFastPeriod, bestSlowPeriod, changed
}

func (ts *tradeParamsService) OptimizeBBands(df *model.DataFrame, n int, k float64, size float64) (float64, int, float64, bool) {
	performance := float64(0)
	bestN := n
	bestK := k

	for n := 20; n <= 21; n++ {
		for k := 2.0; k <= 2.0; k += 0.1 {
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

	changed := n != bestN ||
		k != bestK

	return performance, bestN, bestK, changed
}

func (ts *tradeParamsService) OptimizeIchimoku(df *model.DataFrame, size float64) (float64, bool) {
	signalEvents := ts.dataFrameService.BacktestIchimoku(df, size)
	if signalEvents == nil {
		return 0, false
	}
	performance := signalEvents.EstimateProfit()

	return performance, false
}

func (ts *tradeParamsService) OptimizeRSI(df *model.DataFrame, period int, buyThread, sellThread float64, size float64) (float64, int, float64, float64, bool) {
	performance := float64(0)
	bestPeriod := period
	bestBuyThread, bestSellThread := buyThread, sellThread

	for period := 14; period <= 21; period++ {
		for buyThread := float64(25); buyThread <= 35; buyThread++ {
			for sellThread := float64(65); sellThread <= 75; sellThread++ {
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

	changed := period != bestPeriod ||
		buyThread != bestBuyThread ||
		sellThread != bestSellThread

	return performance, bestPeriod, bestBuyThread, bestSellThread, changed
}

func (ts *tradeParamsService) OptimizeMACD(df *model.DataFrame, fastPeriod, slowPeriod, signalPeriod int, size float64) (float64, int, int, int, bool) {
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

	changed := fastPeriod != bestFastPeriod ||
		slowPeriod != bestSlowPeriod ||
		signalPeriod != bestSignalPeriod

	return performance, bestFastPeriod, bestSlowPeriod, bestSignalPeriod, changed
}

func (ts *tradeParamsService) OptimizeAll(df *model.DataFrame, params *model.TradeParams) (*model.TradeParams, bool) {
	_, emaPeriod1, emaPeriod2, emaChanged := ts.OptimizeEMA(df, params.EMAPeriod1(), params.EMAPeriod2(), params.Size())
	_, bbandsN, bbandsK, bbandsChanged := ts.OptimizeBBands(df, params.BBandsN(), params.BBandsK(), params.Size())
	_, rsiPeriod, rsiBuyThread, rsiSellThread, rsiChanged := ts.OptimizeRSI(df, params.RSIPeriod(), params.RSIBuyThread(), params.RSISellThread(), params.Size())
	_, macdFastPeriod, macdSlowPeriod, macdSignalPeriod, macdChanged := ts.OptimizeMACD(df, params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSignalPeriod(), params.Size())

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

	changed := emaChanged ||
		bbandsChanged ||
		rsiChanged ||
		macdChanged

	return newParams, changed
}
