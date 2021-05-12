const chartOptions = {
  chart: {
    type: 'candlestick',
    height: 350,
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
      series: null,
      chartOptions: chartOptions,
    }
  },
  methods: {
    async getCandle() {
      let params = {
        "limit": 10,
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
    seriesData() {
      const data = this.candle.candles.map(c => {
        return {
          x: new Date(c['time']),
          y: [c['open'], c['high'], c['low'], c['close']],
        }
      })
      return [{
        data: data,
      }]
    },
    async update() {
      // キャンドルデータとインディケータを取得
      this.candle = await this.getCandle()
      if (!this.candle) {
        return
      }
      // apexchartに渡すデータ
      this.series = this.seriesData()
    }
  },
  mounted: async function() {
    this.update()
  },
})
