import { ref } from 'vue'
import { defineStore } from 'pinia'
export const useSearchHistory = defineStore(
  'history',
  () => {
    // 创建一个响应式的列表
    const searchHistoryArray = ref<string[]>([])

    const add = (keyword: string) => {
      searchHistoryArray.value = [
        keyword,
        ...searchHistoryArray.value.filter((item) => item !== keyword).slice(0, 9),
      ]
    }

    const clean = () => {
      searchHistoryArray.value = []
    }
    const cleanIndex = (index: number) => {
      searchHistoryArray.value.splice(index, 1)
    }
    return { searchHistoryArray, add, clean, cleanIndex }
  },
  {
    persist: true,
  }
)
