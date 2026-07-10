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
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import Header from '@/components/Header.vue'
import ToastContainer from '@/components/ToastContainer.vue'

const isDark = ref(localStorage.getItem('theme') === 'dark')

function toggleDark() {
  isDark.value = !isDark.value
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}
</script>
