<template>
  <div class="max-w-md mx-auto px-4 py-16">
    <div class="card">
      <h1 class="text-2xl font-bold mb-6 text-center">
        {{ isRegister ? t('login.createAccount') : t('login.welcomeBack') }}
      </h1>

      <form @submit.prevent="submit" class="space-y-4">
        <div v-if="isRegister">
          <label class="text-sm font-medium block mb-1">{{ t('common.name') }}</label>
          <input v-model="form.name" type="text" class="input" :placeholder="t('login.namePlaceholder')" required />
        </div>
        <div>
          <label class="text-sm font-medium block mb-1">{{ t('common.email') }}</label>
          <input v-model="form.email" type="email" class="input" :placeholder="t('login.emailPlaceholder')" required />
        </div>
        <div>
          <div class="flex items-center justify-between mb-1">
            <label class="text-sm font-medium">{{ t('common.password') }}</label>
            <button
              v-if="!isRegister"
              type="button"
              class="text-xs text-primary-600 hover:underline"
              @click="handleForgotPassword"
              :disabled="forgotLoading"
            >
              {{ t('login.forgotPassword') }}
            </button>
          </div>
          <input v-model="form.password" type="password" class="input" placeholder="••••••••" required />
        </div>

        <div v-if="error" class="text-red-500 text-sm">{{ error }}</div>

        <button type="submit" class="btn-primary w-full" :disabled="loading">
          {{ loading ? t('common.loading') : (isRegister ? t('login.createAccount') : t('common.login')) }}
        </button>
      </form>

      <p class="text-center text-sm mt-4 text-gray-500">
        {{ isRegister ? t('login.haveAccount') : t('login.noAccount') }}
        <RouterLink :to="isRegister ? '/login' : '/register'" class="text-primary-600 hover:underline ml-1">
          {{ isRegister ? t('common.login') : t('common.signUp') }}
        </RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { authAPI } from '@/api/client'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth  = useAuthStore()
const toast = useToast()

const isRegister = computed(() => route.name === 'register')
const form = ref({ name: '', email: '', password: '' })
const loading = ref(false)
const error = ref('')
const forgotLoading = ref(false)

async function submit() {
  loading.value = true
  error.value = ''
  try {
    if (isRegister.value) {
      await auth.register(form.value)
      toast.success(t('login.toastAccountCreated'))
    } else {
      await auth.login(form.value)
    }
    router.push(route.query.redirect || '/dashboard')
  } catch (err) {
    const msg = err?.message || t('login.errorGeneric')
    if (err?.status >= 500) {
      toast.error(t('login.toastServerError'))
    } else {
      error.value = msg
    }
  } finally {
    loading.value = false
  }
}

async function handleForgotPassword() {
  if (!form.value.email) {
    error.value = t('login.forgotPasswordNeedEmail')
    return
  }
  forgotLoading.value = true
  try {
    await authAPI.forgotPassword(form.value.email)
    toast.success(t('login.forgotPasswordSent'))
  } catch (err) {
    toast.error(err?.message || t('login.errorGeneric'))
  } finally {
    forgotLoading.value = false
  }
}
</script>
