new Vue({
  el: '#app',
  delimiters: ['${', '}'],
  vuetify: new Vuetify(),
  data() {
    return {
      candle: null,
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
    drawChart() {
      console.log("drawChart")
    },
    async update() {
      // キャンドルデータとインディケータを取得
      this.candle = await this.getCandle()
      if (!this.candle) {
        return
      }
      // this.candleを使ってグラフ描画
      this.drawChart()
    }
  },
  mounted: async function() {
    this.update()
  },
})
