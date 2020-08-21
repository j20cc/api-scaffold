import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

const routes = [
    { path: '/', name: 'home', component: () => import('../pages/Index.vue') },
    //auth-route-start
    { path: '/login', name: 'login', component: () => import('../pages/auth/Login.vue'), meta: { guest: true } },
    { path: '/register', name: 'register', component: () => import('../pages/auth/Register.vue'), meta: { guest: true } },
    { path: '/forget', name: 'forget', component: () => import('../pages/auth/Forget.vue'), meta: { guest: true } },
    { path: '/password/reset', name: 'reset', component: () => import('../pages/auth/Reset.vue'), meta: { guest: true } },
    { path: '/verification', name: 'verification', component: () => import('../pages/auth/VerifyEmail.vue'), meta: { auth: true } },
    { path: '/user', name: 'user', component: () => import('../pages/auth/User.vue'), meta: { auth: true } },
    //auth-route-end
]

let mode = process.env.NODE_ENV == 'production' ? 'history' : 'hash'
const router = new VueRouter({
    mode: mode,
    routes: routes
})

router.beforeEach((to, from, next) => {
    let user = localStorage.getItem('user')
    if (to.meta.auth && !user) {
        next('/login')
    } else if (to.meta.guest && user) {
        next('/')
    } else {
        next()
    }
})

export default router