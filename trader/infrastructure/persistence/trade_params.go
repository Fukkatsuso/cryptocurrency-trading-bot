package persistence

import (
	"errors"
	"fmt"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type tradeParamsRepository struct {
	db DB
}

func NewTradeParamsRepository(db DB) repository.TradeParamsRepository {
	return &tradeParamsRepository{
		db: db,
	}
}

func (tr *tradeParamsRepository) Save(tp model.TradeParams) error {
	cmd := fmt.Sprintf(`
        INSERT INTO trade_params (
            trade_enable,
            product_code,
            size,
            sma_enable,
            sma_period1,
            sma_period2,
            sma_period3,
            ema_enable,
            ema_period1,
            ema_period2,
            ema_period3,
            bbands_enable,
            bbands_n,
            bbands_k,
            ichimoku_enable,
            rsi_enable,
            rsi_period,
            rsi_buy_thread,
            rsi_sell_thread,
            macd_enable,
            macd_fast_period,
            macd_slow_period,
            macd_signal_period,
            stop_limit_percent
        )
        VALUES (
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?
        )
        `,
	)

	_, err := tr.db.Exec(cmd,
		tp.TradeEnable(),
		tp.ProductCode(),
		tp.Size(),
		tp.SMAEnable(),
		tp.SMAPeriod1(),
		tp.SMAPeriod2(),
		tp.SMAPeriod3(),
		tp.EMAEnable(),
		tp.EMAPeriod1(),
		tp.EMAPeriod2(),
		tp.EMAPeriod3(),
		tp.BBandsEnable(),
		tp.BBandsN(),
		tp.BBandsK(),
		tp.IchimokuEnable(),
		tp.RSIEnable(),
		tp.RSIPeriod(),
		tp.RSIBuyThread(),
		tp.RSISellThread(),
		tp.MACDEnable(),
		tp.MACDFastPeriod(),
		tp.MACDSlowPeriod(),
		tp.MACDSignalPeriod(),
		tp.StopLimitPercent(),
	)
	return err
}

func (tr *tradeParamsRepository) Find(productCode string) (*model.TradeParams, error) {
	// 最後に作成されたパラメータを取得
	// productCodeで絞り込み，そのうちcreated_atが最新のレコードを探す
	cmd := fmt.Sprintf(`
            SELECT
                tp.trade_enable,
                tp.size,
                tp.sma_enable,
                tp.sma_period1,
                tp.sma_period2,
                tp.sma_period3,
                tp.ema_enable,
                tp.ema_period1,
                tp.ema_period2,
                tp.ema_period3,
                tp.bbands_enable,
                tp.bbands_n,
                tp.bbands_k,
                tp.ichimoku_enable,
                tp.rsi_enable,
                tp.rsi_period,
                tp.rsi_buy_thread,
                tp.rsi_sell_thread,
                tp.macd_enable,
                tp.macd_fast_period,
                tp.macd_slow_period,
                tp.macd_signal_period,
                tp.stop_limit_percent
            FROM
                trade_params AS tp
            WHERE
                tp.product_code = ? AND
                tp.created_at = (
                    SELECT
                        MAX(sub_tp.created_at)
                    FROM
                        trade_params AS sub_tp
                    WHERE
                        sub_tp.product_code = ?
                )`,
	)
	row := tr.db.QueryRow(cmd, productCode, productCode)

	var tradeEnable bool
	var size float64
	var smaEnable bool
	var smaPeriod1, smaPeriod2, smaPeriod3 int
	var emaEnable bool
	var emaPeriod1, emaPeriod2, emaPeriod3 int
	var bbandsEnable bool
	var bbandsN int
	var bbandsK float64
	var ichimokuEnable bool
	var rsiEnable bool
	var rsiPeriod int
	var rsiBuyThread, rsiSellThread float64
	var macdEnable bool
	var macdFastPeriod, macdSlowPeriod, macdSignalPeriod int
	var stopLimitPercent float64
	err := row.Scan(
		&tradeEnable,
		&size,
		&smaEnable,
		&smaPeriod1,
		&smaPeriod2,
		&smaPeriod3,
		&emaEnable,
		&emaPeriod1,
		&emaPeriod2,
		&emaPeriod3,
		&bbandsEnable,
		&bbandsN,
		&bbandsK,
		&ichimokuEnable,
		&rsiEnable,
		&rsiPeriod,
		&rsiBuyThread,
		&rsiSellThread,
		&macdEnable,
		&macdFastPeriod,
		&macdSlowPeriod,
		&macdSignalPeriod,
		&stopLimitPercent,
	)
	if err != nil {
		return nil, err
	}

	tradeParams := model.NewTradeParams(
		tradeEnable,
		productCode,
		size,
		smaEnable,
		smaPeriod1,
		smaPeriod2,
		smaPeriod3,
		emaEnable,
		emaPeriod1,
		emaPeriod2,
		emaPeriod3,
		bbandsEnable,
		bbandsN,
		bbandsK,
		ichimokuEnable,
		rsiEnable,
		rsiPeriod,
		rsiBuyThread,
		rsiSellThread,
		macdEnable,
		macdFastPeriod,
		macdSlowPeriod,
		macdSignalPeriod,
		stopLimitPercent,
	)
	if tradeParams == nil {
		return nil, errors.New(fmt.Sprint("invalid trade_params:",
			tradeEnable,
			size,
			smaEnable,
			smaPeriod1,
			smaPeriod2,
			smaPeriod3,
			emaEnable,
			emaPeriod1,
			emaPeriod2,
			emaPeriod3,
			bbandsEnable,
			bbandsN,
			bbandsK,
			ichimokuEnable,
			rsiEnable,
			rsiPeriod,
			rsiBuyThread,
			rsiSellThread,
			macdEnable,
			macdFastPeriod,
			macdSlowPeriod,
			macdSignalPeriod,
			stopLimitPercent,
		))
	}
	return tradeParams, nil
}
