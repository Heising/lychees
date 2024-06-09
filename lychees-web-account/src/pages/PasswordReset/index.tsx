import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormMessage,
} from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import { EnvelopeClosedIcon, LockClosedIcon } from '@radix-ui/react-icons'
import { useState } from 'react'
import { getCaptchaApi, passwordResetApi, verifyCaptchaPasswordResetApi } from '@/apis/user'
import { PublicKey, createPasswordValidation } from '@/models'
import JSEncrypt from 'jsencrypt'

import useStore from '@/store'
import { toastError, toastSuccess } from '@/utils/toast'
import CaptchaPinForm from '@/components/CaptchaPinForm'
import EmailPinForm from '@/components/EmailPinForm'
const pubKey: PublicKey = {
  publicKey: '',
  nanoid: '',
}
const PasswordResetForm = () => {
  const formSchema = z
    .object({
      email: z.string({ required_error: '请输入邮箱' }).email({ message: '邮箱格式不正确' }),
      password: createPasswordValidation('请输入密码'),
      confirm: createPasswordValidation('请再一次输入密码，防止输入错误'),
    })
    .refine((data) => data.password === data.confirm, {
      message: '两次密码不一致，请重新输入',
      // path: ['confirm']是一个用于指示错误路径的部分
      path: ['confirm'], // path of error
    })

  const [disabled, setDisabled] = useState(false)
  const setPindisabled = useStore((state) => state.setPindisabled)
  // 更新图片
  const updateImageSrc = useStore((state) => state.updateImageSrc)

  const captchaId = useStore((state) => state.captchaId)
  const setCaptchaId = useStore((state) => state.setCaptchaId)

  const captchaTimestamp = useStore((state) => state.captchaTimestamp)
  const setCaptchaTimestamp = useStore((state) => state.setCaptchaTimestamp)

  const [visibleCaptcha, setVisibleCaptcha] = useState(false)
  const [visibleEmail, setVisibleEmail] = useState(false)
  // let captchadisabled: NodeJS.Timeout

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: '',
      password: '',
      confirm: '',
    },
  })

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    // console.log(import.meta.env)
    // register(values)
    // console.log(useStore.getState().captchaId)
    // console.log(CaptchaId)

    console.log(values.email)
    // console.log('captchaId', captchaId)
    if (captchaId) {
      setVisibleCaptcha(true)
      return
    }
    await refreshCaptcha()
    setVisibleCaptcha(true)
  }
  // 刷新验证码
  async function refreshCaptcha() {
    if (Date.now() - captchaTimestamp < 1000) {
      toastError('验证码请求频繁')
      return
    }
    const res = await getCaptchaApi()
    res.data.baseCaptcha && updateImageSrc(res.data.baseCaptcha)
    setCaptchaTimestamp(Date.now())

    // captchaId = res.data.captchaId
    setCaptchaId(res.data.captchaId)

    formCaptchaPin.reset()
  }

  // 校验验证码
  const handleCaptchaVerify = async () => {
    // console.log(baseCaptcha)
    // console.log(captchaId)
    if (Date.now() - captchaTimestamp > 60000) {
      toastError('图片验证码过期，请重新输入')
      await refreshCaptcha()
      return
    }

    await verifyCaptchaPasswordResetApi({
      captchaId: captchaId,
      verifyValue: formCaptchaPin.watch('pin'),
      email: form.watch('email'),
    })
      .then((result) => {
        pubKey.nanoid = result.data.nanoid
        pubKey.publicKey = result.data.publicKey
        setVisibleCaptcha(false)

        formEmailPin.reset()
        setCaptchaId('')

        setVisibleEmail(true)
      })
      .catch((err) => {
        if (err.response.data.error === '验证码错误！') {
          refreshCaptcha()

        } else if (err.response.data.error === '用户不存在') {
          form.resetField('email')

          setVisibleCaptcha(false)
        } else if (err.response.data.error === '同一邮箱只能五分钟申请一次邮箱验证码') {
          pubKey.nanoid = err.response.data.data.nanoid
          pubKey.publicKey = err.response.data.data.publicKey
          setVisibleCaptcha(false)

          // 打开邮箱验证码表单组件
          formEmailPin.reset()
          setVisibleEmail(true)
        }
      })
  }
  const FormSchemaPin = z.object({
    pin: z.string().min(6, {
      message: '必须6个字符',
    }),
  })
  const formCaptchaPin = useForm<z.infer<typeof FormSchemaPin>>({
    resolver: zodResolver(FormSchemaPin),
    defaultValues: {
      pin: '',
    },
  })
  const formEmailPin = useForm<z.infer<typeof FormSchemaPin>>({
    resolver: zodResolver(FormSchemaPin),
    defaultValues: {
      pin: '',
    },
  })
  // 处理重置密码
  function onPasswordReset(data: z.infer<typeof FormSchemaPin>) {
    console.log(data)
    setVisibleEmail(false)
    console.log('表单状态', form.formState.errors)
    // 判断表单是否通过校验
    if (Object.keys(form.formState.errors).length === 0) {
      console.log('表单通过校验')
      passwordReset(
        {
          email: form.watch('email'),
          password: form.watch('password'),
          confirm: form.watch('confirm'),
        },
        data.pin
      )
    } else {
      toastError('表单校验未通过')
    }
  }
  // 开始重置密码
  const passwordReset = async (account: z.infer<typeof formSchema>, verifyCode: string) => {
    const email = account.email


    const encryptor = new JSEncrypt() // 新建JSEncrypt对象
    encryptor.setPublicKey(pubKey.publicKey) // 设置公钥
    const encrypted = encryptor.encrypt(pubKey.nanoid + account.password)
    console.log(encrypted)
    const nanoid = pubKey.nanoid
    if (encrypted) {
      passwordResetApi({ verifyCode, email, encrypted, nanoid })
        .then((response) => {
          toastSuccess(response.data.message)
          // 密码重置成功，清空表单
          form.reset()
        })
        .catch((err) => {
          setVisibleEmail(true)
          formEmailPin.reset()
          if (err.response == undefined || err.response.status >= 500) {
            return Promise.reject(err)
          }
          if (err.response.data.error === '邮箱验证码错误') {
            setDisabled(true)

            // 定时器的id
            setPindisabled(
              setTimeout(() => {
                setDisabled(false)
                clearTimeout(useStore.getState().pinDisabled)
              }, 5000)
            )
            toastError('请5秒后再校验邮箱验证码', 5000)
          }
        })
    } else {
      toastError('密码加密失败，请重试')
    }
  }
  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="max-w-screen-sm container sm relative top-36"
        >
          <h3 className="text-center text-slate-500 text-3xl mb-3.5">重置密码</h3>
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Input placeholder="邮箱号" type="email" autoComplete="email" {...field}>
                    <EnvelopeClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请输入你使用的邮箱</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Input
                    placeholder="新密码"
                    type="password"
                    autoComplete="new-password"
                    {...field}
                  >
                    <LockClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请输入你的新密码, 6-16位包含特殊字符和大小写</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="confirm"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Input
                    placeholder="新密码"
                    type="password"
                    autoComplete="new-password"
                    {...field}
                  >
                    <LockClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请再次输入你的密码校验</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Separator />
          <div className="flex justify-between space-y-5">
            <a
              href="https://redli.cn/login"
              className="cursor-pointer select-none inline-flex items-center justify-center whitespace-nowrap rounded-md text-bese font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2 self-end"
            >
              前往登录
            </a>

            <Button type="submit">修改密码</Button>
          </div>
        </form>
      </Form>
      {/* 图片验证码 */}
      <CaptchaPinForm
        visibleCaptcha={visibleCaptcha}
        setVisibleCaptcha={setVisibleCaptcha}
        refreshCaptcha={refreshCaptcha}
        handleSubmit={handleCaptchaVerify}
        formCaptchaPin={formCaptchaPin}
      ></CaptchaPinForm>

      {/* 邮箱验证码 */}
      <EmailPinForm
        visibleEmail={visibleEmail}
        setVisibleEmail={setVisibleEmail}
        email={form.watch('email')}
        formEmailPin={formEmailPin}
        handleSubmit={onPasswordReset}
        disabled={disabled}
      ></EmailPinForm>
    </>
  )
}

export default PasswordResetForm
