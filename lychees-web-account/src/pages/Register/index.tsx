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
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { cn } from '@/lib/utils'
import {
  CalendarIcon,
  EnvelopeClosedIcon,
  Link2Icon,
  LockClosedIcon,
  PersonIcon,
} from '@radix-ui/react-icons'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import { useState } from 'react'
import { getCaptchaApi, registerApi, verifyCaptchaApi } from '@/apis/user'
import { PublicKey, createPasswordValidation, formSchemaPin } from '@/models'
import JSEncrypt from 'jsencrypt'

import useStore from '@/store'
import { toastError, toastSuccess } from '@/utils/toast'
import CaptchaPinForm from '@/components/CaptchaPinForm'
import EmailPinForm from '@/components/EmailPinForm'
const pubKey: PublicKey = {
  publicKey: '',
  nanoid: '',
}
const RegisterForm = () => {
  // 注册表单的规则
  const formSchema = z
    .object({
      nickname: z
        .string({ required_error: '请输入昵称' })
        .min(1, {
          message: '昵称密码须大于等于1位',
        })
        .max(16, {
          message: '昵称须小于等于16位',
        }),
      iconfontLink: z.string(),
      email: z.string({ required_error: '请输入邮箱' }).email({ message: '邮箱格式不正确' }),
      password: createPasswordValidation('请输入密码'),
      confirm: createPasswordValidation('请再一次输入密码，防止输入错误'),
      birthday: z.date({ required_error: '请选择日期' }),
    })
    .refine((data) => data.password === data.confirm, {
      message: '两次密码不一致，请重新输入',
      // path: ['confirm']是一个用于指示错误路径的部分
      path: ['confirm'], // path of error
    })
  // 更新图片
  const { updateImageSrc, setCaptchaId, captchaTimestamp, setCaptchaTimestamp } = useStore()
  // 设置图片验证码表单组件打开
  const [visibleCaptcha, setVisibleCaptcha] = useState(false)
  // 设置邮箱验证码表单组件打开
  const [visibleEmail, setVisibleEmail] = useState(false)
  // let captchadisabled: NodeJS.Timeout

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nickname: '',
      email: '',
      password: '',
      confirm: '',
      iconfontLink: '',
      birthday: undefined,
    },
  })

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    // console.log(import.meta.env)
    // register(values)

    console.log(values.email)
    // console.log('captchaId', captchaId)
    if (useStore.getState().captchaId) {
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

    await verifyCaptchaApi({
      captchaId: useStore.getState().captchaId,
      verifyValue: formCaptchaPin.watch('pin'),
      email: form.watch('email'),
      nickname: form.watch('nickname'),
    })
      .then((result) => {
        pubKey.nanoid = result.data.nanoid
        pubKey.publicKey = result.data.publicKey
        setVisibleCaptcha(false)
        formCaptchaPin.reset()

        // 打开邮箱验证码表单组件
        formEmailPin.reset()
        setVisibleEmail(true)
      })
      .catch((err) => {
        // formCaptchaPin.reset()
        if (err.response == undefined || err.response.status >= 500) {
          return Promise.reject(err)
        }
        if (err.response.data.error === '验证码错误！') {
          // 刷新验证码
          refreshCaptcha()
        } else if (err.response.data.error === '邮箱已经注册，请登录！') {
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

  // 验证码表单
  const formCaptchaPin = useForm<z.infer<typeof formSchemaPin>>({
    resolver: zodResolver(formSchemaPin),
    defaultValues: {
      pin: '',
    },
  })
  // 邮箱验证码表单
  const formEmailPin = useForm<z.infer<typeof formSchemaPin>>({
    resolver: zodResolver(formSchemaPin),
    defaultValues: {
      pin: '',
    },
  })
  // 处理注册
  function onRegister(data: z.infer<typeof formSchemaPin>) {
    // captchaId = ''
    setCaptchaId('')
    console.log(data)
    setVisibleEmail(false)
    console.log('表单状态', form.formState.errors)
    // 判断表单是否通过校验
    if (Object.keys(form.formState.errors).length === 0) {
      console.log('表单通过校验')
      register(
        {
          nickname: form.watch('nickname'),
          iconfontLink: form.watch('iconfontLink'),
          email: form.watch('email'),
          password: form.watch('password'),
          confirm: form.watch('confirm'),
          birthday: form.watch('birthday'),
        },
        data.pin
      )
    } else {
      toastError('表单校验未通过')
    }
  }
  // 开始注册
  const register = async (account: z.infer<typeof formSchema>, verifyCode: string) => {
    const email = account.email
    const nickname = account.nickname
    const iconfontLink = account.iconfontLink
    const birthday = account.birthday

    const encryptor = new JSEncrypt() // 新建JSEncrypt对象
    encryptor.setPublicKey(pubKey.publicKey) // 设置公钥
    const encrypted = encryptor.encrypt(pubKey.nanoid + account.password)
    console.log(encrypted)
    const nanoid = pubKey.nanoid
    if (encrypted) {
      registerApi({ verifyCode, email, nickname, encrypted, nanoid, iconfontLink, birthday })
        .then((response) => {
          toastSuccess(response.data.message)
          // 注册成功，清空表单
          form.reset()
        })
        .catch(() => {
          setVisibleEmail(true)
          formEmailPin.reset()
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
          <h3 className="text-center text-slate-500 text-3xl mb-3.5">创建账号</h3>
          <FormField
            control={form.control}
            name="nickname"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="昵称" type="text" autoComplete="nickname" {...field}>
                    <PersonIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请输入昵称方便我们称呼你</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Input placeholder="邮箱" type="email" autoComplete="email" {...field}>
                    <EnvelopeClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请输入常用邮箱详见白名单</FormDescription>
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
                  <Input placeholder="密码" type="password" autoComplete="new-password" {...field}>
                    <LockClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>6-16位包含特殊字符和大小写</FormDescription>
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
                  <Input placeholder="密码" type="password" autoComplete="new-password" {...field}>
                    <LockClosedIcon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请再次输入你的密码校验</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="iconfontLink"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Input placeholder="链接" type="url" {...field}>
                    <Link2Icon className="ml-auto h-4 w-4 opacity-50 absolute left-3.5 top-1/2 -translate-y-1/2" />
                  </Input>
                </FormControl>
                <FormDescription>请输入第三方图标链接，选填</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="birthday"
            render={({ field }) => (
              <FormItem>
                {/* <FormLabel>Username</FormLabel> */}
                <FormControl>
                  <Popover>
                    <PopoverTrigger asChild>
                      <FormControl>
                        <Button
                          variant={'outline'}
                          className={cn(
                            'w-full font-normal relative pl-8 justify-start text-muted-foreground',
                            !field.value && 'text-muted-foreground'
                          )}
                        >
                          {field.value ? (
                            format(field.value, 'PPP', { locale: zhCN })
                          ) : (
                            <span>生日</span>
                          )}
                          <CalendarIcon className="absolute left-3.5 ml-auto h-4 w-4 opacity-50" />
                        </Button>
                      </FormControl>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={field.value}
                        onSelect={field.onChange}
                        disabled={(date: Date) =>
                          date > new Date() || date < new Date('1970-01-01')
                        }
                        initialFocus
                        locale={zhCN}
                        captionLayout="dropdown-buttons"
                        // defaultMonth={new Date()}
                        fromMonth={new Date('1970-01-01')}
                        toDate={new Date()}
                      />
                    </PopoverContent>
                  </Popover>
                </FormControl>
                <FormDescription>用作生日祝贺相关，可随意填写</FormDescription>
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

            <Button type="submit">开始注册</Button>
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
        handleSubmit={onRegister}
        disabled={false}
      ></EmailPinForm>
    </>
  )
}

export default RegisterForm
