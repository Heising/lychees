import { createBrowserRouter } from 'react-router-dom'
import App from '@/App.tsx'
import { lazy } from 'react'

// import PasswordReset from '@/pages/PasswordReset'
// import Register from '@/pages/Register'
// import EmailUpdate from '@/pages/EmailUpdate'

// React 组件懒加载
// 使用React.lazy动态导入路由组件
const Register = lazy(() => import('@/pages/Register'))
const PasswordReset = lazy(() => import('@/pages/PasswordReset'))
const EmailUpdate = lazy(() => import('@/pages/EmailUpdate'))

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/password_reset',
    element: <PasswordReset />,
  },
  {
    path: '/email_update',
    element: <EmailUpdate />,
  },
])

export default router
