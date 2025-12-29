<template>
  <div class="admin-panel-page">
    <!-- Фон -->
    <div class="background"></div>
    <Header/>
    
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Хлебные крошки -->
      <div class="breadcrumbs">
        <button @click="goBack" class="breadcrumb-back">
          <span class="back-icon">←</span>
          <span>Назад</span>
        </button>
        <span class="breadcrumb-separator">›</span>
        <span class="breadcrumb-current">Админ панель</span>
      </div>

      <!-- Заголовок -->
      <div class="admin-header">
        <div class="admin-title">
          <div class="admin-icon">⚙️</div>
          <h1>Админ панель</h1>
        </div>
        <div class="admin-stats">
          <div class="stat-card">
            <div class="stat-icon">👥</div>
            <div class="stat-content">
              <div class="stat-value">{{ users.length }}</div>
              <div class="stat-label">Пользователей</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">🔧</div>
            <div class="stat-content">
              <div class="stat-value">{{ adminCount }}</div>
              <div class="stat-label">Админов</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Поиск -->
      <div class="admin-search">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск по имени или email..."
            @input="handleSearch"
          />
          <button v-if="searchQuery" @click="clearSearch" class="clear-search-btn">
            ✕
          </button>
        </div>
      </div>

      <!-- Только для админа -->
      <div v-if="!isAdmin" class="access-denied">
        <div class="denied-content">
          <div class="denied-icon">🚫</div>
          <h2>Доступ запрещен</h2>
          <p>Для просмотра этой страницы требуется роль администратора</p>
          <router-link to="/" class="back-home-btn">
            <span>Вернуться на главную</span>
          </router-link>
        </div>
      </div>

      <!-- Лоадер -->
      <div v-else-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Загрузка пользователей...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error-state">
        <div class="error-content">
          <span class="error-icon">⚠️</span>
          <h3>Ошибка загрузки</h3>
          <p>{{ error }}</p>
          <button @click="fetchUsers" class="retry-btn">
            Попробовать снова
          </button>
        </div>
      </div>

      <!-- Таблица пользователей -->
      <div v-else class="users-table">
        <div class="table-header">
          <div class="table-cell table-cell-id">ID</div>
          <div class="table-cell table-cell-name">Имя пользователя</div>
          <div class="table-cell table-cell-email">Email</div>
          <div class="table-cell table-cell-role">Роль</div>
          <div class="table-cell table-cell-actions">Действия</div>
        </div>

        <div 
          v-for="user in filteredUser" 
          :key="user.id"
          class="user-row"
          :class="{ 'admin-row': user.role === 'admin' }"
        >
          <div class="table-cell table-cell-id">
            <span class="cell-label">ID</span>
            <span class="cell-value">{{ user.id }}</span>
          </div>
          
          <div class="table-cell table-cell-name">
            <span class="cell-label">Имя</span>
            <div class="user-info">
              <span class="cell-value">{{ user.username }}</span>
            </div>
          </div>
          
          <div class="table-cell table-cell-email">
            <span class="cell-label">Email</span>
            <span class="cell-value">{{ user.email }}</span>
          </div>
          
          <div class="table-cell table-cell-role">
            <span class="cell-label">Роль</span>
            <span 
              class="role-badge" 
              :class="user.role === 'admin' ? 'role-admin' : 'role-user'"
            >
              {{ user.role === 'admin' ? 'Админ' : 'Пользователь' }}
            </span>
          </div>
          
          <div class="table-cell table-cell-actions">
            <span class="cell-label">Действия</span>
            <button 
              @click="makeAdmin(user.id)"
              :disabled="user.role === 'admin'"
              class="action-btn"
              :class="{ 'btn-disabled': user.role === 'admin', 'btn-active': user.role !== 'admin' }"
            >
              <span class="btn-text">
                {{ user.role === 'admin' ? 'Администратор' : 'Назначить админом' }}
              </span>
            </button>
          </div>
        </div>

        <!-- Сообщение если пользователей нет -->
        <div v-if="!loading && filteredUser.length === 0" class="empty-state">
          <span class="empty-icon">👥</span>
          <h3>Пользователи не найдены</h3>
          <p v-if="searchQuery">
            Попробуйте изменить параметры поиска
          </p>
          <button @click="clearSearch" class="clear-filters-btn">
            Сбросить фильтры
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth.js'
import Header from '@/components/Header.vue'

const authStore = useAuthStore()
const router = useRouter()
const searchQuery = ref('')

// Данные
const users = ref([])
const loading = ref(true)
const error = ref(null)

