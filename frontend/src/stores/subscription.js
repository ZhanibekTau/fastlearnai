import { defineStore } from 'pinia'
import { ref } from 'vue'
import { subscriptionAPI } from '@/api/client'

export const useSubscriptionStore = defineStore('subscription', () => {
  const subscription = ref(null)
  const usage = ref(null)
  const plans = ref([])

  async function fetchMine() {
    subscription.value = await subscriptionAPI.getMine()
    usage.value = await subscriptionAPI.getUsage()
  }

  async function fetchPlans() {
    plans.value = await subscriptionAPI.getPlans()
  }

  async function getUpgradeInfo() {
    return subscriptionAPI.getUpgradeInfo()
  }

  return { subscription, usage, plans, fetchMine, fetchPlans, getUpgradeInfo }
})
