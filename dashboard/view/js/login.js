new Vue({
  el: '#app',
  delimiters: ['${', '}'],
  vuetify: new Vuetify(),
  data() {
    return {
      valid: true,
      userId: '',
      password: '',
      showPassword: false,
      userIdRules: [
        v => (v && v.length > 0)  || 'user ID is required',
        v => (v && v.length <= 50) || 'user ID must be less than 50 characters'
      ],
      passwordRules: [
        v => (v && v.length > 0)  || 'password is required',
      ],
    }
  },
  methods: {
    async submit() {
      const params = new FormData()
      params.append('userId', this.userId)
      params.append('password', this.password)
      await axios.post('/api/login', params, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }).then(res => {
        window.location.href = '/admin'
      }).catch(err => {
        console.log(err)
        window.alert('failed to login')
      })
    },
  },
})