// Проверка админа
const isAdmin = computed(() => {
  return authStore.isAdmin.value || localStorage.getItem('role') === 'admin'
})

// Подсчет админов
const adminCount = computed(() => {
  return users.value.filter(user => user.role === 'admin').length
})

// Загрузка пользователей
const fetchUsers = async () => {
  if (!isAdmin.value) {
    router.push('/profile')
    return
  }
  
  try {
    loading.value = true
    error.value = null
    
    const token = localStorage.getItem('token')
    const response = await fetch('/api/users', {
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) throw new Error('Ошибка загрузки пользователей')
    
    const data = await response.json()
    users.value = Array.isArray(data) ? data : data.users || []
    
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const filteredUser = computed(() => {
  if (!searchQuery.value) return users.value
  
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user => 
    user.username.toLowerCase().includes(query) ||
    (user.email && user.email.toLowerCase().includes(query))
  )
})

const handleSearch = () => {
  // Debounce можно добавить при необходимости
}

const clearSearch = () => {
  searchQuery.value = ''
}

// Сделать админом
const makeAdmin = async (targetUserId) => {
  if (!isAdmin.value) return
  
  const callerId = parseInt(localStorage.getItem('user_id'))
  
  const requestBody = {
    caller_id: callerId,
    target_id: targetUserId
  }
  
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/users/setAdmin', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Ошибка сервера')
    }
    
    // Обновляем локально
    const user = users.value.find(u => u.id === targetUserId)
    if (user) user.role = 'admin'
    
    alert('✅ Пользователь стал администратором!')
    
  } catch (err) {
    alert('❌ ' + err.message)
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/main')
  }
}

// Загрузка
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
/* Основные контейнеры */
.admin-panel-page {
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

.content-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 30px 20px;
  max-width: 1400px;
  margin: 0 auto;
}

/* Хлебные крошки */
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 30px;
  color: white;
}

.breadcrumb-back {
  display: flex;
  align-items: center;
  gap: 5px;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  color: white;
  padding: 8px 16px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.95rem;
}

.breadcrumb-back:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(-3px);
}

.breadcrumb-separator {
  color: rgba(255, 255, 255, 0.5);
}

.breadcrumb-current {
  color: white;
  font-weight: 500;
}

/* Заголовок админ панели */
.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.admin-title {
  display: flex;
  align-items: center;
  gap: 15px;
}

.admin-icon {
  font-size: 2.5rem;
}

.admin-title h1 {
  font-size: 2.2rem;
  color: white;
  margin: 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.admin-stats {
  display: flex;
  gap: 15px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  min-width: 160px;
}

.stat-icon {
  font-size: 1.8rem;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
  color: white;
  line-height: 1;
}

.stat-label {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4px;
}

/* Поиск */
.admin-search {
  margin-bottom: 25px;
}

.search-box {
  position: relative;
  width: 100%;
  max-width: 500px;
}

.search-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  font-size: 1.1rem;
}

.search-box input {
  width: 100%;
  padding: 14px 20px 14px 45px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  font-size: 1rem;
  background: rgba(255, 255, 255, 0.9);
  color: #1f2937;
  transition: all 0.3s ease;
}

.search-box input:focus {
  outline: none;
  border-color: #4f46e5;
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2);
}

.clear-search-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 4px;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.clear-search-btn:hover {
  background: #f3f4f6;
  color: #4b5563;
}

/* Отказано в доступе */
.access-denied {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.denied-content {
  text-align: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 40px;
  border-radius: 20px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  max-width: 500px;
  width: 100%;
}

.denied-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.denied-content h2 {
  color: white;
  font-size: 1.8rem;
  margin-bottom: 10px;
}

.denied-content p {
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 25px;
  font-size: 1rem;
}

.back-home-btn {
  display: inline-block;
  padding: 12px 30px;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.3);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.back-home-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* Лоадер */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 60px 0;
}

.loading-spinner {
  width: 60px;
  height: 60px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: white;
  font-size: 1.1rem;
}

