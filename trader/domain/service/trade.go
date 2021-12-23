package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type TradeService interface {
	Trade(productCode string, pastPeriod int) error
	Buy(events *model.SignalEvents, productCode string, size float64, timeTime time.Time) error
	Sell(events *model.SignalEvents, productCode string, size float64, timeTime time.Time) error
}

type tradeService struct {
	balanceRepository     repository.BalanceRepository
	tickerRepository      repository.TickerRepository
	orderRepository       repository.OrderRepository
	signalEventRepository repository.SignalEventRepository
	candleService         CandleService
	dataFrameService      DataFrameService
	tradeParamsService    TradeParamsService
}

func NewTradeService(
	br repository.BalanceRepository,
	tr repository.TickerRepository,
	or repository.OrderRepository,
	sr repository.SignalEventRepository,
	cs CandleService,
	ds DataFrameService,
	ts TradeParamsService,
) TradeService {
	return &tradeService{
		balanceRepository:     br,
		tickerRepository:      tr,
		orderRepository:       or,
		signalEventRepository: sr,
		candleService:         cs,
		dataFrameService:      ds,
		tradeParamsService:    ts,
	}
}

func (ts *tradeService) Trade(productCode string, pastPeriod int) error {
	params, err := ts.tradeParamsService.Find(productCode)
	if err != nil {
		return err
	}
	if !params.TradeEnable() {
		return errors.New("trade is not enabled")
	}

	candles, err := ts.candleService.FindAll(productCode, int64(pastPeriod))
	if err != nil {
		return err
	}

	events, err := ts.signalEventRepository.FindAll(productCode)
	if err != nil {
		return err
	}
	signalEvents := model.NewSignalEvents(events)
	if signalEvents == nil {
		return errors.New("can't make a SignalEvents instance")
	}

	df := model.NewDataFrame(productCode, candles, signalEvents)

	if params.EMAEnable() {
		ok1 := df.AddEMA(params.EMAPeriod1())
		ok2 := df.AddEMA(params.EMAPeriod2())
		params.EnableEMA(ok1 && ok2)
	}

	if params.BBandsEnable() {
		ok := df.AddBBands(params.BBandsN(), params.BBandsK())
		params.EnableBBands(ok)
	}

	if params.IchimokuEnable() {
		ok := df.AddIchimoku()
		params.EnableIchimoku(ok)
	}

	if params.MACDEnable() {
		ok := df.AddMACD(params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSignalPeriod())
		params.EnableMACD(ok)
	}

	if params.RSIEnable() {
		ok := df.AddRSI(params.RSIPeriod())
		params.EnableRSI(ok)
	}

	now := len(candles) - 1
	buy, sell := ts.dataFrameService.Analyze(df, now, params)

	if buy {
		nowTime := time.Now().UTC()
		err := ts.Buy(signalEvents, productCode, params.Size(), nowTime)
		if err != nil {
			return err
		}
	}

	currentPrice := candles[now].Close()
	if sell ||
		signalEvents.ShouldCutLoss(currentPrice, params.StopLimitPercent()) {
		nowTime := time.Now().UTC()
		err := ts.Sell(signalEvents, productCode, params.Size(), nowTime)
		if err != nil {
			return err
		}

		// パラメータ更新
		var changed bool
		params, changed = ts.tradeParamsService.OptimizeAll(df, params)
		if changed {
			err := ts.tradeParamsService.Save(*params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ts *tradeService) Buy(events *model.SignalEvents, productCode string, size float64, timeTime time.Time) error {
	if !events.CanBuyAt(timeTime) {
		return errors.New("[Buy] can't buy due to signal_event's history")
	}

	// 所持中の現金
	codes := strings.Split(productCode, "_")
	currencyCode := codes[1]
	balance, err := ts.balanceRepository.FetchByCurrencyCode(currencyCode)
	if err != nil {
		return err
	}
	availableCurrency := balance.Available()

	// 現在の価格
	ticker, err := ts.tickerRepository.Fetch(productCode)
	if err != nil {
		return err
	}
	needCurrency := ticker.BestAsk() * size

	// お金が足りないときは購入しない
	if availableCurrency < needCurrency {
		return errors.New(fmt.Sprintf("[Buy] you don't have enough money. available: %f, need: %f", availableCurrency, needCurrency))
	}

	// 買い注文
	order := model.NewBuyOrder(productCode, size)
	if order == nil {
		return errors.New("[Buy] can't make a new order instance")
	}
	fmt.Printf("[Buy] order: %+v\n", order)

	// 注文送信
	completedOrder, err := ts.orderRepository.Send(*order)
	if err != nil {
		fmt.Println("[Buy]", err)
		return err
	}
	fmt.Printf("[Buy] order completed: %+v\n", completedOrder)

	// SignalEvent
	signalEvent := model.NewSignalEvent(timeTime, productCode, model.OrderSideBuy, completedOrder.AveragePrice, completedOrder.Size)
	if signalEvent == nil {
		return errors.New("[Buy] order send, but signal_event is nil")
	}
	events.AddBuySignal(*signalEvent)

	// SingalEventをDBに保存
	err = ts.signalEventRepository.Save(*signalEvent)
	if err != nil {
		return err
	}

	return nil
}

func (ts *tradeService) Sell(events *model.SignalEvents, productCode string, size float64, timeTime time.Time) error {
	if !events.CanSellAt(timeTime) {
		return errors.New("[Sell] can't sell due to signal_event's history")
	}

	// 所持中の仮想通貨
	codes := strings.Split(productCode, "_")
	coinCode := codes[0]
	balance, err := ts.balanceRepository.FetchByCurrencyCode(coinCode)
	if err != nil {
		return err
	}
	availableCoin := balance.Available()

	// パラメータに設定したサイズよりも保有量が足りないときは保有量だけ使う
	if availableCoin < size {
		size = availableCoin
	}

	// 売り注文
	order := model.NewSellOrder(productCode, size)
	if order == nil {
		return errors.New("[Sell] can't make a new order instance")
	}
	fmt.Printf("[Sell] order: %+v\n", order)

	// 注文送信
	completedOrder, err := ts.orderRepository.Send(*order)
	if err != nil {
		fmt.Println("[Sell]", err)
		return err
	}
	fmt.Printf("[Sell] order completed: %+v\n", completedOrder)

	// SignalEvent
	signalEvent := model.NewSignalEvent(timeTime, productCode, model.OrderSideSell, completedOrder.AveragePrice, completedOrder.Size)
	if signalEvent == nil {
		return errors.New("[Sell] order send, but signal_event is nil")
	}
	events.AddSellSignal(*signalEvent)

	// SingalEventをDBに保存
	err = ts.signalEventRepository.Save(*signalEvent)
	if err != nil {
		return err
	}

	return nil
}
