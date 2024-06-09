import { ref } from 'vue'
import { defineStore } from 'pinia'
export const useFocusStore = defineStore('focus', () => {
  const isFog = ref<boolean>(false)
  const isFocus = ref<boolean>(false)
  const isDisplayEngine = ref<boolean>(false)

  function handleBlur() {
    isFog.value = false
    isDisplayEngine.value = false
  }

  function handleFocus() {
    isFog.value = true
  }

  function change() {
    isFog.value = !isFog.value
    if (isDisplayEngine.value) {
      isDisplayEngine.value = false
    }
  }
  return { isFog, isFocus, isDisplayEngine, handleBlur, change, handleFocus }
})
