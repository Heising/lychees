import { create } from 'zustand'
import { combine } from 'zustand/middleware'
const useStore = create(
  combine(
    {
      // 获取到验证码的时间戳
      captchaTimestamp: 0,
      // 验证码id
      captchaId: '',
      // 禁用输入的定时器
      pinDisabled: setTimeout(() => {
        clearTimeout(useStore.getState().pinDisabled)
      }, 0),
      // 图片url
      imageSrc: '',
    },
    (set) => ({
      setCaptchaTimestamp: (newCaptchaTimestamp: number) =>
        set({ captchaTimestamp: newCaptchaTimestamp }),

      setCaptchaId: (newCaptchaId: string) => set({ captchaId: newCaptchaId }),

      setPindisabled: (newPindisabled: NodeJS.Timeout) =>
        set({ pinDisabled: newPindisabled }),

      updateImageSrc: (newImageSrc: string) => set({ imageSrc: newImageSrc }),
    })
  )
)
export default useStore
