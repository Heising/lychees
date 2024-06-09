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
import { formSchemaPin } from '@/models'
import { z } from 'zod'

const EmailPinForm = ({
  visibleEmail,
  setVisibleEmail,
  email,
  formEmailPin,
  handleSubmit,
  disabled,
}: {
  visibleEmail: boolean
  setVisibleEmail: React.Dispatch<React.SetStateAction<boolean>>
  email: string
  formEmailPin: UseFormReturn<{
    pin: string
  }>
  handleSubmit: (data: z.infer<typeof formSchemaPin>) => void
  disabled: boolean
}) => {
  return (
    <>
      {/* 邮箱验证码 */}
      <AlertDialog open={visibleEmail} onOpenChange={setVisibleEmail}>
        <AlertDialogContent className="sm:max-w-md">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-center">
              输入 {email} 收到的数字验证码
            </AlertDialogTitle>
          </AlertDialogHeader>
          <Form {...formEmailPin}>
            <form
              onSubmit={formEmailPin.handleSubmit(handleSubmit)}
              className="w-full space-y-1 relative"
            >
              <FormField
                control={formEmailPin.control}
                name="pin"
                render={({ field }) => (
                  <FormItem>
                    {/* <FormLabel>One-Time Password</FormLabel> */}
                    <FormControl>
                      <div className="grid flex-1 gap-2 justify-center">
                        <InputOTP maxLength={6} {...field} disabled={disabled}>
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

export default EmailPinForm
