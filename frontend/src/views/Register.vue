<template>
  <div class="register-page">
    <!-- Фон -->
    <div class="background"></div>
    
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Логотип и название на одной строке -->
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Регистрация</h1>
      </div>
      
      <!-- Подзаголовок -->
      <p class="subtitle">Создайте аккаунт для доступа ко всем функциям</p>
      
      <!-- Форма регистрации -->
      <form class="register-form" @submit.prevent="handleSubmit">
        <!-- Имя пользователя -->
        <div class="form-group">
          <label for="username">Имя пользователя *</label>
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
        
        <!-- Email -->
        <div class="form-group">
          <label for="email">Email *</label>
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
          <div v-else class="hint">
            На этот email придёт подтверждение
          </div>
        </div>
        
        <!-- Пароль -->
        <div class="form-group">
          <label for="password">Пароль *</label>
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
            Минимум 6 символов. Пароль будет захеширован на сервере
          </div>
        </div>
        
        <!-- Подтверждение пароля -->
        <div class="form-group">
          <label for="confirmPassword">Подтвердите пароль *</label>
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
        
        <!-- Кнопка регистрации -->
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
        
        <!-- Сообщение об успехе -->
        <div v-if="successMessage" class="success-message">
          <div class="success-icon">✅</div>
          <p>{{ successMessage }}</p>
          <p class="success-subtext">
            Теперь вы можете <router-link to="/login" class="success-link">войти</router-link> в систему
          </p>
        </div>
      </form>
      
      <!-- Ссылка на вход -->
      <div class="auth-link">
        <p>Уже есть аккаунт?</p>
        <router-link to="/login" class="link">
          Войти
        </router-link>
      </div>
      
      <!-- Футер -->
      <div class="footer">
        <p>QuizMaster © 2024</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'RegisterView',
  
  setup() {
    const router = useRouter()
    
    // Состояния формы
    const form = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })
    
    // Состояния UI
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const loading = ref(false)
    const errors = reactive({})
    const successMessage = ref('')
    
    // Базовый URL API (замените на ваш)
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    
    // Методы
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
      
      // Валидация имени пользователя (по требованиям бэкенда)
      if (!form.username.trim()) {
        newErrors.username = 'Имя пользователя обязательно'
      } else if (form.username.length < 3) {
        newErrors.username = 'Минимум 3 символа'
      } else if (form.username.length > 50) {
        newErrors.username = 'Максимум 50 символов'
      } else if (!/^[a-zA-Z0-9_]+$/.test(form.username)) {
        newErrors.username = 'Только латинские буквы, цифры и подчеркивание'
      }
      
      // Валидация email
      if (!form.email.trim()) {
        newErrors.email = 'Email обязателен'
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
        newErrors.email = 'Введите корректный email'
      } else if (form.email.length > 100) {
        newErrors.email = 'Максимум 100 символов'
      }
      
      // Валидация пароля
      if (!form.password) {
        newErrors.password = 'Пароль обязателен'
      } else if (form.password.length < 6) {
        newErrors.password = 'Минимум 6 символов'
      } else if (form.password.length > 72) {
        newErrors.password = 'Максимум 72 символа'
      }
      
      // Валидация подтверждения пароля
      if (form.password !== form.confirmPassword) {
        newErrors.confirmPassword = 'Пароли не совпадают'
      }
      
      // Обновляем ошибки
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
        const response = await fetch(`${API_BASE_URL}/register`, {
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
        
        // Успешная регистрация
        successMessage.value = 'Регистрация успешна! Вы можете войти в систему.'
        
        // Очистка формы
        form.username = ''
        form.email = ''
        form.password = ''
        form.confirmPassword = ''
        
        // Автоматический переход на страницу входа через 3 секунды
        setTimeout(() => {
          router.push('/login')
        }, 3000)
        
      } catch (error) {
        console.error('Ошибка регистрации:', error)
        
        // Отображение ошибки пользователю
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
/* Базовые стили */
.register-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0;
  padding: 0;
  overflow: auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, 
    Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  min-height: 100vh;
}

/* Фон */
.background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    linear-gradient(135deg, #1e1b4b 0%, #312e81 25%, #4f46e5 50%, #7c3aed 75%, #a78bfa 100%),
    repeating-linear-gradient(
      45deg,
      transparent,
      transparent 20px,
      rgba(255, 255, 255, 0.03) 20px,
      rgba(255, 255, 255, 0.03) 40px
    );
  z-index: -1;
}

/* Центральный контейнер */
.content-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px 20px;
  max-width: 500px;
  margin: 0 auto;
  text-align: center;
}

