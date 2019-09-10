import Vue from 'vue'
import vuetify from './plugins/vuetify'
import VueRouter from 'vue-router'
import $auth from "./components/auth";
import App from './components/App';
import Login from './components/Login';
import NewAccount from './components/NewAccount';
import Dashboard from './components/Dashboard';
import UserDashboard from './components/UsersDashboard';
import Report from './components/Report';

Vue.use(VueRouter);
Vue.use(vuetify);

function requireAuth(to, from, next) {
    if (!$auth.isLoggedIn()) {
        next({
            path: '/login',
            query: {redirect: to.fullPath}
        })
    } else {
        next()
    }
}

const router = new VueRouter({
    mode: 'history',
    base: '/savings/',
    routes: [
        {path: '/login', component: Login},
        {path: '/new', component: NewAccount},
        {path: '/dashboard', component: Dashboard, alias: '/', beforeEnter: requireAuth},
        {path: '/users', component: UserDashboard, beforeEnter: requireAuth},
        {path: '/report', component: Report, beforeEnter: requireAuth},
        {
            path: '/logout', beforeEnter(to, from, next) {
                $auth.logout();
                next('/login');
            }
        },
    ]
});

new Vue({
    el: '#app',
    router,
    vuetify,
    render: h => h(App),
});
