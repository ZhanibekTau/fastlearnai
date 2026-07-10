import { defineStore } from 'pinia'
import { ref } from 'vue'
import { videoAPI } from '@/api/client'

export const useVideoStore = defineStore('video', () => {
  const videos = ref([])
  const currentVideo = ref(null)
  const loading = ref(false)
  const uploadProgress = ref(0)
  const totalCount = ref(0)

  async function fetchList(page = 1) {
    loading.value = true
    try {
      const res = await videoAPI.list(page)
      videos.value = res.videos || []
      totalCount.value = res.total_count || 0
    } finally {
      loading.value = false
    }
  }

  async function fetchVideo(id) {
    loading.value = true
    try {
      currentVideo.value = await videoAPI.get(id)
    } finally {
      loading.value = false
    }
  }

  async function upload(file, youtubeURL, language) {
    const formData = new FormData()
    if (file) formData.append('file', file)
    if (youtubeURL) formData.append('youtube_url', youtubeURL)
    formData.append('language', language || 'auto')

    uploadProgress.value = 0
    const video = await videoAPI.upload(formData, (p) => { uploadProgress.value = p })
    videos.value.unshift(video)
    return video
  }

  async function remove(id) {
    await videoAPI.delete(id)
    videos.value = videos.value.filter((v) => v.id !== id)
  }

  async function pollStatus(id, onDone) {
    const interval = setInterval(async () => {
      const s = await videoAPI.status(id)
      if (s.status === 'completed') {
        clearInterval(interval)
        onDone?.(s)
      } else if (s.status === 'failed') {
        clearInterval(interval)
        onDone?.(s)
      }
    }, 5000)
    return () => clearInterval(interval)
  }

  return { videos, currentVideo, loading, uploadProgress, totalCount, fetchList, fetchVideo, upload, remove, pollStatus }
})
