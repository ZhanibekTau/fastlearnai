<template>
  <div class="max-w-3xl mx-auto px-4 py-10">
    <h1 class="text-2xl font-bold mb-8">{{ t('admin.title') }}</h1>

    <div v-if="denied" class="card text-center text-gray-500">
      {{ t('admin.accessDenied') }}
    </div>

    <template v-else>
      <form @submit.prevent="search" class="flex gap-2 mb-6">
        <input
          v-model="query"
          type="text"
          class="input flex-1"
          :placeholder="t('admin.searchPlaceholder')"
        />
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? t('common.loading') : t('admin.searchButton') }}
        </button>
      </form>

      <p v-if="searched && !loading && results.length === 0" class="text-center text-gray-500">
        {{ t('admin.noResults') }}
      </p>

      <div v-for="r in results" :key="r.id" class="card mb-4">
        <div class="flex items-center justify-between mb-2">
          <div>
            <div class="font-semibold">{{ r.name }}</div>
            <div class="text-sm text-gray-500">{{ r.email }}</div>
          </div>
          <div class="text-right text-sm">
            <div>{{ t('admin.currentPlan') }}: <span class="font-semibold capitalize">{{ r.subscription?.plan || 'free' }}</span></div>
            <div class="text-gray-500">
              {{ r.subscription?.current_period_end
                ? t('admin.expiresOn') + ' ' + new Date(r.subscription.current_period_end).toLocaleDateString(locale)
                : t('admin.noExpiry') }}
            </div>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-2 mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
          <select v-model="forms[r.id].plan" class="input w-auto">
            <option value="free">Free</option>
            <option value="pro">Pro</option>
            <option value="premium">Premium</option>
          </select>
          <input
            v-if="forms[r.id].plan !== 'free'"
            v-model.number="forms[r.id].months"
            type="number"
            min="1"
            class="input w-24"
            :placeholder="t('admin.months')"
          />
          <button class="btn-primary text-sm" :disabled="forms[r.id].submitting" @click="grant(r)">
            {{ t('admin.grant') }}
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/client'
import { useToast } from '@/composables/useToast'

const { t, locale } = useI18n()
const toast = useToast()

const query = ref('')
const results = ref([])
const forms = reactive({})
const loading = ref(false)
const searched = ref(false)
const denied = ref(false)

function ensureForm(id) {
  if (!forms[id]) forms[id] = { plan: 'pro', months: 1, submitting: false }
}

async function search() {
  if (!query.value.trim()) return
  loading.value = true
  searched.value = true
  try {
    results.value = await adminAPI.searchUsers(query.value.trim())
    results.value.forEach((r) => ensureForm(r.id))
  } catch (err) {
    if (err?.status === 403) {
      denied.value = true
    } else {
      toast.error(err?.message || t('login.errorGeneric'))
    }
  } finally {
    loading.value = false
  }
}

async function grant(r) {
  const form = forms[r.id]
  form.submitting = true
  try {
    const sub = await adminAPI.grantPlan(r.id, form.plan, form.months)
    r.subscription = sub
    toast.success(t('admin.grantSuccess', { email: r.email }))
  } catch (err) {
    if (err?.status === 403) {
      denied.value = true
    } else {
      toast.error(err?.message || t('login.errorGeneric'))
    }
  } finally {
    form.submitting = false
  }
}
</script>
