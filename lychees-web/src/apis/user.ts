// 封装所有和用户相关的接口函数
import request from '@/utils/http'
import type { LoginUser, LoginResponse, PublicKey, LoginDevicesType } from '@/models'

export const loginApi = ({ email, encrypted, nanoid }: LoginUser) => {
  return request<LoginResponse>({
    url: '/signin',
    method: 'POST',
    data: {
      email,
      encrypted,
      nanoid,
    },
  })
}
export const logOutApi = () => {
  return request({
    url: '/auth/logout',
    method: 'DELETE',
  })
}

export const getDevicesApi = () => {
  return request<LoginDevicesType>({
    url: '/auth/devices',
  })
}

export const expelDeviceApi = (token: string) => {
  return request({
    url: '/auth/device',
    data: {
      token,
    },
    method: 'DELETE',
  })
}

export const expelDevicesApi = () => {
  return request({
    url: '/auth/devices',
    method: 'DELETE',
  })
}

export const getKeyApi = async () => {
  return request<PublicKey>({
    url: '/getkey',
  })
}
