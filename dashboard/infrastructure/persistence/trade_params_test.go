package persistence_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func newTradeParamsList() []model.TradeParams {
	table := []struct {
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
	}{
		{
			tradeEnable:      false,
			productCode:      config.ProductCode,
			size:             0,
			smaEnable:        false,
			smaPeriod1:       0,
			smaPeriod2:       0,
			smaPeriod3:       0,
			emaEnable:        false,
			emaPeriod1:       0,
			emaPeriod2:       0,
			emaPeriod3:       0,
			bbandsEnable:     false,
			bbandsN:          0,
			bbandsK:          0,
			ichimokuEnable:   false,
			rsiEnable:        false,
			rsiPeriod:        0,
			rsiBuyThread:     0,
			rsiSellThread:    0,
			macdEnable:       false,
			macdFastPeriod:   0,
			macdSlowPeriod:   0,
			macdSignalPeriod: 0,
			stopLimitPercent: 0,
		},
		{
			tradeEnable:      true,
			productCode:      config.ProductCode,
			size:             0.01,
			smaEnable:        true,
			smaPeriod1:       7,
			smaPeriod2:       14,
			smaPeriod3:       50,
			emaEnable:        true,
			emaPeriod1:       7,
			emaPeriod2:       14,
			emaPeriod3:       50,
			bbandsEnable:     true,
			bbandsN:          20,
			bbandsK:          2.2,
			ichimokuEnable:   true,
			rsiEnable:        true,
			rsiPeriod:        14,
			rsiBuyThread:     30.5,
			rsiSellThread:    70.5,
			macdEnable:       true,
			macdFastPeriod:   12,
			macdSlowPeriod:   26,
			macdSignalPeriod: 9,
			stopLimitPercent: 0.75,
		},
	}

	tradeParamsList := make([]model.TradeParams, 0)
	for _, t := range table {
		tradeParams := model.NewTradeParams(
			t.tradeEnable,
			t.productCode,
			t.size,
			t.smaEnable,
			t.smaPeriod1,
			t.smaPeriod2,
			t.smaPeriod3,
			t.emaEnable,
			t.emaPeriod1,
			t.emaPeriod2,
			t.emaPeriod3,
			t.bbandsEnable,
			t.bbandsN,
			t.bbandsK,
			t.ichimokuEnable,
			t.rsiEnable,
			t.rsiPeriod,
			t.rsiBuyThread,
			t.rsiSellThread,
			t.macdEnable,
			t.macdFastPeriod,
			t.macdSlowPeriod,
			t.macdSignalPeriod,
			t.stopLimitPercent,
		)
		if tradeParams == nil {
			continue
		}
		tradeParamsList = append(tradeParamsList, *tradeParams)
	}
	return tradeParamsList
}

func TestTradeParams(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)

	tradeParamsList := newTradeParamsList()

	t.Run("save trade_params", func(t *testing.T) {
		for _, tradeParams := range tradeParamsList {
			err := tradeParamsRepository.Save(tradeParams)
			if err != nil {
				t.Fatal(err.Error())
			}
		}
	})

	t.Run("find trade_params", func(t *testing.T) {
		lastTradeParams := tradeParamsList[len(tradeParamsList)-1]
		productCode := lastTradeParams.ProductCode()
		tradeParams, err := tradeParamsRepository.Find(productCode)
		if err != nil {
			t.Fatal(err.Error())
		}
		if tradeParams == nil ||
			*tradeParams != lastTradeParams {
			t.Fatalf("%+v != %+v", *tradeParams, lastTradeParams)
		}
	})
}
