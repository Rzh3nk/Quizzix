import { ref, computed } from 'vue'
import axios from 'axios'

// Настройка axios
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Интерцептор для добавления токена к запросам
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Интерцептор для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Если токен просрочен или невалиден
    if (error.response?.status === 401) {
      // Очищаем аутентификацию
      user.value = null
      token.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      // Можно добавить редирект на страницу входа
      if (window.location.pathname !== '/login') {
        window.location.href = '/login?session=expired'
      }
    }
    return Promise.reject(error)
  }
)

// Состояние (реактивные переменные)
const user = ref(JSON.parse(localStorage.getItem('user')) || null)
const token = ref(localStorage.getItem('token') || null)
const loading = ref(false)
const error = ref(null)

// Геттеры (вычисляемые свойства)
const isAuthenticated = computed(() => !!token.value)
const currentUser = computed(() => user.value)

// Методы (экшены)
const login = async (credentials) => {
  try {
    loading.value = true
    error.value = null
    
    console.log('Отправка запроса на вход...', credentials)
    
    // Реальный запрос к API
    const response = await api.post('/auth/login', credentials)
    
    console.log('Ответ от API:', response.data)
    
    const { token: authToken, user: userData } = response.data
    
    // Обновляем состояние
    token.value = authToken
    user.value = userData
    
    // Сохраняем в localStorage
    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(userData))
    
    // Добавляем токен в заголовки axios
    api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
    
    return response.data
    
  } catch (err) {
    console.error('Ошибка входа:', err)
    
    // Обработка ошибок API
    if (err.response) {
      // Сервер ответил с ошибкой
      const apiError = err.response.data?.message || 
                      err.response.data?.error || 
                      `Ошибка ${err.response.status}`
      error.value = apiError
      throw new Error(apiError)
    } else if (err.request) {
      // Запрос был сделан, но ответа нет
      error.value = 'Нет ответа от сервера. Проверьте подключение.'
      throw new Error('Нет ответа от сервера')
    } else {
      // Что-то пошло не так при настройке запроса
      error.value = 'Ошибка при отправке запроса'
      throw err
    }
  } finally {
    loading.value = false
  }
}

const register = async (userData) => {
  try {
    loading.value = true
    error.value = null
    
    console.log('Регистрация пользователя:', userData)
    
    // Реальный запрос к API
    const response = await api.post('/auth/register', userData)
    
    console.log('Ответ от API:', response.data)
    
    const { token: authToken, user: userDataFromServer } = response.data
    
    token.value = authToken
    user.value = userDataFromServer
    
    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(userDataFromServer))
    
    api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
    
    return response.data
    
  } catch (err) {
    console.error('Ошибка регистрации:', err)
    
    // Специфичные ошибки регистрации
    if (err.response?.status === 409) {
      error.value = 'Пользователь с таким email уже существует'
      throw new Error('Пользователь с таким email уже существует')
    } else if (err.response?.status === 400) {
      const validationError = err.response.data?.errors?.[0]?.msg || 
                            err.response.data?.message ||
                            'Неверные данные для регистрации'
      error.value = validationError
      throw new Error(validationError)
    } else if (err.response) {
      error.value = err.response.data?.message || `Ошибка ${err.response.status}`
      throw new Error(error.value)
    } else {
      error.value = 'Ошибка при регистрации'
      throw err
    }
  } finally {
    loading.value = false
  }
}

const logout = async () => {
  try {
    // Отправляем запрос на выход (если API поддерживает)
    await api.post('/auth/logout')
  } catch (err) {
    console.warn('Ошибка при выходе:', err)
    // Продолжаем даже если запрос не удался
  } finally {
    // Очищаем состояние независимо от результата запроса
    user.value = null
    token.value = null
    error.value = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    delete api.defaults.headers.common['Authorization']
    
    console.log('Пользователь вышел')
  }
}

const checkAuth = async () => {
  const storedToken = localStorage.getItem('token')
  const storedUser = localStorage.getItem('user')
  
  if (storedToken && storedUser) {
    /*try {
      // Проверяем валидность токена через API
      const response = await api.get('/auth/verify', {
        headers: {
          Authorization: `Bearer ${storedToken}`
        }
      })
      
      if (response.data.valid) {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
        api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
        return true
      } else {
        // Токен невалиден, очищаем
        clearAuth()
        return false
      }
    } catch (err) {
      console.warn('Ошибка проверки токена:', err)
      // Если API недоступно, используем локальные данные
      token.value = storedToken
      user.value = JSON.parse(storedUser)
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      return true
    }*/
   return true
  }
  return false
}

const refreshToken = async () => {
  try {
    const refreshToken = localStorage.getItem('refreshToken')
    if (!refreshToken) return false
    
    const response = await api.post('/auth/refresh', {
      refreshToken
    })
    
    const { token: newToken } = response.data
    
    token.value = newToken
    localStorage.setItem('token', newToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    
    return true
  } catch (err) {
    console.error('Ошибка обновления токена:', err)
    clearAuth()
    return false
  }
}

const updateProfile = async (profileData) => {
  try {
    loading.value = true
    error.value = null
    
    const response = await api.put('/auth/profile', profileData)
    
    const updatedUser = response.data.user
    user.value = updatedUser
    
    localStorage.setItem('user', JSON.stringify(updatedUser))
    
    return response.data
    
  } catch (err) {
    console.error('Ошибка обновления профиля:', err)
    
    if (err.response?.status === 401) {
      error.value = 'Сессия истекла. Пожалуйста, войдите снова.'
      clearAuth()
    } else if (err.response) {
      error.value = err.response.data?.message || `Ошибка ${err.response.status}`
    } else {
      error.value = 'Ошибка обновления профиля'
    }
    
    throw err
  } finally {
    loading.value = false
  }
}

const clearAuth = () => {
  user.value = null
  token.value = null
  error.value = null
  
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  localStorage.removeItem('refreshToken')
  
  delete api.defaults.headers.common['Authorization']
}

// Вспомогательные функции
const setAuthHeader = (authToken) => {
  api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
}

const clearAuthHeader = () => {
  delete api.defaults.headers.common['Authorization']
}

// Инициализация - проверяем аутентификацию при загрузке
checkAuth()

// Экспортируем публичный интерфейс
export const useAuthStore = () => {
  return {
    // Состояние
    user,
    token,
    loading,
    error,
    
    // Геттеры
    isAuthenticated,
    currentUser,
    
    // Действия (основные)
    login,
    register,
    logout,
    checkAuth,
    updateProfile,
    refreshToken,
    clearAuth,
    
    // Вспомогательные
    setAuthHeader,
    clearAuthHeader
  }
}

// Экспорт axios instance для использования в других местах
export { api }