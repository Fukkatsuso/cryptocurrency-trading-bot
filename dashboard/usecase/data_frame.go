package usecase

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
)

type DataFrameUsecase interface {
	Get(params *model.TradeParams, candleLimit int64, backtestEnable bool) (*model.DataFrame, error)
}

type dataFrameUsecase struct {
	candleService      service.CandleService
	signalEventService service.SignalEventService
	dataFrameService   service.DataFrameService
}

func NewDataFrameUsecase(cs service.CandleService, ss service.SignalEventService, ds service.DataFrameService) DataFrameUsecase {
	return &dataFrameUsecase{
		candleService:      cs,
		signalEventService: ss,
		dataFrameService:   ds,
	}
}

func (du *dataFrameUsecase) Get(params *model.TradeParams, candleLimit int64, backtestEnable bool) (*model.DataFrame, error) {
	candles, err := du.candleService.FindAll(params.ProductCode(), candleLimit)
	if err != nil {
		return nil, err
	}

	var events []model.SignalEvent
	if len(candles) > 0 {
		timeTime := candles[0].Time().Time()
		events, err = du.signalEventService.FindAllAfterTime(params.ProductCode(), timeTime)
		if err != nil {
			return nil, err
		}
	}
	signalEvents := model.NewSignalEvents(events)

	df := model.NewDataFrame(params.ProductCode(), candles, signalEvents)
	if err != nil {
		return nil, err
	}

	if params.SMAEnable() {
		ok1 := df.AddSMA(params.SMAPeriod1())
		ok2 := df.AddSMA(params.SMAPeriod2())
		ok3 := df.AddSMA(params.SMAPeriod3())
		params.EnableSMA(ok1 && ok2 && ok3)
	}

	if params.EMAEnable() {
		ok1 := df.AddEMA(params.EMAPeriod1())
		ok2 := df.AddEMA(params.EMAPeriod2())
		ok3 := df.AddEMA(params.EMAPeriod3())
		params.EnableEMA(ok1 && ok2 && ok3)
	}

	if params.BBandsEnable() {
		ok := df.AddBBands(params.BBandsN(), params.BBandsK())
		params.EnableBBands(ok)
	}

	if params.IchimokuEnable() {
		ok := df.AddIchimoku()
		params.EnableIchimoku(ok)
	}

	if params.RSIEnable() {
		ok := df.AddRSI(params.RSIPeriod())
		params.EnableRSI(ok)
	}

	if params.MACDEnable() {
		ok := df.AddMACD(params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSignalPeriod())
		params.EnableMACD(ok)
	}

	if backtestEnable {
		du.dataFrameService.Backtest(df, params)
	}

	return df, nil
}
