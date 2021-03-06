const chartOptionsBase = {
  chart: {
    type: 'candlestick',
    height: 400,
  },
  title: {
    text: 'CandleStick Chart',
    align: 'left',
  },
  xaxis: { // https://apexcharts.com/docs/options/xaxis/
    type: 'datetime',
  },
  yaxis: {
    tooltip: {
      enabled: true,
    },
  },
}

new Vue({
  el: '#app',
  delimiters: ['${', '}'],
  vuetify: new Vuetify(),
  components: {
    apexchart: VueApexCharts,
  },
  data() {
    return {
      candle: null,
      config: {
        limit: 30,
        sma: {
          enable: false,
          periods: [7, 14, 50],
        },
        ema: {
          enable: false,
          periods: [7, 14, 50],
        },
        bbands: {
          enable: false,
          n: 20,
          k: 2,
        },
        ichimoku: {
          enable: false,
        },
        rsi: {
          enable: false,
          period: 14,
        },
        macd: {
          enable: false,
          periods: [12, 26, 9],
        },
        backtest: {
          enable: false,
        }
      }
    }
  },
  methods: {
    async getCandle() {
      let params = {
        "limit": this.config.limit,
        "sma": this.config.sma.enable,
        "smaPeriod1": this.config.sma.periods[0],
        "smaPeriod2": this.config.sma.periods[1],
        "smaPeriod3": this.config.sma.periods[2],
        "ema": this.config.ema.enable,
        "emaPeriod1": this.config.ema.periods[0],
        "emaPeriod2": this.config.ema.periods[1],
        "emaPeriod3": this.config.ema.periods[2],
        "bbands": this.config.bbands.enable,
        "bbandsN": this.config.bbands.n,
        "bbandsK": this.config.bbands.k,
        "ichimoku": this.config.ichimoku.enable,
        "rsi": this.config.rsi.enable,
        "rsiPeriod": this.config.rsi.period,
        "macd": this.config.macd.enable,
        "macdPeriod1": this.config.macd.periods[0],
        "macdPeriod2": this.config.macd.periods[1],
        "macdPeriod3": this.config.macd.periods[2],
        "backtest": this.config.backtest.enable,
      }
      return await axios.get('/api/candle', {
        params: params,
      }).then(res => {
        return res.data
      }).catch(err => {
        console.log(err)
        return null
      })
    },
    async update() {
      // ?????????????????????????????????????????????????????????
      this.candle = await this.getCandle()
    },
  },
  computed: {
    series() {
      if (!this.candle || !this.candle.candles) {
        return [{
          data: []
        }]
      }
      const data = this.candle.candles.map(c => {
        return {
          x: new Date(c['time']).getTime(),
          y: [c['open'], c['high'], c['low'], c['close']],
        }
      })
      return [{
        data: data,
      }]
    },
    chartOptions() {
      const annotations = {
        xaxis: [
          ...this.tradeEventAnnotationXaxis,
          ...this.backTestEventAnnotationXaxis,
        ],
      }
      const options = {
        ...chartOptionsBase,
        annotations: annotations,
      }
      return options
    },
    tradeEventAnnotationXaxis() {
      const color = '#00E396'
      if (this.candle && this.candle.events && this.candle.events.signals) {
        const xaxis = this.candle.events.signals.map(s => {
          return {
            x: new Date(s['time']).getTime(),
            borderColor: color,
            label: {
              borderColor: color,
              style: {
                fontSize: '12px',
                color: '#fff',
                background: color,
              },
              orientation: 'horizontal',
              offsetY: 10,
              text: s['side'],
            },
          }
        })
        return xaxis
      }
      return []
    },
    backTestEventAnnotationXaxis() {
      const color = '#3C90EB'
      if (this.candle && this.candle.backtestEvents && this.candle.backtestEvents.signals) {
        const xaxis = this.candle.backtestEvents.signals.map(s => {
          return {
            x: new Date(s['time']).getTime(),
            borderColor: color,
            label: {
              borderColor: color,
              style: {
                fontSize: '12px',
                color: '#fff',
                background: color,
              },
              orientation: 'horizontal',
              offsetY: chartOptionsBase.chart.height-92,
              text: s['side'],
            },
          }
        })
        return xaxis
      }
      return []
    },
    // ???????????????????????????????????????????????????????????????
    backtestCurrentHold() {
      if (!this.candle.backtestEvents || !this.candle.backtestEvents.signals) {
        return 0
      }
      let hold = 0
      for (const signal of this.candle.backtestEvents.signals) {
        if (signal.side == "BUY") {
          hold -= signal.size
        } else if (signal.side == "SELL") {
          hold += signal.size
        }
      }
      return hold
    }
  },
  mounted: async function() {
    this.update()
  },
})
