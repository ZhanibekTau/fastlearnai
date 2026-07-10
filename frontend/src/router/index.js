import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/', component: () => import('@/views/HomePage.vue'), name: 'home' },
  { path: '/login', component: () => import('@/views/LoginPage.vue'), name: 'login' },
  { path: '/register', component: () => import('@/views/LoginPage.vue'), name: 'register' },
  {
    path: '/dashboard',
    component: () => import('@/views/DashboardPage.vue'),
    name: 'dashboard',
    meta: { requiresAuth: true },
  },
  {
    path: '/videos/:id',
    component: () => import('@/views/VideoDetailPage.vue'),
    name: 'video-detail',
    meta: { requiresAuth: true },
  },
  { path: '/pricing', component: () => import('@/views/PricingPage.vue'), name: 'pricing' },
  {
    path: '/account',
    component: () => import('@/views/AccountPage.vue'),
    name: 'account',
    meta: { requiresAuth: true },
  },
  { path: '/public/:id', component: () => import('@/views/VideoDetailPage.vue'), name: 'public-video' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
})

export default router
