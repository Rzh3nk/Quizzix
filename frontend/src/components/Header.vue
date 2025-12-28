<template>
  <header class="main-header">
    <div class="header-container">
      <!-- Логотип -->
      <router-link to="/" class="logo">
        <div class="logo-icon">🧠</div>
        <div class="logo-text">QuizMaster</div>
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
/* Основные стили шапки */
.main-header {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 1px 15px rgba(0, 0, 0, 0.05);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Логотип */
.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  color: inherit;
  font-weight: 700;
  font-size: 1.5rem;
  transition: transform 0.2s;
}

.logo:hover {
  transform: scale(1.02);
}

.logo-icon {
  font-size: 1.8rem;
  color: #4f46e5;
}

.logo-text {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Навигация */
.main-nav {
  display: flex;
  gap: 1rem;
  margin: 0 auto;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  text-decoration: none;
  color: #4a5568;
  border-radius: 10px;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 0.95rem;
}

.nav-link:hover {
  background: rgba(79, 70, 229, 0.08);
  color: #4f46e5;
  transform: translateY(-1px);
}

.nav-link.active {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.nav-icon {
  font-size: 1.1rem;
}

/* Правая часть */
.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

/* Для авторизованных */
.auth-section {
  position: relative;
}

.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  background: none;
  border: 1px solid rgba(79, 70, 229, 0.2);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 150px;
}

.user-btn:hover {
  background: rgba(79, 70, 229, 0.05);
  border-color: rgba(79, 70, 229, 0.4);
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.1);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
}

.user-name {
  font-weight: 500;
  color: #2d3748;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100px;
}

.dropdown-arrow {
  font-size: 0.7rem;
  color: #a0aec0;
  transition: transform 0.2s;
}

.user-btn:hover .dropdown-arrow {
  transform: rotate(180deg);
}

/* Выпадающее меню профиля */
.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 200px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(0, 0, 0, 0.08);
  z-index: 1001;
  overflow: hidden;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  text-decoration: none;
  color: #4a5568;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  font-size: 0.95rem;
  transition: all 0.2s;
}

.dropdown-item:hover {
  background: rgba(79, 70, 229, 0.05);
  color: #4f46e5;
}

.dropdown-icon {
  font-size: 1.1rem;
  width: 20px;
  text-align: center;
}

.dropdown-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.08);
  margin: 0.25rem 0;
}

.dropdown-item.logout {
  color: #f56565;
}

.dropdown-item.logout:hover {
  background: rgba(245, 101, 101, 0.05);
  color: #e53e3e;
}

/* Для неавторизованных */
.guest-section {
  display: flex;
  gap: 0.75rem;
}

.auth-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.login-btn {
  background: rgba(79, 70, 229, 0.08);
  color: #4f46e5;
  border-color: rgba(79, 70, 229, 0.2);
}

.login-btn:hover {
  background: rgba(79, 70, 229, 0.15);
  border-color: rgba(79, 70, 229, 0.4);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.2);
}

.register-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 70, 229, 0.4);
}

/* Мобильное меню */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #4a5568;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.mobile-menu-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.mobile-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.mobile-menu {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 320px;
  background: white;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  box-shadow: -5px 0 25px rgba(0, 0, 0, 0.1);
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
  }
  to {
    transform: translateX(0);
  }
}

.mobile-menu-header {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mobile-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.mobile-logo-icon {
  font-size: 1.5rem;
  color: #4f46e5;
}

.mobile-logo-text {
  font-weight: 700;
  font-size: 1.2rem;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.mobile-close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #4a5568;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.mobile-close-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.mobile-menu-content {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 0;
  display: flex;
  flex-direction: column;
}

.mobile-nav {
  margin-bottom: 2rem;
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
  text-decoration: none;
  color: #4a5568;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.mobile-nav-item:hover {
  background: rgba(79, 70, 229, 0.05);
  color: #4f46e5;
}

.mobile-nav-icon {
  font-size: 1.2rem;
  width: 24px;
  text-align: center;
}

.mobile-auth-section {
  margin-top: auto;
  padding: 1.5rem;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.mobile-user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.mobile-user-details {
  flex: 1;
}

.mobile-user-name {
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.25rem;
}

.mobile-user-email {
  font-size: 0.85rem;
  color: #718096;
}

.mobile-user-links {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.mobile-user-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  text-decoration: none;
  color: #4a5568;
  border-radius: 8px;
  transition: background-color 0.2s;
  font-size: 0.95rem;
}

.mobile-user-link:hover {
  background: rgba(79, 70, 229, 0.05);
  color: #4f46e5;
}

.mobile-user-link.logout {
  color: #f56565;
}

.mobile-user-link.logout:hover {
  background: rgba(245, 101, 101, 0.05);
  color: #e53e3e;
}

.link-icon {
  font-size: 1.1rem;
  width: 20px;
  text-align: center;
}

.mobile-auth-buttons {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.mobile-auth-btn {
  padding: 0.875rem;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 500;
  text-align: center;
  transition: all 0.3s ease;
}

.mobile-auth-btn.login {
  background: rgba(79, 70, 229, 0.08);
  color: #4f46e5;
  border: 2px solid rgba(79, 70, 229, 0.2);
}

.mobile-auth-btn.login:hover {
  background: rgba(79, 70, 229, 0.15);
}

.mobile-auth-btn.register {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

/* Адаптивность */
@media (max-width: 992px) {
  .main-nav {
    display: none;
  }
  
  .mobile-menu-btn {
    display: block;
  }
  
  .user-name {
    display: none;
  }
  
  .user-btn {
    min-width: auto;
    padding: 0.5rem;
  }
}

@media (max-width: 768px) {
  .header-container {
    padding: 0 1rem;
  }
  
  .guest-section {
    display: none;
  }
  
  .auth-btn {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .mobile-menu {
    width: 100%;
  }
}
</style>