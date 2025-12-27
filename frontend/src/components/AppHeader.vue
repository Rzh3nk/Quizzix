<template>
  <header class="app-header">
    <div class="container">
      <!-- Логотип и навигация -->
      <div class="header-left">
        <!-- Логотип - будет ссылкой на главную -->
        <div class="logo">
          <div class="logo-icon">Q</div>
          <span class="logo-text">QuizMaster</span>
          <!-- TODO: Добавить router-link когда будет страница Home -->
          <!-- <router-link to="/" class="logo-link"> -->
          <!--   <div class="logo-icon">Q</div> -->
          <!--   <span class="logo-text">QuizMaster</span> -->
          <!-- </router-link> -->
        </div>
        
        <!-- Основная навигация -->
        <nav class="main-nav">
          <!-- Главная -->
          <div class="nav-link">
            <div class="nav-icon">🏠</div>
            <span>Главная</span>
            <!-- TODO: Добавить router-link когда будет страница Home -->
            <!-- <router-link to="/" class="nav-link-content"> -->
            <!--   <div class="nav-icon">🏠</div> -->
            <!--   <span>Главная</span> -->
            <!-- </router-link> -->
          </div>
          
          <!-- Категории квизов -->
          <div class="nav-link">
            <div class="nav-icon">📁</div>
            <span>Категории</span>
            <!-- TODO: Добавить router-link когда будет страница Categories -->
            <!-- <router-link to="/categories" class="nav-link-content"> -->
            <!--   <div class="nav-icon">📁</div> -->
            <!--   <span>Категории</span> -->
            <!-- </router-link> -->
          </div>
          
          <!-- Создать квиз -->
          <div class="nav-link">
            <div class="nav-icon">➕</div>
            <span>Создать</span>
            <!-- TODO: Добавить router-link когда будет страница CreateQuiz -->
            <!-- <router-link to="/create-quiz" class="nav-link-content"> -->
            <!--   <div class="nav-icon">➕</div> -->
            <!--   <span>Создать</span> -->
            <!-- </router-link> -->
          </div>
          
          <!-- Лидерборд -->
          <div class="nav-link">
            <div class="nav-icon">🏆</div>
            <span>Рейтинг</span>
            <!-- TODO: Добавить router-link когда будет страница Leaderboard -->
            <!-- <router-link to="/leaderboard" class="nav-link-content"> -->
            <!--   <div class="nav-icon">🏆</div> -->
            <!--   <span>Рейтинг</span> -->
            <!-- </router-link> -->
          </div>
        </nav>
      </div>
      
      <!-- Правая часть: поиск и профиль -->
      <div class="header-right">
        <!-- Поиск квизов -->
        <div class="search-container">
          <div class="search-icon">🔍</div>
          <input 
            type="text" 
            placeholder="Поиск квизов..." 
            class="search-input"
            @input="handleSearch"
          />
          <!-- TODO: Добавить функционал поиска когда будет API -->
        </div>
        
        <!-- Уведомления -->
        <button 
          class="notification-btn"
          @click="showNotifications"
        >
          <div class="nav-icon">🔔</div>
          <span v-if="unreadCount > 0" class="notification-badge">
            {{ unreadCount }}
          </span>
          <!-- TODO: Добавить страницу уведомлений -->
          <!-- <router-link to="/notifications" class="notification-link">...</router-link> -->
        </button>
        
        <!-- Профиль пользователя -->
        <div class="user-menu">
          <button 
            class="user-btn"
            @click="toggleUserMenu"
          >
            <div class="user-avatar">
              <span class="avatar-text">{{ userInitials }}</span>
            </div>
            <span class="user-name">{{ userName }}</span>
            <div class="chevron-icon" :class="{ rotated: showUserMenu }">▼</div>
          </button>
          
          <!-- Выпадающее меню профиля -->
          <div 
            v-if="showUserMenu" 
            class="dropdown-menu"
            @click.stop
          >
            <!-- Профиль -->
            <div class="dropdown-item">
              <div class="dropdown-icon">👤</div>
              <span>Мой профиль</span>
              <!-- TODO: Добавить router-link когда будет страница Profile -->
              <!-- <router-link to="/profile" class="dropdown-link"> -->
              <!--   <div class="dropdown-icon">👤</div> -->
              <!--   <span>Мой профиль</span> -->
              <!-- </router-link> -->
            </div>
            
            <!-- Мои квизы -->
            <div class="dropdown-item">
              <div class="dropdown-icon">📝</div>
              <span>Мои квизы</span>
              <!-- TODO: Добавить router-link когда будет страница MyQuizzes -->
              <!-- <router-link to="/my-quizzes" class="dropdown-link"> -->
              <!--   <div class="dropdown-icon">📝</div> -->
              <!--   <span>Мои квизы</span> -->
              <!-- </router-link> -->
            </div>
            
            <!-- Настройки -->
            <div class="dropdown-item">
              <div class="dropdown-icon">⚙️</div>
              <span>Настройки</span>
              <!-- TODO: Добавить router-link когда будет страница Settings -->
              <!-- <router-link to="/settings" class="dropdown-link"> -->
              <!--   <div class="dropdown-icon">⚙️</div> -->
              <!--   <span>Настройки</span> -->
              <!-- </router-link> -->
            </div>
            
            <div class="dropdown-divider"></div>
            
            <!-- Выйти -->
            <button 
              class="dropdown-item logout"
              @click="handleLogout"
            >
              <div class="dropdown-icon">🚪</div>
              <span>Выйти</span>
            </button>
          </div>
        </div>
        
        <!-- Мобильное меню -->
        <button 
          class="mobile-menu-btn"
          @click="toggleMobileMenu"
        >
          <div v-if="!showMobileMenu" class="menu-icon">☰</div>
          <div v-else class="menu-icon">✕</div>
        </button>
      </div>
    </div>
    
    <!-- Мобильное меню (выпадающее) -->
    <div 
      v-if="showMobileMenu" 
      class="mobile-menu"
      @click="closeMobileMenu"
    >
      <div class="mobile-menu-content" @click.stop>
        <div class="mobile-user-info">
          <div class="mobile-user-avatar">
            <span class="avatar-text">{{ userInitials }}</span>
          </div>
          <div>
            <p class="mobile-user-name">{{ userName }}</p>
            <p class="mobile-user-email">{{ userEmail }}</p>
          </div>
        </div>
        
        <div class="mobile-menu-divider"></div>
        
        <!-- Главная -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">🏠</div>
          <span>Главная</span>
          <!-- TODO: Добавить router-link когда будет страница Home -->
          <!-- <router-link to="/" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <!-- Категории -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">📁</div>
          <span>Категории</span>
          <!-- TODO: Добавить router-link когда будет страница Categories -->
          <!-- <router-link to="/categories" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <!-- Создать квиз -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">➕</div>
          <span>Создать</span>
          <!-- TODO: Добавить router-link когда будет страница CreateQuiz -->
          <!-- <router-link to="/create-quiz" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <!-- Рейтинг -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">🏆</div>
          <span>Рейтинг</span>
          <!-- TODO: Добавить router-link когда будет страница Leaderboard -->
          <!-- <router-link to="/leaderboard" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <!-- Профиль -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">👤</div>
          <span>Профиль</span>
          <!-- TODO: Добавить router-link когда будет страница Profile -->
          <!-- <router-link to="/profile" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <!-- Настройки -->
        <div class="mobile-nav-link">
          <div class="mobile-nav-icon">⚙️</div>
          <span>Настройки</span>
          <!-- TODO: Добавить router-link когда будет страница Settings -->
          <!-- <router-link to="/settings" class="mobile-nav-link-content">...</router-link> -->
        </div>
        
        <div class="mobile-menu-divider"></div>
        
        <!-- Выйти -->
        <button 
          class="mobile-nav-link logout"
          @click="handleLogout"
        >
          <div class="mobile-nav-icon">🚪</div>
          <span>Выйти</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  name: 'AppHeader',
  
  setup() {
    // Состояния
    const unreadCount = ref(3)
    const showUserMenu = ref(false)
    const showMobileMenu = ref(false)
    const searchQuery = ref('')
    
    // Моковые данные пользователя (заменить на реальные)
    const user = ref({
      name: 'Иван Петров',
      email: 'ivan@example.com',
      initials: 'ИП'
    })
    
    // Вычисляемые свойства
    const userName = computed(() => user.value.name)
    const userEmail = computed(() => user.value.email)
    const userInitials = computed(() => user.value.initials)
    
    // Методы
    const handleSearch = () => {
      console.log('Поиск:', searchQuery.value)
      // TODO: Реализовать поиск когда будет готово API
    }
    
    const showNotifications = () => {
      console.log('Показать уведомления')
      // TODO: Реализовать переход к уведомлениям
      // router.push('/notifications')
    }
    
    const toggleUserMenu = () => {
      showUserMenu.value = !showUserMenu.value
    }
    
    const closeUserMenu = () => {
      showUserMenu.value = false
    }
    
    const toggleMobileMenu = () => {
      showMobileMenu.value = !showMobileMenu.value
    }
    
    const closeMobileMenu = () => {
      showMobileMenu.value = false
    }
    
    const handleLogout = () => {
      console.log('Выход из системы')
      // TODO: Реализовать выход когда будет авторизация
      // authStore.logout()
      // router.push('/login')
    }
    
    // Закрываем меню при клике вне его области
    const handleClickOutside = (event) => {
      if (!event.target.closest('.user-menu') && showUserMenu.value) {
        showUserMenu.value = false
      }
      if (!event.target.closest('.mobile-menu-btn') && 
          !event.target.closest('.mobile-menu') && 
          showMobileMenu.value) {
        showMobileMenu.value = false
      }
    }
    
    // Вешаем обработчики
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })
    
    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })
    
    return {
      // Данные
      unreadCount,
      showUserMenu,
      showMobileMenu,
      searchQuery,
      userName,
      userEmail,
      userInitials,
      
      // Методы
      handleSearch,
      showNotifications,
      toggleUserMenu,
      closeUserMenu,
      toggleMobileMenu,
      closeMobileMenu,
      handleLogout
    }
  }
}
</script>

