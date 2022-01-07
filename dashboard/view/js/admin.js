new Vue({
  el: '#app',
  delimiters: ['${', '}'],
  vuetify: new Vuetify(),
  data() {
    return {
      validParams: true,
      productCode: 'ETH_JPY',
      tradeParams: null,
      newTradeParams: null,
      balance: null,
      tradeParamsRules: {
        size: [
          v => !!v || 'size is required',
          v => (v && parseFloat(v) >= 0) || 'size must be more than 0',
        ],
        smaPeriod1: [
          v => !!v || 'smaPeriod1 is required',
          v => (v && v > 0) || 'smaPeriod1 must be more than 0',
        ],
        smaPeriod2: [
          v => !!v || 'smaPeriod2 is required',
          v => (v && v > 0) || 'smaPeriod2 must be more than 0',
        ],
        smaPeriod3: [
          v => !!v || 'smaPeriod3 is required',
          v => (v && v > 0) || 'smaPeriod3 must be more than 0',
        ],
        emaPeriod1: [
          v => !!v || 'emaPeriod1 is required',
          v => (v && v > 0) || 'emaPeriod1 must be more than 0',
        ],
        emaPeriod2: [
          v => !!v || 'emaPeriod2 is required',
          v => (v && v > 0) || 'emaPeriod2 must be more than 0',
        ],
        emaPeriod3: [
          v => !!v || 'emaPeriod3 is required',
          v => (v && v > 0) || 'emaPeriod3 must be more than 0',
        ],
        bbandsN: [
          v => !!v || 'bbandsN is required',
          v => (v && v > 0) || 'bbandsN is must be more than 0',
        ],
        bbandsK: [
          v => !!v || 'bbandsK is required',
          v => (v && parseFloat(v) > 0) || 'bbandsK is must be more than 0',
        ],
        rsiPeriod: [
          v => !!v || 'rsiPeriod is required',
          v => (v && v > 0) || 'rsiPeriod is must be more than 0',
        ],
        rsiBuyThread: [
          v => !!v || 'rsiBuyThread is required',
          v => (v && v >= 0) || 'rsiBuyThread is must be more than 0',
          v => (v && v <= 100) || 'rsiBuyThread is must be less than 100',
        ],
        rsiSellThread: [
          v => !!v || 'rsiSellThread is required',
          v => (v && parseFloat(v) >= 0) || 'rsiSellThread is must be more than 0',
          v => (v && parseFloat(v) <= 100) || 'rsiSellThread is must be less than 100',
        ],
        macdFastPeriod: [
          v => !!v || 'macdFastPeriod is required',
          v => (v && v > 0) || 'macdFastPeriod is must be more than 0',
        ],
        macdSlowPeriod: [
          v => !!v || 'macdSlowPeriod is required',
          v => (v && v > 0) || 'macdSlowPeriod is must be more than 0',
        ],
        macdSignalPeriod: [
          v => !!v || 'macdSignalPeriod is required',
          v => (v && v > 0) || 'macdSignalPeriod is must be more than 0',
        ],
        stopLimitPercent: [
          v => !!v || 'stopLimitPercent is required',
          v => (v && parseFloat(v) >= 0) || 'stopLimitPercent is must be more than 0',
          v => (v && parseFloat(v) <= 1) || 'stopLimitPercent is must be less than 100',
        ],
      },
    }
  },
  methods: {
    async getTradeParams() {
      const params = {
        "productCode": this.productCode,
      }
      return await axios.get('/admin/api/trade-params', {
        params: params,
      }).then(res => {
        return res.data
      }).catch(err => {
        console.log(err)
        return null
      })
    },
    async postTradeParams() {
      return await axios.post('/admin/api/trade-params', {
        ...this.newTradeParams,
      }).then(res => {
        return res.data
      }).catch(err => {
        console.log(err)
        return null
      })
    },
    async updateTradeParams() {
      const res = await this.postTradeParams()
      if (!res) {
        alert('failed to update')
        return
      }
      // 表示するパラメータも更新
      const tradeParams = await this.getTradeParams()
      this.tradeParams = _.cloneDeep(tradeParams)
      this.newTradeParams = _.cloneDeep(tradeParams)
    },
    resetTradeParams() {
      this.newTradeParams = _.cloneDeep(this.tradeParams)
    },
    async getBalance() {
      return await axios.get('/admin/api/balance', {
      }).then(res => {
        return res.data
      }).catch(err => {
        console.log(err)
        return null
      })
    },
  },
  mounted: async function() {
    const tradeParams = await this.getTradeParams()
    this.tradeParams = _.cloneDeep(tradeParams)
    this.newTradeParams = _.cloneDeep(tradeParams)

    this.balance = await this.getBalance()
  },
})
