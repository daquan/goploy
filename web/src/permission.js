import router from './route';
import store from './store';
import NProgress from 'nprogress'; // progress bar
import 'nprogress/nprogress.css';// progress bar style
import {getToken} from '@/utils/auth';

NProgress.configure({showSpinner: false});// NProgress Configuration
const whiteList = ['/login'];// no redirect whitelist

router.beforeEach((to, from, next) => {
  NProgress.start(); // start progress bar
  if (getToken()) {
    if (to.path === '/login') {
      next({path: '/'});
      NProgress.done();
    } else {
      store.dispatch('GetInfo').then((res) => {
        next();
      }).catch(() => {
        next({path: '/login'});
      });
    }
  } else {
    if (whiteList.indexOf(to.path) !== -1) { // 在免登录白名单，直接进入
      next();
    } else {
      next(`/login?redirect=${to.path}`); // 否则全部重定向到登录页
      NProgress.done();
    }
  }
});

router.afterEach(() => {
  NProgress.done(); // finish progress bar
});