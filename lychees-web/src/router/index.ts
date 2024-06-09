

import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/views/Layout/index.vue'


const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Layout,
      meta: {
        title: '红荔书签页',
        auth: true,
      },
    },
    {
      path: '/login',
      component: () => import('@/views/Login/index.vue'),
      meta: {
        title: '登录到红荔书签页',
      },
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})
router.beforeEach((to, _, next) => {
  const title = to.meta.title as string
  if (title) {
    document.title = title
  }
  next()
})

export default router
