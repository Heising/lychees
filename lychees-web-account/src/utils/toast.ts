import { toast } from 'sonner'

export function toastError(message: string | React.ReactNode, duration: number = 4000) {
  toast.error(message, {
    position: 'top-center',
    duration: duration,
    style: { background: 'rgb(251, 113, 133)', color: 'white' },
  })
}
// 默认持续duration10秒
export function toastSuccess(message: string | React.ReactNode) {
  toast.success(message, {
    duration: 10000,
    position: 'top-center',
    style: { background: '#01aebc', color: 'white' },
  })
}
