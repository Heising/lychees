// axios基础的封装
import axios, { type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/userStore'
import type { Payload } from '@/models'
import { unix } from '@/utils/tools'
import router from '@/router'
// 是否正在刷新的标记
let isRefreshing = false

//重试中
let Refreshing: Promise<AxiosResponse<any, any>>

const httpInstance = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
  timeout: 5000,
})
// axios请求拦截器
httpInstance.interceptors.request.use(
  async (config) => {
    const userStore = useUserStore()
    const token = userStore.userInfo.token
    if (token) {
      const jwtArray = token.split('.')
      // 拿到负载段，并且base64转json，再json解码
      const payload: Payload = JSON.parse(atob(jwtArray[1]))

      // token过期
      if (payload.tokenInfo.expireUnix < unix()) {
        // 1. 清除本地用户数据
        // 2. 跳转到登录页\
        // 统一错误提示
        ElMessage({
          type: 'warning',
          message: '状态过期，请登录！',
        })
        userStore.clearUserInfo()

        const controller = new AbortController() // 创建新的 AbortController 实例
        // 更新config中的signal为新的controller.signal
        config.signal = controller.signal
        controller.abort() // 取消本次的请求

        router.push('/login')
        return config
      }
      // 如果访问的是刷新的路径
      if (config.url?.startsWith('/refresh')) {
        // 1. 从pinia获取token数据
        // 2. 按照后端的要求拼接token数据
        if (userStore.userInfo.token) {
          config.headers['x-jwt'] = userStore.userInfo.token
        }
        // 在判断isRefreshing之前防止无限递归
        return config
      }
      // 如果在刷新中，则等待
      if (isRefreshing) {
        await Refreshing
      } else if (payload.exp < unix()) {
        // 请求重发
        // jwt已经过期则靠token刷新jwt
        isRefreshing = true
        Refreshing = httpInstance.get('/refresh')
        await Refreshing.finally(() => {
          isRefreshing = false
        })
      } else if (payload.exp - unix() < 2 * 60 * 60) {
        // 两小时刷新
        // 快过期则靠jwt刷新
        isRefreshing = true
        Refreshing = httpInstance.get('/refreshjwt')
        await Refreshing.finally(() => {
          isRefreshing = false
        })
      }

      // 如果访问需有鉴权的路径
      if (config.url?.startsWith('/auth/')) {
        // 1. 从pinia获取token数据
        // 2. 按照后端的要求拼接token数据

        config.headers['x-jwt'] = userStore.userInfo.token
      }

      return config
    }

    // 如果访问需有鉴权的路径
    if (config.url?.startsWith('/auth/')) {
      // 1. 从pinia获取token数据
      // 2. 按照后端的要求拼接token数据

      ElMessage({
        type: 'warning',
        message: '请登录！',
      })

      userStore.clearUserInfo()

      const controller = new AbortController() // 创建新的 AbortController 实例
      // 更新config中的signal为新的controller.signal
      config.signal = controller.signal
      controller.abort() // 取消本次的请求

      router.push('/login')
    }

    return config
  },
  (e) => Promise.reject(e)
)

// axios响应式拦截器
httpInstance.interceptors.response.use(
  (response) => {
    const userStore = useUserStore()

    if (response.headers['x-jwt']) {
      userStore.userInfo.token = response.headers['x-jwt']
    }
    return response.data
  },
  (error) => {
    if (axios.isCancel(error)) {
      // 请求被AbortController终止的清理工作
      // console.log('请求被AbortController终止，执行清理工作')
      // 执行清理操作，例如重置状态或清除数据
      return Promise.reject(error)
    }
    const userStore = useUserStore()

    if (error.response == undefined || error.response.status >= 500) {
      ElMessage({
        grouping: true,
        type: 'error',
        message: '服务器可能存在问题，请待会重试',
      })
      return Promise.reject(error)
    }
    // 统一错误提示
    ElMessage({
      type: 'warning',
      message: error.response.data.error,
      grouping: true,
    })

    // 401token失效处理
    // 1. 清除本地用户数据
    // 2. 跳转到登录页\
    if (error.response.status === 401) {
      userStore.clearUserInfo()

      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default httpInstance
