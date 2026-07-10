<template>
  <header class="sticky top-0 z-50 bg-white dark:bg-gray-950 border-b border-gray-100 dark:border-gray-800">
    <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
      <RouterLink to="/" class="text-xl font-bold text-primary-600">FastLearn</RouterLink>

      <nav class="hidden md:flex items-center gap-6">
        <RouterLink to="/" class="text-sm hover:text-primary-600 transition-colors">{{ t('header.home') }}</RouterLink>
        <RouterLink to="/pricing" class="text-sm hover:text-primary-600 transition-colors">{{ t('header.pricing') }}</RouterLink>
        <RouterLink v-if="auth.isLoggedIn" to="/dashboard" class="text-sm hover:text-primary-600 transition-colors">{{ t('header.dashboard') }}</RouterLink>
        <RouterLink v-if="auth.isAdmin" to="/admin" class="text-sm hover:text-primary-600 transition-colors">{{ t('header.admin') }}</RouterLink>
      </nav>

      <div class="flex items-center gap-3">
        <select
          :value="locale"
          @change="setLocale($event.target.value)"
          class="text-sm bg-transparent border border-gray-200 dark:border-gray-700 rounded-lg px-2 py-1 cursor-pointer"
        >
          <option value="kk">ҚАЗ</option>
          <option value="ru">РУС</option>
          <option value="en">ENG</option>
        </select>

        <button @click="$emit('toggle-dark')" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
          <span class="text-lg">{{ isDark ? '☀️' : '🌙' }}</span>
        </button>

        <template v-if="auth.isLoggedIn">
          <RouterLink to="/account" class="text-sm font-medium hover:text-primary-600">{{ auth.user?.name }}</RouterLink>
          <button @click="handleLogout" class="btn-secondary text-sm">{{ t('common.logout') }}</button>
        </template>
        <template v-else>
          <RouterLink to="/login" class="text-sm hover:text-primary-600">{{ t('common.login') }}</RouterLink>
          <RouterLink to="/register" class="btn-primary text-sm">{{ t('common.signUp') }}</RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { setLocale } from '@/i18n'

defineEmits(['toggle-dark'])

const { t, locale } = useI18n()
const auth = useAuthStore()
const router = useRouter()

const isDark = computed(() => document.documentElement.classList.contains('dark'))

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>
