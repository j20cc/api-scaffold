import axios from 'axios'
import { Notification, Loading } from 'element-ui'
import router from './routes'
import store from './store'

let loadingInstance = null
const instance = axios.create({
  timeout: 10000,
  // baseURL: process.env.NODE_ENV === 'production' ? prohost : devhost
  baseURL: '/api'
})

let httpCode = {
  400: '请求参数错误',
  401: '登录状态失效',
  403: '服务器拒绝本次访问',
  404: '请求资源未找到',
  500: '内部服务器错误',
  501: '服务器不支持该请求中使用的方法',
  502: '网关错误',
  504: '网关超时'
}

instance.interceptors.request.use(config => {
  let user = store.state.user
  if (user != null) {
    config.headers['Authorization'] = "Bearer " + user.token
  }
  loadingInstance = Loading.service({
    text: "加载中...",
    spinner: "el-icon-loading"
  });

  return config
}, error => {
  return Promise.reject(error)
})

instance.interceptors.response.use(response => {
  loadingInstance.close()
  return Promise.resolve(response.data)
}, error => {
  loadingInstance.close()
  if (error.response) {
    let tips = error.response.status in httpCode ? httpCode[error.response.status] : error.response.data.error
    Notification.error({
      title: '错误',
      message: error.response.data.error || tips
    });
    if (error.response.status === 401) {
      router.push({
        path: `/login`
      })
    }
    return Promise.reject(error)
  } else {
    Notification.error({
      title: "错误",
      message: '请求超时, 请刷新重试'
    });
    return Promise.reject(new Error('请求超时, 请刷新重试'))
  }
})

const get = (url, params, config = {}) => {
  return new Promise((resolve, reject) => {
    instance({
      method: 'get',
      url,
      params,
      ...config
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}

const post = (url, data, config = {}) => {
  return new Promise((resolve, reject) => {
    instance({
      method: 'post',
      url,
      data,
      ...config
    }).then(response => {
      resolve(response)
    }).catch(error => {
      reject(error)
    })
  })
}

export default {
  get,
  post
}
