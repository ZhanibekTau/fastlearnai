<template>
  <div>
    <div
      class="border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-xl p-10 text-center transition-colors"
      :class="{ 'border-primary-500 bg-primary-50 dark:bg-primary-950': isDragging }"
      @dragover.prevent="isDragging = true"
      @dragleave="isDragging = false"
      @drop.prevent="onDrop"
    >
      <div class="text-4xl mb-3">🎬</div>
      <p class="font-medium mb-1">{{ t('upload.dropZone') }}</p>
      <p class="text-sm text-gray-500 mb-4">{{ t('upload.formats') }}</p>
      <button class="btn-primary" @click="fileInput.click()">{{ t('upload.chooseFile') }}</button>
      <input ref="fileInput" type="file" class="hidden" accept="video/*,audio/*" @change="onFileChange" />
    </div>

    <div class="my-4 flex items-center gap-3">
      <div class="flex-1 h-px bg-gray-200 dark:bg-gray-700" />
      <span class="text-sm text-gray-400">{{ t('upload.orYoutube') }}</span>
      <div class="flex-1 h-px bg-gray-200 dark:bg-gray-700" />
    </div>

    <input
      v-model="youtubeURL"
      type="url"
      class="input"
      placeholder="https://www.youtube.com/watch?v=..."
    />

    <div class="mt-4 flex gap-3 items-center">
      <label class="text-sm font-medium">{{ t('upload.languageLabel') }}</label>
      <select v-model="language" class="input max-w-xs">
        <option value="auto">{{ t('upload.autoDetect') }}</option>
        <option value="en">English</option>
        <option value="ru">Русский</option>
        <option value="kk">Қазақша</option>
      </select>
    </div>

    <div v-if="selectedFile" class="mt-3 text-sm text-gray-600 dark:text-gray-400">
      {{ t('upload.selected', { name: selectedFile.name, size: formatSize(selectedFile.size) }) }}
    </div>

    <div v-if="uploading" class="mt-4">
      <div class="flex justify-between text-sm mb-1">
        <span>{{ t('upload.uploading') }}</span>
        <span>{{ store.uploadProgress }}%</span>
      </div>
      <div class="h-2 bg-stone-200 rounded-full overflow-hidden">
        <div class="h-full bg-primary-500 transition-all" :style="{ width: store.uploadProgress + '%' }" />
      </div>
    </div>

    <button
      class="btn-primary w-full mt-6"
      :disabled="uploading || (!selectedFile && !youtubeURL)"
      @click="submit"
    >
      {{ uploading ? t('upload.uploading') : t('upload.generateSummary') }}
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useVideoStore } from '@/stores/video'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'

const emit = defineEmits(['uploaded'])

const { t } = useI18n()
const store = useVideoStore()
const router = useRouter()
const toast  = useToast()

const isDragging = ref(false)
const selectedFile = ref(null)
const youtubeURL = ref('')
const language = ref('auto')
const fileInput = ref(null)
const uploading = ref(false)

function onDrop(e) {
  isDragging.value = false
  selectedFile.value = e.dataTransfer.files[0] || null
}

function onFileChange(e) {
  selectedFile.value = e.target.files[0] || null
}

function formatSize(bytes) {
  return bytes > 1_000_000 ? `${(bytes / 1_000_000).toFixed(1)} MB` : `${(bytes / 1024).toFixed(0)} KB`
}

async function submit() {
  if (!selectedFile.value && !youtubeURL.value.trim()) {
    toast.warn(t('upload.toastChooseFileOrUrl'))
    return
  }
  uploading.value = true
  try {
    const video = await store.upload(selectedFile.value, youtubeURL.value, language.value)
    emit('uploaded', video)
    toast.success(t('upload.toastUploaded'))
    router.push(`/videos/${video.id}`)
  } catch (err) {
    if (err?.status === 402) {
      toast.warn(t('upload.toastLimitReached'))
    } else if (err?.status === 401) {
      toast.error(t('upload.toastLoginRequired'))
    } else {
      toast.error(err?.message || t('upload.toastUploadFailed'))
    }
  } finally {
    uploading.value = false
  }
}
</script>
