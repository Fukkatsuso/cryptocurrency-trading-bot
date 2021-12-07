package model

type TradeParams struct {
	tradeEnable      bool
	productCode      string
	size             float64
	smaEnable        bool
	smaPeriod1       int
	smaPeriod2       int
	smaPeriod3       int
	emaEnable        bool
	emaPeriod1       int
	emaPeriod2       int
	emaPeriod3       int
	bbandsEnable     bool
	bbandsN          int
	bbandsK          float64
	ichimokuEnable   bool
	rsiEnable        bool
	rsiPeriod        int
	rsiBuyThread     float64
	rsiSellThread    float64
	macdEnable       bool
	macdFastPeriod   int
	macdSlowPeriod   int
	macdSignalPeriod int
	stopLimitPercent float64
}

func NewTradeParams(tradeEnable bool, productCode string, size float64,
	smaEnable bool, smaPeriod1, smaPeriod2, smaPeriod3 int,
	emaEnable bool, emaPeriod1, emaPeriod2, emaPeriod3 int,
	bbandsEnable bool, bbandsN int, bbandsK float64,
	ichimokuEnable bool,
	rsiEnable bool, rsiPeriod int, rsiBuyThread, rsiSellThread float64,
	macdEnable bool, macdFastPeriod, macdSlowPeriod, macdSignalPeriod int,
	stopLimitPercent float64) *TradeParams {
	if productCode == "" {
		return nil
	}

	if size < 0 {
		return nil
	}

	if smaEnable &&
		(smaPeriod1 <= 0 ||
			smaPeriod2 <= 0 ||
			smaPeriod3 <= 0) {
		return nil
	}

	if emaEnable &&
		(emaPeriod1 <= 0 ||
			emaPeriod2 <= 0 ||
			emaPeriod3 <= 0) {
		return nil
	}

	if bbandsEnable &&
		(bbandsN <= 0 ||
			bbandsK <= 0) {
		return nil
	}

	if rsiEnable &&
		(rsiPeriod <= 0 ||
			rsiBuyThread < 0 || 100 < rsiBuyThread ||
			rsiSellThread < 0 || 100 < rsiSellThread) {
		return nil
	}

	if macdEnable &&
		(macdFastPeriod <= 0 ||
			macdSlowPeriod <= 0 ||
			macdSignalPeriod <= 0) {
		return nil
	}

	if stopLimitPercent < 0 || 100 < stopLimitPercent {
		return nil
	}

	return &TradeParams{
		tradeEnable:      tradeEnable,
		productCode:      productCode,
		size:             size,
		smaEnable:        smaEnable,
		smaPeriod1:       smaPeriod1,
		smaPeriod2:       smaPeriod2,
		smaPeriod3:       smaPeriod3,
		emaEnable:        emaEnable,
		emaPeriod1:       emaPeriod1,
		emaPeriod2:       emaPeriod2,
		emaPeriod3:       emaPeriod3,
		bbandsEnable:     bbandsEnable,
		bbandsN:          bbandsN,
		bbandsK:          bbandsK,
		ichimokuEnable:   ichimokuEnable,
		rsiEnable:        rsiEnable,
		rsiPeriod:        rsiPeriod,
		rsiBuyThread:     rsiBuyThread,
		rsiSellThread:    rsiSellThread,
		macdEnable:       macdEnable,
		macdFastPeriod:   macdFastPeriod,
		macdSlowPeriod:   macdSlowPeriod,
		macdSignalPeriod: macdSignalPeriod,
		stopLimitPercent: stopLimitPercent,
	}
}

func (tp *TradeParams) TradeEnable() bool {
	return tp.tradeEnable
}

func (tp *TradeParams) ProductCode() string {
	return tp.productCode
}

func (tp *TradeParams) Size() float64 {
	return tp.size
}

func (tp *TradeParams) SMAEnable() bool {
	return tp.smaEnable
}

func (tp *TradeParams) SMAPeriod1() int {
	return tp.smaPeriod1
}

func (tp *TradeParams) SMAPeriod2() int {
	return tp.smaPeriod2
}

func (tp *TradeParams) SMAPeriod3() int {
	return tp.smaPeriod3
}

func (tp *TradeParams) EMAEnable() bool {
	return tp.emaEnable
}

func (tp *TradeParams) EMAPeriod1() int {
	return tp.emaPeriod1
}

func (tp *TradeParams) EMAPeriod2() int {
	return tp.emaPeriod2
}

func (tp *TradeParams) EMAPeriod3() int {
	return tp.emaPeriod3
}

func (tp *TradeParams) BBandsEnable() bool {
	return tp.bbandsEnable
}

func (tp *TradeParams) BBandsN() int {
	return tp.bbandsN
}

func (tp *TradeParams) BBandsK() float64 {
	return tp.bbandsK
}

func (tp *TradeParams) IchimokuEnable() bool {
	return tp.ichimokuEnable
}

func (tp *TradeParams) RSIEnable() bool {
	return tp.rsiEnable
}

func (tp *TradeParams) RSIPeriod() int {
	return tp.rsiPeriod
}

func (tp *TradeParams) RSIBuyThread() float64 {
	return tp.rsiBuyThread
}

func (tp *TradeParams) RSISellThread() float64 {
	return tp.rsiSellThread
}

func (tp *TradeParams) MACDEnable() bool {
	return tp.macdEnable
}

func (tp *TradeParams) MACDFastPeriod() int {
	return tp.macdFastPeriod
}

func (tp *TradeParams) MACDSlowPeriod() int {
	return tp.macdSlowPeriod
}

func (tp *TradeParams) MACDSignalPeriod() int {
	return tp.macdSignalPeriod
}

func (tp *TradeParams) StopLimitPercent() float64 {
	return tp.stopLimitPercent
}

func (tp *TradeParams) EnableSMA(enable bool) {
	tp.smaEnable = enable
}

func (tp *TradeParams) EnableEMA(enable bool) {
	tp.emaEnable = enable
}

func (tp *TradeParams) EnableBBands(enable bool) {
	tp.bbandsEnable = enable
}

func (tp *TradeParams) EnableIchimoku(enable bool) {
	tp.ichimokuEnable = enable
}

func (tp *TradeParams) EnableRSI(enable bool) {
	tp.rsiEnable = enable
}

func (tp *TradeParams) EnableMACD(enable bool) {
	tp.macdEnable = enable
}

func NewBasicTradeParams(productCode string, size float64) *TradeParams {
	return NewTradeParams(
		true,
		productCode,
		size,
		true,
		7,
		14,
		50,
		true,
		7,
		14,
		50,
		true,
		20,
		2,
		true,
		true,
		14,
		30,
		70,
		true,
		12,
		26,
		9,
		0.75,
	)
}
