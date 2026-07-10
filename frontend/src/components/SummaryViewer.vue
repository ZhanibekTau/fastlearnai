<template>
  <div v-if="summary" class="space-y-6">
    <div class="card">
      <h2 class="text-2xl font-bold mb-2">{{ summary.title }}</h2>
      <div class="flex gap-4 text-sm text-gray-500 dark:text-gray-400 mb-4">
        <span>{{ summary.language?.toUpperCase() }}</span>
        <span v-if="summary.duration">{{ summary.duration }}</span>
      </div>
      <p class="text-gray-700 dark:text-gray-300 leading-relaxed">{{ summary.summary?.overview }}</p>
    </div>

    <div class="card">
      <h3 class="font-semibold text-lg mb-4">{{ t('summary.keyPoints') }}</h3>
      <div v-for="(point, i) in summary.summary?.key_points" :key="i" class="mb-4 last:mb-0">
        <h4 class="font-medium text-primary-600 mb-1">{{ point.topic }}</h4>
        <p class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">{{ point.explanation }}</p>
      </div>
    </div>

    <div class="card">
      <h3 class="font-semibold text-lg mb-3">{{ t('summary.keyTakeaways') }}</h3>
      <ul class="space-y-2">
        <li
          v-for="(takeaway, i) in summary.summary?.key_takeaways"
          :key="i"
          class="flex gap-2 text-sm"
        >
          <span class="text-primary-500 mt-0.5">✓</span>
          <span>{{ takeaway }}</span>
        </li>
      </ul>
    </div>

    <div v-if="summary.qa?.length" class="card">
      <h3 class="font-semibold text-lg mb-4">{{ t('summary.qa') }}</h3>
      <div v-for="(item, i) in summary.qa" :key="i" class="mb-4 last:mb-0">
        <p class="font-medium mb-1">{{ t('summary.questionPrefix') }} {{ item.question }}</p>
        <p class="text-sm text-gray-600 dark:text-gray-400">{{ t('summary.answerPrefix') }} {{ item.answer }}</p>
      </div>
    </div>

    <div class="flex gap-3">
      <button class="btn-secondary text-sm" @click="$emit('download', 'pdf')">{{ t('summary.downloadPdf') }}</button>
      <button class="btn-secondary text-sm" @click="$emit('download', 'md')">{{ t('summary.downloadMarkdown') }}</button>
      <button class="btn-secondary text-sm" @click="copyLink">{{ t('summary.copyLink') }}</button>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const props = defineProps({ summary: Object })
defineEmits(['download'])

const { t } = useI18n()

function copyLink() {
  navigator.clipboard.writeText(window.location.href)
}
</script>
