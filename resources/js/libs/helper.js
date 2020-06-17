import store from './store'

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

export default {
    getUser,
    setUser,
    removeUser
}