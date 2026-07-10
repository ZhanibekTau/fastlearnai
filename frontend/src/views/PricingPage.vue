<template>
  <div class="max-w-4xl mx-auto px-4 py-16">
    <h1 class="text-3xl font-bold text-center mb-3">{{ t('pricing.title') }}</h1>
    <p class="text-center text-gray-500 mb-12">{{ t('pricing.subtitle') }}</p>

    <div class="grid md:grid-cols-3 gap-6">
      <div
        v-for="plan in plans"
        :key="plan.id"
        class="card flex flex-col"
        :class="{ 'ring-2 ring-primary-500': plan.id === 'pro' }"
      >
        <div v-if="plan.id === 'pro'" class="text-xs font-semibold text-primary-600 uppercase mb-2">{{ t('pricing.mostPopular') }}</div>
        <h2 class="text-xl font-bold mb-1">{{ plan.name }}</h2>
        <div class="text-3xl font-bold mb-4">
          <span v-if="plan.price === 0">{{ t('pricing.free') }}</span>
          <span v-else>{{ plan.price.toLocaleString('ru-RU') }} ₸<span class="text-base font-normal text-gray-400">{{ t('pricing.perMonth') }}</span></span>
        </div>

        <ul class="space-y-2 text-sm flex-1 mb-6">
          <li class="flex gap-2">
            <span class="text-primary-500">✓</span>
            {{ plan.videos_month === -1 ? t('common.unlimited') : plan.videos_month }} {{ t('pricing.videosPerMonth') }}
          </li>
          <li class="flex gap-2">
            <span class="text-primary-500">✓</span>
            {{ plan.max_duration === -1 ? t('common.unlimited') : plan.max_duration + ' min' }} {{ t('pricing.maxDuration') }}
          </li>
          <li v-for="feature in plan.features" :key="feature" class="flex gap-2">
            <span class="text-primary-500">✓</span>
            {{ feature }}
          </li>
        </ul>

        <a
          v-if="plan.id !== 'free'"
          href="https://wa.me/77479255091"
          target="_blank"
          class="btn-primary text-center"
        >
          {{ t('pricing.get', { name: plan.name }) }}
        </a>
        <RouterLink v-else to="/register" class="btn-secondary text-center">{{ t('pricing.getStarted') }}</RouterLink>
      </div>
    </div>

    <p class="text-center text-sm text-gray-400 mt-8">
      {{ t('pricing.footnote') }}
    </p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { subscriptionAPI } from '@/api/client'

const { t } = useI18n()
const plans = ref([])

onMounted(async () => {
  plans.value = await subscriptionAPI.getPlans()
})
</script>
