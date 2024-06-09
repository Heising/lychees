import { UseFormReturn } from 'react-hook-form'

import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'

import { InputOTP, InputOTPGroup, InputOTPSlot } from '@/components/ui/input-otp'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import useStore from '@/store'

const CaptchaPinForm = ({
  visibleCaptcha,
  setVisibleCaptcha,
  handleSubmit,
  refreshCaptcha,
  formCaptchaPin,
}: {
  visibleCaptcha: boolean
  setVisibleCaptcha: React.Dispatch<React.SetStateAction<boolean>>
  handleSubmit: () => Promise<void>
  refreshCaptcha: () => Promise<void>
  formCaptchaPin: UseFormReturn<{
    pin: string
  }>
}) => {

  return (
    <>
      {/* 图片验证码 */}
      <AlertDialog open={visibleCaptcha} onOpenChange={setVisibleCaptcha}>
        <AlertDialogContent className="sm:max-w-md">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-center">
              请输入数字校验码，点击图片可以刷新
            </AlertDialogTitle>
          </AlertDialogHeader>
          <div className="flex justify-center">
            <img
              src={useStore((state) => state.imageSrc)}
              width="240"
              height="80"
              onClick={refreshCaptcha}
            />
          </div>
          <Form {...formCaptchaPin}>
            <form
              onSubmit={formCaptchaPin.handleSubmit(handleSubmit)}
              className="w-full space-y-1 relative"
            >
              <FormField
                control={formCaptchaPin.control}
                name="pin"
                render={({ field }) => (
                  <FormItem>
                    {/* <FormLabel>One-Time Password</FormLabel> */}
                    <FormControl>
                      <div className="grid flex-1 gap-2 justify-center">
                        <InputOTP maxLength={6} {...field}>
                          <InputOTPGroup>
                            <InputOTPSlot index={0} />
                            <InputOTPSlot index={1} />
                            <InputOTPSlot index={2} />
                            <InputOTPSlot index={3} />
                            <InputOTPSlot index={4} />
                            <InputOTPSlot index={5} />
                          </InputOTPGroup>
                        </InputOTP>
                      </div>
                    </FormControl>
                    {/* <FormDescription>
                    Please enter the one-time password sent to your phone.
                  </FormDescription> */}
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogCancel>取消</AlertDialogCancel>
              <Button type="submit" className="absolute right-0 bottom-0">
                提交
              </Button>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}

export default CaptchaPinForm
