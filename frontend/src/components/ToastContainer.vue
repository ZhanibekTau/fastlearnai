<template>
  <Teleport to="body">
    <div class="fixed bottom-5 right-5 z-50 flex flex-col gap-2 w-80 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in items"
          :key="toast.id"
          class="flex items-start gap-3 px-4 py-3 rounded-2xl shadow-lg pointer-events-auto cursor-pointer"
          :class="toastClass(toast.type)"
          @click="dismiss(toast.id)"
        >
          <span class="text-lg leading-none mt-0.5 shrink-0">{{ icon(toast.type) }}</span>
          <p class="text-sm font-medium leading-snug">{{ toast.message }}</p>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { useToast } from '@/composables/useToast'

const { items, dismiss } = useToast()

function toastClass(type) {
  return {
    error:   'bg-red-50 text-red-800 border border-red-200',
    success: 'bg-green-50 text-green-800 border border-green-200',
    info:    'bg-blue-50 text-blue-800 border border-blue-200',
    warn:    'bg-amber-50 text-amber-800 border border-amber-200',
  }[type] ?? 'bg-stone-50 text-stone-800 border border-stone-200'
}

function icon(type) {
  return { error: '✕', success: '✓', info: 'ℹ', warn: '⚠' }[type] ?? 'ℹ'
}
</script>

<style scoped>
.toast-enter-active  { transition: all 0.25s ease; }
.toast-leave-active  { transition: all 0.2s ease; }
.toast-enter-from    { opacity: 0; transform: translateX(100%); }
.toast-leave-to      { opacity: 0; transform: translateX(100%); }
</style>
