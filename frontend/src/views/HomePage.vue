<template>
  <div class="max-w-3xl mx-auto px-4 py-16">
    <div class="text-center mb-12">
      <h1 class="text-4xl font-bold mb-4">
        {{ t('home.titleLine1') }}
        <span class="text-primary-600">{{ t('home.titleLine2') }}</span>
      </h1>
      <p class="text-lg text-gray-500 dark:text-gray-400">
        {{ t('home.subtitle') }}
      </p>
    </div>

    <div class="card">
      <VideoUpload @uploaded="onUploaded" />
    </div>

    <div class="mt-12 grid grid-cols-3 gap-6 text-center">
      <div class="card">
        <div class="text-3xl mb-2">⚡</div>
        <h3 class="font-semibold mb-1">{{ t('home.featureFastTitle') }}</h3>
        <p class="text-sm text-gray-500">{{ t('home.featureFastDesc') }}</p>
      </div>
      <div class="card">
        <div class="text-3xl mb-2">🌐</div>
        <h3 class="font-semibold mb-1">{{ t('home.featureMultiTitle') }}</h3>
        <p class="text-sm text-gray-500">{{ t('home.featureMultiDesc') }}</p>
      </div>
      <div class="card">
        <div class="text-3xl mb-2">📄</div>
        <h3 class="font-semibold mb-1">{{ t('home.featureExportTitle') }}</h3>
        <p class="text-sm text-gray-500">{{ t('home.featureExportDesc') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import VideoUpload from '@/components/VideoUpload.vue'

const { t } = useI18n()

function onUploaded(video) {
  // Router push already handled in VideoUpload
}
</script>
