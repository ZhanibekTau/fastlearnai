import { reactive } from 'vue'

const state = reactive({ items: [] })
let _id = 0

export function useToast() {
  function show(message, type = 'error', duration = 4500) {
    const id = ++_id
    state.items.push({ id, message, type })
    setTimeout(() => dismiss(id), duration)
  }

  function dismiss(id) {
    const i = state.items.findIndex(t => t.id === id)
    if (i !== -1) state.items.splice(i, 1)
  }

  return {
    items: state.items,
    error:   (msg) => show(msg, 'error'),
    success: (msg) => show(msg, 'success'),
    info:    (msg) => show(msg, 'info'),
    warn:    (msg) => show(msg, 'warn'),
    dismiss,
  }
}
