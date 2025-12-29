<template>
  <div class="register-page">
    <div class="background"></div>
    <div class="content-wrapper">
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Регистрация</h1>
      </div>
      
      <p class="subtitle">Создайте аккаунт для доступа ко всем функциям</p>
      
      <!--Форма-->
      <form class="register-form" @submit.prevent="handleSubmit">
    
        <div class="form-group">
          <label for="username">Имя пользователя</label>
          <div class="input-wrapper">
            <div class="input-icon">👤</div>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="Придумайте логин"
              :class="{ 'error': errors.username }"
              @input="clearError('username')"
              required
              maxlength="50"
            />
          </div>
          <div v-if="errors.username" class="error-message">
            {{ errors.username }}
          </div>
          <div v-else class="hint">
            От 3 до 50 символов (латинские буквы, цифры, подчеркивание)
          </div>
        </div>
        
        <div class="form-group">
          <label for="email">Email</label>
          <div class="input-wrapper">
            <div class="input-icon">📧</div>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="example@email.com"
              :class="{ 'error': errors.email }"
              @input="clearError('email')"
              required
              maxlength="100"
            />
          </div>
          <div v-if="errors.email" class="error-message">
            {{ errors.email }}
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
              placeholder="Не менее 6 символов"
              :class="{ 'error': errors.password }"
              @input="clearError('password')"
              required
              minlength="6"
              maxlength="72"
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
          <div v-else class="hint">
            Минимум 6 символов
          </div>
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">Подтвердите пароль</label>
          <div class="input-wrapper">
            <div class="input-icon">✓</div>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="Повторите пароль"
              :class="{ 'error': errors.confirmPassword }"
              @input="clearError('confirmPassword')"
              required
            />
            <button
              type="button"
              class="password-toggle"
              @click="toggleConfirmPasswordVisibility"
            >
              {{ showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
            </button>
          </div>
          <div v-if="errors.confirmPassword" class="error-message">
            {{ errors.confirmPassword }}
          </div>
        </div>
       
        <button 
          type="submit" 
          class="submit-btn"
          :disabled="loading"
        >
          <span v-if="!loading">Создать аккаунт</span>
          <span v-else class="loading-text">
            <span class="spinner"></span>
            Регистрация...
          </span>
        </button>
       
        <div v-if="successMessage" class="success-message">
          <div class="success-icon">✅</div>
          <p>{{ successMessage }}</p>
          <p class="success-subtext">
            Теперь вы можете <router-link to="/login" class="success-link">войти</router-link> в систему
          </p>
        </div>
      </form>
      
      <div class="auth-link">
        <p>Уже есть аккаунт?</p>
        <router-link to="/login" class="link">
          Войти
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'Register',
  
  setup() {
    const router = useRouter()
    
    const form = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })
    
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const loading = ref(false)
    const errors = reactive({})
    const successMessage = ref('')
    
  
    const togglePasswordVisibility = () => {
      showPassword.value = !showPassword.value
    }
    
    const toggleConfirmPasswordVisibility = () => {
      showConfirmPassword.value = !showConfirmPassword.value
    }
    
    const clearError = (field) => {
      if (errors[field]) {
        errors[field] = ''
      }
    }
    
    const validateForm = () => {
      const newErrors = {}
      
      if (!form.username.trim()) {
        newErrors.username = 'Имя пользователя обязательно'
      } else if (form.username.length < 3) {
        newErrors.username = 'Минимум 3 символа'
      } else if (form.username.length > 50) {
        newErrors.username = 'Максимум 50 символов'
      } else if (!/^[a-zA-Z0-9_]+$/.test(form.username)) {
        newErrors.username = 'Только латинские буквы, цифры и подчеркивание'
      }
      
      if (!form.email.trim()) {
        newErrors.email = 'Email обязателен'
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
        newErrors.email = 'Введите корректный email'
      } else if (form.email.length > 100) {
        newErrors.email = 'Максимум 100 символов'
      }
      
      if (!form.password) {
        newErrors.password = 'Пароль обязателен'
      } else if (form.password.length < 6) {
        newErrors.password = 'Минимум 6 символов'
      } else if (form.password.length > 72) {
        newErrors.password = 'Максимум 72 символа'
      }
      
      if (form.password !== form.confirmPassword) {
        newErrors.confirmPassword = 'Пароли не совпадают'
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
      
      try {
        const response = await fetch(`/api/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: form.username,
            email: form.email,
            password: form.password
          })
        })
        
        const data = await response.json()
        
        if (!response.ok) {
          // Обработка ошибок с бэкенда
          if (data.error && data.error.includes('exists')) {
            throw new Error('Пользователь с таким именем или email уже существует')
          } else if (data.error && data.error.includes('hash')) {
            throw new Error('Ошибка обработки пароля')
          } else if (data.error && data.error.includes('invalid')) {
            throw new Error('Неверные данные')
          } else {
            throw new Error(data.error || 'Ошибка регистрации')
          }
        }
        
        successMessage.value = 'Регистрация успешна! Вы можете войти в систему.'
        
        form.username = ''
        form.email = ''
        form.password = ''
        form.confirmPassword = ''
        
        //Переход на страницу входа через 3 секунды
        setTimeout(() => {
          router.push('/login')
        }, 3000)
        
      } catch (error) {
        console.error('Ошибка регистрации:', error)
        
        if (error.message.includes('уже существует')) {
          errors.username = error.message
        } else {
          errors.general = error.message || 'Ошибка регистрации. Попробуйте позже.'
          alert(error.message || 'Ошибка регистрации. Попробуйте позже.')
        }
      } finally {
        loading.value = false
      }
    }
    
    return {
      form,
      showPassword,
      showConfirmPassword,
      loading,
      errors,
      successMessage,
      togglePasswordVisibility,
      toggleConfirmPasswordVisibility,
      clearError,
      handleSubmit
    }
  }
}
</script>

<style scoped>
/* Стили только для этого компонента */
@import '@/assets/reg.css';
</style>