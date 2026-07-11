<template>
  <div class="max-w-4xl mx-auto px-4 py-10">
    <RouterLink to="/dashboard" class="text-sm text-gray-400 hover:text-gray-600 mb-6 inline-block">← {{ t('common.back') }}</RouterLink>

    <LoadingSpinner
      v-if="store.loading || isPolling"
      :message="isPolling ? t('videoDetail.processingMessage') : t('common.loading')"
    />

    <div v-else-if="store.currentVideo">
      <div v-if="store.currentVideo.status === 'failed'" class="card bg-red-50 dark:bg-red-950 border-red-200">
        <p class="text-red-600 font-medium">{{ t('videoDetail.processingFailed') }}</p>
        <p class="text-sm text-red-400 mt-1">{{ store.currentVideo.error_message }}</p>
        <button class="btn-primary mt-4" @click="regenerate">{{ t('videoDetail.retry') }}</button>
      </div>

      <SummaryViewer
        v-else-if="store.currentVideo.status === 'completed' && summary"
        :summary="summary"
        @download="handleDownload"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import { useVideoStore } from '@/stores/video'
import { videoAPI } from '@/api/client'
import SummaryViewer from '@/components/SummaryViewer.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const { t } = useI18n()
const route = useRoute()
const store = useVideoStore()
const isPolling = ref(false)
const summary = ref(null)
let stopPolling = null

onMounted(async () => {
  const id = route.params.id
  await store.fetchVideo(id)
  const video = store.currentVideo

  if (video?.status === 'completed') {
    summary.value = video
  } else if (video?.status === 'pending' || video?.status === 'processing') {
    isPolling.value = true
    stopPolling = await store.pollStatus(id, async (s) => {
      isPolling.value = false
      if (s.status === 'completed') {
        await store.fetchVideo(id)
        summary.value = store.currentVideo
      }
    })
  }
})

onUnmounted(() => stopPolling?.())

async function regenerate() {
  await videoAPI.regenerate(route.params.id)
  isPolling.value = true
  stopPolling = await store.pollStatus(route.params.id, async () => {
    isPolling.value = false
    await store.fetchVideo(route.params.id)
    summary.value = store.currentVideo
  })
}

async function handleDownload(format) {
  try {
    const blob = await videoAPI.download(route.params.id, format)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${store.currentVideo?.title || 'summary'}.${format}`
    a.click()
    URL.revokeObjectURL(url)
  } catch (err) {
    alert(err.message)
  }
}
</script>
