import { Button } from '@/components/ui/button'
import { Link } from 'react-router-dom'
const App = () => {
  return (
    <div className="container grid gap-4 grid-cols-2 mt-48">
      <Button asChild>
        <Link to="https://redli.cn/">前往使用</Link>
      </Button>
      <Button asChild>
        <Link to="https://www.redli.cn/">查看文档</Link>
      </Button>
      <Button asChild>
        <Link to="/register">注册账号</Link>
      </Button>
      <Button asChild>
        <Link to="/password_reset">重置密码</Link>
      </Button>
      <Button asChild>
        <Link to="/email_update">修改邮箱</Link>
      </Button>
    </div>
  )
}

export default App
