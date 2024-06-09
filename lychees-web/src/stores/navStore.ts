import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Bookmarks } from '@/models'
import { bookmarksFindApi } from '@/apis/bookmark'
import { useUserStore } from '@/stores/userStore'

export const useNavStore = defineStore(
  'nav',
  () => {
    const bookmarks = ref<Bookmarks>({
      updateAt: 0,
      iconfontLink: '',
      arrayBookmarks: [],
      nickname: '',
    })
    const isMask = ref(true)
    async function getBookmarks() {
      useUserStore().userInfo.token
      const result = await bookmarksFindApi(bookmarks.value.updateAt)
      if (result) {
        bookmarks.value = result.data
      }
    }

    function clearBookmarks() {
      bookmarks.value = {
        nickname: '',
        updateAt: 0,
        iconfontLink: '',
        arrayBookmarks: [],
      }
    }
    return { bookmarks, isMask, getBookmarks, clearBookmarks }
  },
  {
    persist: true,
  }
)
