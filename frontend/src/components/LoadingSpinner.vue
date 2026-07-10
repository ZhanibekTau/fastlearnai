<template>
  <div class="flex flex-col items-center gap-4 py-12 text-gray-500">
    <div class="w-10 h-10 border-4 border-gray-200 border-t-primary-500 rounded-full animate-spin" />
    <p class="text-sm">{{ message || t('common.processing') }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

defineProps({ message: { type: String, default: null } })
const { t } = useI18n()
</script>
