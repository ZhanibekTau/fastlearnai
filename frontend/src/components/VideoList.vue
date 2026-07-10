<template>
  <div>
    <div v-if="!videos.length" class="text-center py-12 text-gray-400">
      {{ t('videoList.empty') }}
    </div>
    <div v-else class="space-y-3">
      <div
        v-for="video in videos"
        :key="video.id"
        class="card flex items-center justify-between hover:shadow-md transition-shadow cursor-pointer"
        @click="$router.push(`/videos/${video.id}`)"
      >
        <div class="flex items-center gap-4">
          <div class="w-10 h-10 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center text-xl">
            🎬
          </div>
          <div>
            <p class="font-medium">{{ video.title || video.original_filename || video.youtube_url || t('videoList.untitled') }}</p>
            <p class="text-sm text-gray-400">{{ formatDate(video.created_at) }} · {{ video.language?.toUpperCase() }}</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <span
            class="text-xs px-2 py-1 rounded-full font-medium"
            :class="statusClass(video.status)"
          >
            {{ statusLabel(video.status) }}
          </span>
          <button
            class="text-red-400 hover:text-red-600 p-1"
            @click.stop="$emit('delete', video.id)"
          >
            ✕
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

defineProps({ videos: { type: Array, default: () => [] } })
defineEmits(['delete'])

const { t, locale } = useI18n()

function formatDate(iso) {
  return new Date(iso).toLocaleDateString(locale.value)
}

function statusLabel(status) {
  return t(`videoList.status.${status}`)
}

function statusClass(status) {
  return {
    completed: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
    processing: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
    failed: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
    pending: 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300',
  }[status] || 'bg-gray-100 text-gray-600'
}
</script>