/* Логотип и название на одной строке */
.logo-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin: 20px 0 10px 0;
  width: 100%;
}

.logo-circle {
  width: 60px;
  height: 60px;
  background: white;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.2rem;
  font-weight: 900;
  color: #4f46e5;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  flex-shrink: 0;
}

h1 {
  font-size: 2.2rem;
  color: white;
  margin: 0;
  font-weight: 700;
  text-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

/* Подзаголовок */
.subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 30px 0;
  font-weight: 400;
  text-align: center;
  width: 100%;
}

/* Форма */
.register-form {
  width: 100%;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  margin-bottom: 30px;
}

/* Группы формы */
.form-group {
  margin-bottom: 25px;
  text-align: left;
}

.form-group label {
  display: block;
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.form-group label::after {
  content: ' *';
  color: #ef4444;
}

/* Поля ввода */
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 15px;
  font-size: 1.2rem;
  color: #6b7280;
  z-index: 1;
}

input {
  width: 100%;
  padding: 14px 15px 14px 45px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: #f9fafb;
  color: #1f2937;
}

input:focus {
  outline: none;
  border-color: #4f46e5;
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

input.error {
  border-color: #ef4444;
  background: #fef2f2;
}

input.error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.password-toggle {
  position: absolute;
  right: 15px;
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #6b7280;
  padding: 5px;
  transition: color 0.3s ease;
}

.password-toggle:hover {
  color: #4f46e5;
}

/* Подсказки и ошибки */
.hint {
  font-size: 0.8rem;
  color: #6b7280;
  margin-top: 5px;
}

.error-message {
  font-size: 0.85rem;
  color: #ef4444;
  margin-top: 5px;
  font-weight: 500;
}

/* Сообщение об успехе */
.success-message {
  margin-top: 20px;
  padding: 20px;
  background: #d1fae5;
  border-radius: 12px;
  border: 2px solid #10b981;
  text-align: center;
}

.success-icon {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.success-message p {
  color: #065f46;
  margin: 0 0 10px 0;
  font-weight: 600;
}

.success-subtext {
  font-size: 0.9rem;
  color: #047857;
  font-weight: 400;
}

.success-link {
  color: #059669;
  font-weight: 600;
  text-decoration: none;
}

.success-link:hover {
  text-decoration: underline;
}

/* Кнопка отправки */
.submit-btn {
  width: 100%;
  padding: 16px;
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 10px;
}

.submit-btn:hover:not(:disabled) {
  background: #3730a3;
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Ссылка на вход */
.auth-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 30px;
}

.auth-link p {
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-size: 1rem;
}

.auth-link .link {
  color: white;
  font-weight: 600;
  text-decoration: none;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.auth-link .link:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
}

/* Футер */
.footer {
  margin-top: auto;
  padding-top: 20px;
}

.footer p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
  margin: 0;
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 20px 15px;
  }
  
  .register-form {
    padding: 25px;
  }
  
  .logo-row {
    gap: 12px;
    margin-top: 10px;
  }
  
  .logo-circle {
    width: 50px;
    height: 50px;
    font-size: 1.8rem;
  }
  
  h1 {
    font-size: 1.8rem;
  }
  
  .subtitle {
    font-size: 1rem;
    margin-bottom: 25px;
  }
  
  .form-group {
    margin-bottom: 20px;
  }
  
  input {
    padding: 12px 15px 12px 45px;
  }
}

@media (max-width: 480px) {
  .content-wrapper {
    padding: 15px 10px;
  }
  
  .register-form {
    padding: 20px;
  }
  
  .logo-row {
    flex-direction: column;
    gap: 10px;
    text-align: center;
  }
  
  h1 {
    font-size: 1.6rem;
  }
  
  .submit-btn {
    padding: 14px;
    font-size: 1rem;
  }
}

@media (max-height: 700px) {
  .content-wrapper {
    padding-top: 10px;
    padding-bottom: 10px;
  }
  
  .register-form {
    margin-bottom: 20px;
  }
  
  .auth-link {
    margin-bottom: 20px;
  }
}
</style>