<style scoped>
.app-header {
  background: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  height: 70px;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* Левая часть */
.header-left {
  display: flex;
  align-items: center;
  gap: 40px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #4f46e5;
  font-weight: 700;
  font-size: 1.5rem;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: bold;
}

.logo-text {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* TODO: Раскомментировать когда добавите router-link
.logo-link {
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 10px;
}
*/

.main-nav {
  display: flex;
  gap: 24px;
}

@media (max-width: 1024px) {
  .main-nav {
    display: none;
  }
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  color: #4b5563;
  transition: all 0.3s ease;
  font-weight: 500;
}

.nav-link:hover {
  background: rgba(79, 70, 229, 0.1);
  color: #4f46e5;
}

.nav-icon {
  font-size: 1.2rem;
}

/* TODO: Раскомментировать когда добавите router-link
.nav-link-content {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: inherit;
  width: 100%;
}
*/

/* Правая часть */
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-container {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  font-size: 1.1rem;
  color: #9ca3af;
}

.search-input {
  padding: 10px 12px 10px 40px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  width: 250px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: #f9fafb;
}

.search-input:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  background: white;
}

@media (max-width: 768px) {
  .search-container {
    display: none;
  }
}

.notification-btn {
  position: relative;
  background: none;
  border: none;
  padding: 8px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.3s ease;
}

.notification-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.notification-badge {
  position: absolute;
  top: 0;
  right: 0;
  background: #ef4444;
  color: white;
  font-size: 12px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Меню пользователя */
.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  background: none;
  border: none;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.user-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
}

.user-name {
  font-weight: 500;
  font-size: 14px;
}

.chevron-icon {
  font-size: 12px;
  transition: transform 0.3s ease;
  margin-left: 4px;
}

.chevron-icon.rotated {
  transform: rotate(180deg);
}

/* Выпадающее меню */
.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  min-width: 220px;
  z-index: 1001;
  overflow: hidden;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #1f2937;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  transition: background 0.3s ease;
  font-size: 0.95rem;
}

