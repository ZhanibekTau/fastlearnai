import axios from 'axios'

const BASE = '/fastlearn/v1'

const api = axios.create({
  baseURL: BASE,
  timeout: 60000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

function extractMessage(data) {
  if (!data) return null
  if (typeof data === 'string') return data
  return data.message || data.error || data.detail || null
}

api.interceptors.response.use(
  (res) => res.data,
  async (error) => {
    const status = error.response?.status
    const data   = error.response?.data

    if (status === 401) {
      const refresh = localStorage.getItem('refresh_token')
      if (refresh) {
        try {
          const res = await axios.post(`${BASE}/auth/refresh`, { refresh_token: refresh })
          localStorage.setItem('access_token', res.data.access_token)
          localStorage.setItem('refresh_token', res.data.refresh_token)
          error.config.headers.Authorization = `Bearer ${res.data.access_token}`
          return api.request(error.config)
        } catch {
          localStorage.clear()
          window.location.href = '/login'
        }
      } else {
        localStorage.clear()
        window.location.href = '/login'
      }
    }

    const msg = extractMessage(data)
    const err = new Error(msg || networkMessage(error))
    err.status = status
    err.data   = data
    return Promise.reject(err)
  }
)

function networkMessage(error) {
  if (!error.response) return 'Network error — check your connection'
  const s = error.response.status
  if (s === 403) return 'Access denied'
  if (s === 404) return 'Not found'
  if (s === 429) return 'Too many requests — slow down'
  if (s >= 500) return 'Server error — try again later'
  return 'Something went wrong'
}

export const authAPI = {
  register:       (data)  => api.post('/auth/register', data),
  login:          (data)  => api.post('/auth/login', data),
  logout:         ()      => api.post('/auth/logout'),
  me:             ()      => api.get('/auth/me'),
  forgotPassword: (email) => api.post('/auth/forgot-password', { email }),
}

export const videoAPI = {
  upload: (formData, onProgress) =>
    api.post('/videos', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 300000,
      onUploadProgress: (e) => onProgress?.(Math.round((e.loaded * 100) / e.total)),
    }),
  list:       (page = 1) => api.get('/videos', { params: { page } }),
  get:        (id)       => api.get(`/videos/${id}`),
  delete:     (id)       => api.delete(`/videos/${id}`),
  status:     (id)       => api.get(`/videos/${id}`),
  regenerate: (id)       => api.post(`/videos/${id}/regenerate`),
  getPublic:  (id)       => api.get(`/videos/${id}/public`),
}

export const subscriptionAPI = {
  getPlans:       () => api.get('/subscriptions/plans'),
  getMine:        () => api.get('/subscriptions/me'),
  getUsage:       () => api.get('/subscriptions/usage'),
  getUpgradeInfo: () => api.get('/subscriptions/upgrade'),
}

export const adminAPI = {
  searchUsers: (email) => api.get('/admin/users', { params: { email } }),
  grantPlan:   (userId, plan, months) => api.post(`/admin/users/${userId}/subscription`, { plan, months }),
}

export default api
