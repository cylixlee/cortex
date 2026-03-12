import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Login from '@/views/Login.vue'
import ChatView from '@/views/ChatView.vue'
import SkillList from '@/views/SkillList.vue'
import SkillUpload from '@/views/SkillUpload.vue'
import SkillDetail from '@/views/SkillDetail.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/chat',
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { guest: true },
    },
    {
      path: '/chat',
      name: 'chat',
      component: ChatView,
      meta: { requiresAuth: true },
    },
    {
      path: '/chat/:id',
      name: 'chat-conversation',
      component: ChatView,
      meta: { requiresAuth: true },
    },
    {
      path: '/skills',
      name: 'skills',
      component: SkillList,
      meta: { requiresAuth: true },
    },
    {
      path: '/skills/upload',
      name: 'skill-upload',
      component: SkillUpload,
      meta: { requiresAuth: true },
    },
    {
      path: '/skills/:id',
      name: 'skill-detail',
      component: SkillDetail,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const isAuthenticated = userStore.isLoggedIn

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.guest && isAuthenticated) {
    next('/chat')
  } else {
    next()
  }
})

export default router
