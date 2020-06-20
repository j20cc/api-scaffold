import store from './store'
import router from './routes'
import { Notification } from 'element-ui'

const getUser = () => {
  return JSON.parse(window.localStorage.getItem('user'))
}

const setUser = (user) => {
  store.dispatch('setUser', user)
  window.localStorage.setItem('user', JSON.stringify(user))
}

const removeUser = () => {
  store.dispatch('setUser', null)
  window.localStorage.removeItem('user')
}

const notifyAndRedirect = (type, message, path, duration) => {
  if (type == 'success') {
    Notification.success({
      title: "成功",
      message: message
    });
  } else {
    Notification.error({
      title: "失败",
      message: message
    });
  }
  setTimeout(() => {
    router.push({
      path: path
    });
  }, duration);
}

export default {
  getUser,
  setUser,
  removeUser,
  notifyAndRedirect
}