/* Ошибка */
.error-state {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.error-content {
  text-align: center;
  background: rgba(239, 68, 68, 0.1);
  backdrop-filter: blur(10px);
  padding: 40px;
  border-radius: 20px;
  border: 2px solid rgba(239, 68, 68, 0.3);
  max-width: 500px;
  width: 100%;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.error-content h3 {
  color: white;
  font-size: 1.5rem;
  margin-bottom: 10px;
}

.error-content p {
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 20px;
}

.retry-btn {
  padding: 10px 25px;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.3);
  color: white;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* Таблица пользователей */
.users-table {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  margin-bottom: 40px;
}

.table-header {
  display: grid;
  grid-template-columns: 80px minmax(200px, 1fr) minmax(250px, 1.5fr) 120px 200px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  padding: 20px;
  color: white;
  font-weight: 600;
  font-size: 0.95rem;
  gap: 20px;
  align-items: center;
}

.table-cell {
  display: flex;
  align-items: center;
  min-height: 44px;
  overflow: hidden;
}

.table-cell-id {
  width: 80px;
  min-width: 80px;
  max-width: 80px;
  justify-content: center;
}

.table-cell-name {
  min-width: 200px;
  flex: 1;
}

.table-cell-email {
  min-width: 250px;
  flex: 1.5;
}

.table-cell-role {
  width: 120px;
  min-width: 120px;
  max-width: 120px;
  justify-content: center;
}

.table-cell-actions {
  width: 200px;
  min-width: 200px;
  max-width: 200px;
  justify-content: center;
}

.user-row {
  display: grid;
  grid-template-columns: 80px minmax(200px, 1fr) minmax(250px, 1.5fr) 120px 200px;
  padding: 18px 20px;
  border-bottom: 1px solid #e5e7eb;
  transition: all 0.3s ease;
  gap: 20px;
  align-items: center;
}

.user-row:hover {
  background: #f9fafb;
}

.admin-row {
  background: rgba(79, 70, 229, 0.05);
}

.cell-label {
  display: none;
  font-weight: 500;
  color: #6b7280;
  font-size: 0.85rem;
  margin-right: 8px;
  min-width: 50px;
}

.cell-value {
  font-size: 0.95rem;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.table-cell-id .cell-value {
  font-weight: 600;
  color: #4f46e5;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.table-cell-email .cell-value {
  color: #6b7280;
}

.role-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  white-space: nowrap;
  text-align: center;
  width: 100px;
}

.role-admin {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.role-user {
  background: #f3f4f6;
  color: #6b7280;
}

.action-btn {
  padding: 10px 16px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9rem;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  width: 180px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-active {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.btn-active:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.4);
}

.btn-disabled {
  background: linear-gradient(135deg, #6b7280 0%, #9ca3af 100%);
  color: white;
  cursor: not-allowed;
  opacity: 0.8;
}

.btn-text {
  display: block;
  text-align: center;
  width: 100%;
}

/* Сообщение если пользователей нет */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  color: #9ca3af;
}

.empty-state h3 {
  color: #1f2937;
  font-size: 1.5rem;
  margin: 0;
}

.empty-state p {
  color: #6b7280;
  margin: 0 0 20px 0;
}

.clear-filters-btn {
  padding: 10px 25px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.clear-filters-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(79, 70, 229, 0.3);
}

/* Адаптивность */
@media (max-width: 1200px) {
  .table-header,
  .user-row {
    grid-template-columns: 70px minmax(180px, 1fr) minmax(200px, 1fr) 110px 180px;
    gap: 15px;
  }
  
  .table-cell-id,
  .table-cell-role,
  .table-cell-actions {
    width: auto;
    min-width: auto;
    max-width: none;
  }
  
  .action-btn {
    width: 160px;
    padding: 8px 12px;
    font-size: 0.85rem;
  }
}

@media (max-width: 992px) {
  .content-wrapper {
    padding: 20px;
  }
  
  .admin-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .admin-stats {
    width: 100%;
    justify-content: flex-start;
  }
  
  .stat-card {
    min-width: 140px;
  }
  
  .table-header {
    display: none;
  }
  
  .user-row {
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding: 20px;
    border-bottom: 2px solid #e5e7eb;
  }
  
  .table-cell {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    padding: 5px 0;
    border-bottom: 1px solid #f3f4f6;
  }
  
  .table-cell:last-child {
    border-bottom: none;
  }
  
  .cell-label {
    display: block;
  }
  
  .user-info {
    justify-content: flex-end;
  }
  
  .role-badge,
  .action-btn {
    width: auto;
    min-width: 140px;
  }
  
  .action-btn {
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 15px;
  }
  
  .admin-title h1 {
    font-size: 1.8rem;
  }
  
  .stat-card {
    min-width: 120px;
    padding: 12px 15px;
  }
  
  .stat-value {
    font-size: 1.5rem;
  }
  
  .search-box {
    max-width: 100%;
  }
}

@media (max-width: 480px) {
  .table-cell {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
  
  .cell-label {
    align-self: flex-start;
  }
  
  .user-info,
  .table-cell-role,
  .table-cell-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .role-badge,
  .action-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>