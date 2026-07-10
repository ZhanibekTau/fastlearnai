<template>
  <div class="max-w-4xl mx-auto px-4 py-10">
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-2xl font-bold">{{ t('dashboard.myVideos') }}</h1>
      <RouterLink to="/" class="btn-primary">{{ t('dashboard.newVideo') }}</RouterLink>
    </div>

    <div v-if="usage" class="card mb-6 flex gap-6 text-sm">
      <div>
        <span class="text-gray-500">{{ t('common.plan') }}:</span>
        <span class="font-semibold ml-2 capitalize">{{ usage.plan }}</span>
      </div>
      <div>
        <span class="text-gray-500">{{ t('dashboard.videosThisMonth') }}</span>
        <span class="font-semibold ml-2">{{ usage.videos_count }} / {{ usage.limits?.max_videos === -1 ? '∞' : usage.limits?.max_videos }}</span>
      </div>
      <RouterLink to="/pricing" class="text-primary-600 hover:underline ml-auto">{{ t('dashboard.upgrade') }}</RouterLink>
    </div>

    <LoadingSpinner v-if="store.loading" :message="t('dashboard.loadingVideos')" />

    <VideoList
      v-else
      :videos="store.videos"
      @delete="deleteVideo"
    />

    <div v-if="store.totalCount > 20" class="mt-6 flex justify-center gap-2">
      <button class="btn-secondary text-sm" :disabled="page <= 1" @click="loadPage(page - 1)">← {{ t('common.prev') }}</button>
      <button class="btn-secondary text-sm" @click="loadPage(page + 1)">{{ t('common.next') }} →</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { useVideoStore } from '@/stores/video'
import { useSubscriptionStore } from '@/stores/subscription'
import VideoList from '@/components/VideoList.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const { t } = useI18n()
const store = useVideoStore()
const subStore = useSubscriptionStore()
const page = ref(1)
const usage = ref(null)

onMounted(async () => {
  await store.fetchList(1)
  await subStore.fetchMine()
  usage.value = subStore.usage
})

async function loadPage(p) {
  page.value = p
  await store.fetchList(p)
}

async function deleteVideo(id) {
  if (!confirm(t('dashboard.confirmDelete'))) return
  await store.remove(id)
}
</script>
