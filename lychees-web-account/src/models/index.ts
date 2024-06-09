import { z } from 'zod'

// 书签操作后的的响应
export interface UpdateResponse {
  updateAt?: number
  error?: string
}
// 拖拽事件的序号
export interface SwapItem {
  rowIndex: number
  colIndex: number
}
export interface Item {
  isSvg: boolean
  icon: string
  turn: boolean
  color?: string
  // backgroundColor?: string

  title: string
  url: string
  iconSize?: string
}

export interface AddItem extends UpdateResponse {
  data: Item
}
export interface Bookmarks {
  updateAt: number
  iconfontLink: string
  wallpaperUrl: string
  arrayBookmarks: Item[][]
}

export interface User {
  id: number
  email: string
  birthday: number
  nickname: string
  createdAt: number
}

// 登录时的数据定义
export interface LoginUser {
  email: string
  encrypted: string
  nanoid?: string
}
// 注册时的数据定义
export interface RegisterUser extends LoginUser {
  nickname: string
  verifyCode: string
  iconfontLink: string
  birthday: Date
}
// 更新邮箱时的数据定义
export interface UpdateEmail extends LoginUser {
  verifyCode: string
  newEmail: string
}

// 修改密码时的数据定义
export interface PasswordReset extends LoginUser {
  verifyCode: string
}
// 注册后的响应
export interface RegisterResponse {
  message: string
  userId: number
  bookmarksId: number
}

export interface PublicKey {
  publicKey: string
  nanoid: string
}
// 获取图片验证码
export interface CaptchaResponse {
  baseCaptcha: string
  captchaId: string
}

// 验证时的发送类型
export interface VerifyKey {
  captchaId: string
  //验证码
  verifyValue: string
  email: string
  newEmail?: string
  nickname?: string
}
// 创建密码验证规则函数
export const createPasswordValidation = (required_error: string) => {
  return z
    .string({ required_error })
    .min(6, {
      message: '密码须大于等于6位',
    })
    .max(16, {
      message: '密码须小于等于16位',
    })
    .regex(new RegExp('.*[A-Z].*'), '至少一个大写字母')
    .regex(new RegExp('.*[a-z].*'), '至少一个小写字母')
    .regex(new RegExp('.*\\d.*'), '至少一个数字')
    .regex(new RegExp('.*[`~<>?,./!@#$%^&*()\\-_+="\'|{}\\[\\];:\\\\].*'), '至少一个特殊字符')
}

// pin表单规则
export const formSchemaPin = z.object({
  pin: z.string().min(6, {
    message: '必须6个字符',
  }),
})
