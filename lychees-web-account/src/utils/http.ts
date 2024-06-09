// axios基础的封装
import { toast } from 'sonner'
import axios from 'axios'

const httpInstance = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
  timeout: 5000,
})
// axios响应式拦截器
httpInstance.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response == undefined || error.response.status >= 500) {
      toast.error('服务器可能存在问题，请待会重试', {
        position: 'top-center',
        style: { background: 'rgb(251, 113, 133)', color: 'white' },
      })
      return Promise.reject(error)
    }
    // 统一错误提示
    toast.error(error.response.data.error, {
      position: 'top-center',
      style: { background: 'rgb(251, 113, 133)', color: 'white' },
    })
    // 401token失效处理
    // 1. 清除本地用户数据
    // 2. 跳转到登录页\
    return Promise.reject(error)
  }
)
export default httpInstance
