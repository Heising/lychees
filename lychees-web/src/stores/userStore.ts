// 管理用户数据相关
import { defineStore } from 'pinia'
import { loginApi } from '@/apis/user'
import { useNavStore } from './navStore'
import { useSearchHistory } from '@/stores/historyStore'

import { cleanData } from '@/stores/db'
import { ref } from 'vue'
import type { LoginUser, LoginResponse } from '@/models'
export const useUserStore = defineStore(
  'user',
  () => {
    // 1. 定义管理用户数据的state
    const userInfo = ref<LoginResponse>({
      token: '',
      userInfo: {
        id: 0,
        email: '',
        nickname: '',
        createdAt: 0,
      },
    })
    // 2. 定义获取接口数据的action函数
    async function getUserInfo({ email, encrypted, nanoid }: LoginUser) {
      const result = await loginApi({
        email,
        encrypted,
        nanoid,
      })
      userInfo.value = result.data
    }

    // 重置用户信息
    async function clearUserInfo() {
      userInfo.value = {
        token: '',
        userInfo: {
          id: 0,
          email: '',
          nickname: '',
          createdAt: 0,
        },
      }
      useNavStore().clearBookmarks()
      useSearchHistory().clean()
      cleanData()
    }
    return {
      userInfo,
      getUserInfo,
      clearUserInfo,
    }
  },
  {
    persist: true,
  }
)
