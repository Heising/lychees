import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import 'virtual:svg-icons-register'
import 'element-plus/theme-chalk/base.css'
import 'element-plus/theme-chalk/el-message.css'
import 'element-plus/theme-chalk/el-notification.css'

import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import router from './router'
const app = createApp(App)

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia()

// 注册持久化插件
pinia.use(piniaPluginPersistedstate)
// 注册pinia
app.use(pinia)
// ElementPlus改为中文
app.use(ElementPlus, {
  locale: zhCn,
})
// 注册路由
app.use(router)

app.mount('#app')