.dropdown-item:hover {
  background: rgba(0, 0, 0, 0.05);
}

.dropdown-item.logout {
  color: #ef4444;
}

.dropdown-item.logout:hover {
  background: rgba(239, 68, 68, 0.1);
}

.dropdown-icon {
  font-size: 1.1rem;
  width: 20px;
  text-align: center;
}

.dropdown-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 8px 0;
}

/* TODO: Раскомментировать когда добавите router-link
.dropdown-link {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: inherit;
  width: 100%;
}
*/

/* Мобильное меню */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  padding: 8px;
  cursor: pointer;
  font-size: 1.5rem;
}

@media (max-width: 1024px) {
  .mobile-menu-btn {
    display: block;
  }
}

.menu-icon {
  font-size: 1.5rem;
}

/* Мобильное выпадающее меню */
.mobile-menu {
  position: fixed;
  top: 70px;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.mobile-menu-content {
  background: white;
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
  margin-bottom: 16px;
}

.mobile-user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.mobile-user-name {
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 4px;
}

.mobile-user-email {
  font-size: 14px;
  color: #6b7280;
}

.mobile-menu-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 20px 0;
}

.mobile-nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  color: #1f2937;
  border-radius: 8px;
  transition: background 0.3s ease;
  cursor: pointer;
  font-size: 1rem;
}

.mobile-nav-link:hover {
  background: rgba(0, 0, 0, 0.05);
}

.mobile-nav-link.logout {
  color: #ef4444;
}

.mobile-nav-link.logout:hover {
  background: rgba(239, 68, 68, 0.1);
}

.mobile-nav-icon {
  font-size: 1.3rem;
  width: 24px;
  text-align: center;
}

/* TODO: Раскомментировать когда добавите router-link
.mobile-nav-link-content {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: inherit;
  width: 100%;
}
*/

/* Активное состояние (когда добавите роутер) */
/*
.nav-link.active {
  background: rgba(79, 70, 229, 0.1);
  color: #4f46e5;
}

.nav-link.active .nav-icon {
  color: #4f46e5;
}
*/

/* Анимации */
@media (max-width: 480px) {
  .container {
    padding: 0 16px;
  }
  
  .logo-text {
    display: none;
  }
  
  .user-name {
    display: none;
  }
  
  .chevron-icon {
    display: none;
  }
}
</style>