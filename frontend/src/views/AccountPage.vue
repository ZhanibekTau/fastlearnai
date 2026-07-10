<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <h1 class="text-2xl font-bold mb-8">{{ t('account.title') }}</h1>

    <div class="card mb-6">
      <h2 class="font-semibold mb-4">{{ t('account.profile') }}</h2>
      <div class="space-y-2 text-sm">
        <div><span class="text-gray-500">{{ t('common.name') }}:</span> <span class="ml-2">{{ auth.user?.name }}</span></div>
        <div><span class="text-gray-500">{{ t('common.email') }}:</span> <span class="ml-2">{{ auth.user?.email }}</span></div>
      </div>
    </div>

    <div class="card mb-6" v-if="sub">
      <h2 class="font-semibold mb-4">{{ t('account.subscription') }}</h2>
      <div class="space-y-2 text-sm">
        <div><span class="text-gray-500">{{ t('common.plan') }}:</span> <span class="ml-2 capitalize font-medium">{{ sub.plan }}</span></div>
        <div><span class="text-gray-500">{{ t('common.status') }}:</span> <span class="ml-2 capitalize">{{ sub.status }}</span></div>
        <div v-if="sub.current_period_end">
          <span class="text-gray-500">{{ t('account.renews') }}</span>
          <span class="ml-2">{{ new Date(sub.current_period_end).toLocaleDateString(locale) }}</span>
        </div>
      </div>
      <a href="https://wa.me/77479255091" target="_blank" class="btn-primary mt-4 inline-block text-sm">
        {{ t('account.upgradeViaWhatsapp') }}
      </a>
    </div>

    <div class="card" v-if="usage">
      <h2 class="font-semibold mb-4">{{ t('account.usageThisMonth') }}</h2>
      <div class="space-y-3">
        <div>
          <div class="flex justify-between text-sm mb-1">
            <span>{{ t('account.videos') }}</span>
            <span>{{ usage.videos_count }} / {{ usage.limits?.max_videos === -1 ? '∞' : usage.limits?.max_videos }}</span>
          </div>
          <div v-if="usage.limits?.max_videos > 0" class="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              class="h-full bg-primary-500 transition-all"
              :style="{ width: Math.min(100, (usage.videos_count / usage.limits.max_videos) * 100) + '%' }"
            />
          </div>
        </div>
      </div>
    </div>

    <button class="btn-secondary mt-8 text-sm" @click="handleLogout">{{ t('common.logout') }}</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscription'

const { t, locale } = useI18n()
const auth = useAuthStore()
const subStore = useSubscriptionStore()
const router = useRouter()
const sub = ref(null)
const usage = ref(null)

onMounted(async () => {
  await subStore.fetchMine()
  sub.value = subStore.subscription
  usage.value = subStore.usage
})

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>
