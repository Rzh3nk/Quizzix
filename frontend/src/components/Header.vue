<template>
  <header class="main-header">
    <div class="header-container">
      <!-- Логотип -->
      <router-link to="/" class="logo">
        <div class="logo-icon">Q</div>
        <div class="logo-text">Quizzix</div>
      </router-link>

      <!-- Навигация для авторизованных пользователей -->
      <nav v-if="isAuthenticated" class="main-nav">
        <router-link to="/categories" class="nav-link" active-class="active">
          <span class="nav-icon">📚</span>
          <span>Категории</span>
        </router-link>
        
        <router-link to="/my-quizzes" class="nav-link" active-class="active">
          <span class="nav-icon">📝</span>
          <span>Мои квизы</span>
        </router-link>
        
        <router-link to="/create-quiz" class="nav-link" active-class="active">
          <span class="nav-icon">✨</span>
          <span>Создать квиз</span>
        </router-link>
        
        <router-link to="/leaderboard" class="nav-link" active-class="active">
          <span class="nav-icon">🏆</span>
          <span>Рейтинг</span>
        </router-link>
      </nav>

      <!-- Правая часть -->
      <div class="header-right">
        <!-- Для авторизованных: профиль -->
        <div v-if="isAuthenticated" class="auth-section">
          <div class="user-menu" ref="userMenuRef">
            <button class="user-btn" @click="toggleUserMenu">
              <div class="user-avatar">
                <img 
                  v-if="user?.avatar" 
                  :src="user.avatar" 
                  :alt="user.name"
                  class="avatar-image"
                />
                <div v-else class="avatar-placeholder">
                  {{ getUserInitials() }}
                </div>
              </div>
              <span class="user-name">{{ user?.name || 'Профиль' }}</span>
              <span class="dropdown-arrow">▼</span>
            </button>
            
            <!-- Выпадающее меню -->
            <div v-if="showUserMenu" class="dropdown-menu">
              <router-link to="/profile" class="dropdown-item" @click="closeMenu">
                <span class="dropdown-icon">👤</span>
                <span>Мой профиль</span>
              </router-link>
              
              <router-link to="/settings" class="dropdown-item" @click="closeMenu">
                <span class="dropdown-icon">⚙️</span>
                <span>Настройки</span>
              </router-link>
              
              <div class="dropdown-divider"></div>
              
              <button @click="logout" class="dropdown-item logout">
                <span class="dropdown-icon">🚪</span>
                <span>Выйти</span>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Для неавторизованных: кнопки входа/регистрации -->
        <div v-else class="guest-section">
          <router-link to="/login" class="auth-btn login-btn">
            <span>Войти</span>
          </router-link>
          <router-link to="/register" class="auth-btn register-btn">
            <span>Зарегистрироваться</span>
          </router-link>
        </div>
      </div>
      
      <!-- Мобильное меню (бургер) -->
      <button class="mobile-menu-btn" @click="toggleMobileMenu">
        <span class="burger-icon">☰</span>
      </button>
    </div>

    <!-- Мобильное меню -->
    <div v-if="showMobileMenu" class="mobile-menu-overlay" @click="toggleMobileMenu"></div>
    
    <div v-if="showMobileMenu" class="mobile-menu">
      <div class="mobile-menu-header">
        <div class="mobile-logo">
          <div class="mobile-logo-icon">🧠</div>
          <div class="mobile-logo-text">QuizMaster</div>
        </div>
        <button class="mobile-close-btn" @click="toggleMobileMenu">✕</button>
      </div>
      
      <div class="mobile-menu-content">
        <!-- Навигация -->
        <nav class="mobile-nav">
          <router-link to="/" class="mobile-nav-item" @click="closeMobileMenu">
            <span class="mobile-nav-icon">🏠</span>
            <span>Главная</span>
          </router-link>
          
          <router-link v-if="isAuthenticated" to="/categories" class="mobile-nav-item" @click="closeMobileMenu">
            <span class="mobile-nav-icon">📚</span>
            <span>Категории</span>
          </router-link>
          
          <router-link v-if="isAuthenticated" to="/my-quizzes" class="mobile-nav-item" @click="closeMobileMenu">
            <span class="mobile-nav-icon">📝</span>
            <span>Мои квизы</span>
          </router-link>
          
          <router-link v-if="isAuthenticated" to="/create-quiz" class="mobile-nav-item" @click="closeMobileMenu">
            <span class="mobile-nav-icon">✨</span>
            <span>Создать квиз</span>
          </router-link>
          
          <router-link v-if="isAuthenticated" to="/leaderboard" class="mobile-nav-item" @click="closeMobileMenu">
            <span class="mobile-nav-icon">🏆</span>
            <span>Рейтинг</span>
          </router-link>
        </nav>
        
        <!-- Авторизация в мобильном меню -->
        <div class="mobile-auth-section">
          <template v-if="isAuthenticated">
            <div class="mobile-user-info">
              <div class="mobile-user-avatar">
                <div v-if="user?.avatar" class="avatar-image">
                  <img :src="user.avatar" :alt="user.name" />
                </div>
                <div v-else class="avatar-placeholder">
                  {{ getUserInitials() }}
                </div>
              </div>
              <div class="mobile-user-details">
                <div class="mobile-user-name">{{ user?.name || 'Пользователь' }}</div>
                <div class="mobile-user-email">{{ user?.email || '' }}</div>
              </div>
            </div>
            
            <div class="mobile-user-links">
              <router-link to="/profile" class="mobile-user-link" @click="closeMobileMenu">
                <span class="link-icon">👤</span>
                <span>Профиль</span>
              </router-link>
              
              <router-link to="/settings" class="mobile-user-link" @click="closeMobileMenu">
                <span class="link-icon">⚙️</span>
                <span>Настройки</span>
              </router-link>
              
              <button @click="logout" class="mobile-user-link logout">
                <span class="link-icon">🚪</span>
                <span>Выйти</span>
              </button>
            </div>
          </template>
          
          <template v-else>
            <div class="mobile-auth-buttons">
              <router-link to="/login" class="mobile-auth-btn login" @click="closeMobileMenu">
                <span>Войти</span>
              </router-link>
              <router-link to="/register" class="mobile-auth-btn register" @click="closeMobileMenu">
                <span>Зарегистрироваться</span>
              </router-link>
            </div>
          </template>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth.js'  // ✅ .js обязательно!

