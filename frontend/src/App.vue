<template>
  <div :class="{ dark: isDark }">
    <Header @toggle-dark="toggleDark" />
    <main class="min-h-screen">
      <RouterView />
    </main>
    <ToastContainer />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import Header from '@/components/Header.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { useAuthStore } from '@/stores/auth'

const isDark = ref(localStorage.getItem('theme') === 'dark')
const auth = useAuthStore()

function toggleDark() {
  isDark.value = !isDark.value
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  if (auth.isLoggedIn) auth.fetchMe().catch(() => {})
})
</script>
