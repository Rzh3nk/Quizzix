<template>
  <div class="login-page">
    <div class="background"></div>
    <div class="content-wrapper">
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Вход в систему</h1>
      </div>
      
      <p class="subtitle">Войдите в свой аккаунт Quizzix</p>
      
      <!--Форма-->
      <form class="login-form" @submit.prevent="handleSubmit">

        <div class="form-group">
          <label for="username">Имя пользователя</label>
          <div class="input-wrapper">
            <div class="input-icon">👤</div>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="Введите ваш логин"
              :class="{ 'error': errors.username }"
              @input="clearError('username')"
              required
            />
          </div>
          <div v-if="errors.username" class="error-message">
            {{ errors.username }}
          </div>
        </div>
        
        <div class="form-group">
          <label for="password">Пароль</label>
          <div class="input-wrapper">
            <div class="input-icon">🔒</div>
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Введите ваш пароль"
              :class="{ 'error': errors.password }"
              @input="clearError('password')"
              required
            />
            <button
              type="button"
              class="password-toggle"
              @click="togglePasswordVisibility"
            >
              {{ showPassword ? '👁️' : '👁️‍🗨️' }}
            </button>
          </div>
          <div v-if="errors.password" class="error-message">
            {{ errors.password }}
          </div>
        </div>
        
      
        <button 
          type="submit" 
          class="submit-btn"
          :disabled="loading"
        >
          <span v-if="!loading">Войти</span>
          <span v-else class="loading-text">
            <span class="spinner"></span>
            Вход...
          </span>
        </button>
        
        <div v-if="errors.general" class="error-general">
          <div class="error-icon">❌</div>
          <p>{{ errors.general }}</p>
        </div>
        
        <div v-if="successMessage" class="success-message">
          <div class="success-icon">✅</div>
          <p>{{ successMessage }}</p>
        </div>
      </form>
      
    
      <div class="auth-link">
        <p>Нет аккаунта?</p>
        <router-link to="/register" class="link">
          Зарегистрироваться
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'Login',
  
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    // Состояния формы
    const form = reactive({
      username: '',  
      password: '',
    })
    
    // Состояния UI
    const showPassword = ref(false)
    const loading = ref(false)
    const errors = reactive({})
    const successMessage = ref('')
  

    const togglePasswordVisibility = () => {
      showPassword.value = !showPassword.value
    }
    
    const clearError = (field) => {
      if (errors[field]) {
        errors[field] = ''
      }
      if (errors.general) {
        errors.general = ''
      }
    }
    
    const validateForm = () => {
      const newErrors = {}
      
      if (!form.username.trim()) {
        newErrors.username = 'Имя пользователя обязательно'
      }
      
      if (!form.password) {
        newErrors.password = 'Пароль обязателен'
      } else if (form.password.length < 1) {
        newErrors.password = 'Введите пароль'
      }
      
      Object.keys(errors).forEach(key => {
        if (!newErrors[key]) errors[key] = ''
      })
      Object.assign(errors, newErrors)
      
      return Object.keys(newErrors).length === 0
    }
    
    const handleSubmit = async () => {
      if (!validateForm()) return
      
      loading.value = true
      successMessage.value = ''
      errors.general = ''
      
      try {
        //Запрос на бэкенд
        const response = await fetch(`/api/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: form.username, 
            password: form.password,
          })
        })
        
        const data = await response.json()
        
        if (!response.ok) {
          //Обработка ошибок с бэкенда
          if (data.error && data.error.includes('credentials')) {
            throw new Error('Неверное имя пользователя или пароль')
          } else if (data.error && data.error.includes('invalid')) {
            throw new Error('Неверные данные')
          } else {
            throw new Error(data.error || 'Ошибка входа')
          }
        }
        

        const { token, user_id, username, email, role } = data
        
        localStorage.setItem('token', token)
        localStorage.setItem('user_id', user_id)
        localStorage.setItem('role', role)
        localStorage.setItem('username', username)
        localStorage.setItem('email', email)
        
        
        authStore.token.value = token
        authStore.login({
          id: user_id,
          username: username,
          email: email,
          token: token,
          role: role
        })
        authStore.checkAuth()
      
        successMessage.value = 'Вход выполнен успешно!'
        
        //Перенаправляем на страницу категорий через секунду
        setTimeout(() => {
          router.push('/main')
        }, 1000)
        
      } catch (error) {
        console.error('Ошибка входа:', error)
        
        errors.general = error.message || 'Ошибка входа. Проверьте данные.'
        
      } finally {
        loading.value = false
      }
    }
    
    return {
      form,
      showPassword,
      loading,
      errors,
      successMessage,
      togglePasswordVisibility,
      clearError,
      handleSubmit,
    }
  }
}
</script>

<style scoped>
/* Стили только для этого компонента */
@import '@/assets/login.css';
</style>