const router = useRouter()
const authStore = useAuthStore()

// ✅ ИСПРАВЬТЕ computed:
const isAuthenticated = computed(() => {
  // auth.js возвращает isAuthenticated как функцию computed()
  return authStore.isAuthenticated.value || !!authStore.user?.value
})

const user = computed(() => {
  return authStore.user?.value || null  // ✅ .value для ref
})

// Реактивные данные
const showUserMenu = ref(false)
const showMobileMenu = ref(false)
const userMenuRef = ref(null)


// Методы
const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
}

const toggleMobileMenu = () => {
  showMobileMenu.value = !showMobileMenu.value
  // Блокируем скролл при открытом меню
  if (showMobileMenu.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const closeMenu = () => {
  showUserMenu.value = false
}

const closeMobileMenu = () => {
  showMobileMenu.value = false
  document.body.style.overflow = ''
}

const getUserInitials = () => {
  if (!user.value?.name) return 'U'
  return user.value.name
    .split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

const logout = async () => {
  try {
    await authStore.logout()
    closeMenu()
    closeMobileMenu()
    router.push('/')
  } catch (error) {
    console.error('Ошибка при выходе:', error)
  }
}

// Закрытие меню при клике вне его
const handleClickOutside = (event) => {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target)) {
    showUserMenu.value = false
  }
}

// Хуки жизненного цикла
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.body.style.overflow = '' // Восстанавливаем скролл
})
</script>

<style scoped>
@import '@/assets/header.css';
</style>