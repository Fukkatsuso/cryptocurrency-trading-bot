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
      chartOptions: chartOptions,
      config: {
        limit: 30,
      }
    }
  },
  methods: {
    async getCandle() {
      let params = {
        "limit": this.config.limit,
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
      // キャンドルデータとインディケータを取得
      this.candle = await this.getCandle()
    }
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
          x: new Date(c['time']),
          y: [c['open'], c['high'], c['low'], c['close']],
        }
      })
      return [{
        data: data,
      }]
    }
  },
  mounted: async function() {
    this.update()
  },
})
