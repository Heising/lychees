// 封装所有和用户相关的接口函数
import request from '@/utils/http'
import type {
  RegisterUser,
  PublicKey,
  RegisterResponse,
  VerifyKey,
  CaptchaResponse,
  PasswordReset,
  UpdateEmail,
} from '@/models'

// 注册用户
export const registerApi = (registerUser: RegisterUser) => {
  return request<RegisterResponse>({
    url: '/signup',
    method: 'POST',
    data: {
      ...registerUser,
    },
  })
}

// 密码重置
export const passwordResetApi = (registerUser: PasswordReset) => {
  return request<RegisterResponse>({
    url: '/password_reset',
    method: 'POST',
    data: {
      ...registerUser,
    },
  })
}

// 更新用户邮箱
export const updateEmailApi = (updateEmail: UpdateEmail) => {
  return request<RegisterResponse>({
    url: '/email_update',
    method: 'POST',
    data: {
      ...updateEmail,
    },
  })
}

// 获取加密用的公钥，如果不是https访问的情况下
export const getKeyApi = async () => {
  return request<PublicKey>({
    url: '/getkey',
  })
}
// 获取验证码
export const getCaptchaApi = async () => {
  return request<CaptchaResponse>({
    url: '/captcha',
  })
}

// 校验注册验证码是否正确 需携带昵称
export const verifyCaptchaApi = async (verify: VerifyKey) => {
  return request<PublicKey>({
    url: '/verify_captcha_register',
    method: 'POST',
    data: { ...verify },
  })
}
// 校验验证码是否正确 密码重置 不用携带昵称
export const verifyCaptchaPasswordResetApi = async (verify: VerifyKey) => {
  return request<PublicKey>({
    url: '/verify_captcha_password_reset',
    method: 'POST',
    data: { ...verify },
  })
}
// 校验验证码是否正确 邮箱变更 不用携带昵称
export const verifyCaptchaEmaiUpdatelApi = async (verify: VerifyKey) => {
  return request<PublicKey>({
    url: '/verify_captcha_email_update',
    method: 'POST',
    data: { ...verify },
  })
}