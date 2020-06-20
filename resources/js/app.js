import Vue from 'vue'
import App from './App.vue'
import router from './libs/routes'
import store from './libs/store'

Vue.config.productionTip = false

//element
import { Notification, Loading, MessageBox } from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
Vue.prototype.$notify = Notification
Vue.prototype.$loading = Loading.service
Vue.prototype.$confirm = MessageBox.confirm
//http
import http from './libs/http'
Vue.prototype.$http = http

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
