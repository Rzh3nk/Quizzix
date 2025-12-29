<template>
  <div class="admin-panel">
    <!-- Заголовок -->
    <div class="header">
      <h1>👑 Админ панель</h1>
      <button @click="goBack">← Назад</button>
    </div>

    <!-- Только для админа -->
    <div v-if="!isAdmin" class="access-denied">
      <h2>🚫 Доступ запрещен</h2>
      <p>Требуется роль администратора</p>
    </div>

    <!-- Загрузка -->
    <div v-else-if="loading" class="loading">Загрузка пользователей...</div>

    <!-- Ошибка -->
    <div v-else-if="error" class="error">
      {{ error }}
      <button @click="fetchUsers">Повторить</button>
    </div>

    <!-- Список пользователей -->
    <div v-else class="users-list">
      <div class="table-header">
        <span>ID</span>
        <span>Username</span>
        <span>Email</span>
        <span>Баллы</span>
        <span>Роль</span>
        <span>Действия</span>
      </div>
      
      <div 
        v-for="user in users" 
        :key="user.id"
        class="user-row"
        :class="{ 'is-admin': user.role === 'admin' }"
      >
        <span>{{ user.id }}</span>
        <span>{{ user.username }}</span>
        <span>{{ user.email }}</span>
        <span>{{ formatNumber(user.total_score || 0) }}</span>
        <span :class="user.role === 'admin' ? 'admin-role' : 'user-role'">
          {{ user.role }}
        </span>
        <button 
          @click="makeAdmin(user.id)"
          :disabled="user.role === 'admin'"
          class="admin-btn"
        >
          {{ user.role === 'admin' ? '👑 Админ' : '👑 Сделать админом' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth.js'  
const authStore = useAuthStore()
const router = useRouter()

// Данные
const users = ref([])
const loading = ref(true)
const error = ref(null)

// Проверка админа
const isAdmin = () => {
 return authStore.isAdmin.value || localStorage.getItem('role') === 'admin'
}

// Загрузка пользователей
const fetchUsers = async () => {
  if (!isAdmin()) {
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

// Сделать админом
const makeAdmin = async (targetUserId) => {
  if (!isAdmin()) return
  
  const callerId = parseInt(localStorage.getItem('user_id'))
  
  const requestBody = {
    caller_id: callerId,   // 1 (админ)
    target_id: targetUserId  // 5 (делаем админом)
  }
  
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/users/setAdmin', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)  // ✅ {caller_id:1, target_id:5}
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Ошибка сервера')
    }
    
    const user = users.value.find(u => u.id === targetUserId)
    if (user) user.role = 'admin'
    
    alert('✅ Пользователь стал администратором!')
    
  } catch (err) {
    alert('❌ ' + err.message)
  }
}

const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const goBack = () => {
  router.push('/profile')
}

// Загрузка
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.admin-panel {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h1 {
  margin: 0;
}

.back-btn {
  padding: 10px 20px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.access-denied, .loading, .error {
  text-align: center;
  padding: 40px;
}

.table-header {
  display: grid;
  grid-template-columns: 60px 200px 250px 100px 100px 180px;
  background: #f3f4f6;
  padding: 15px;
  font-weight: bold;
  border-radius: 8px 8px 0 0;
}

.user-row {
  display: grid;
  grid-template-columns: 60px 200px 250px 100px 100px 180px;
  padding: 15px;
  border-bottom: 1px solid #e5e7eb;
  align-items: center;
}

.user-row:hover {
  background: #f9fafb;
}

.user-row.is-admin {
  background: #fef3c7;
}

.admin-role {
  color: #dc2626;
  font-weight: bold;
}

.user-role {
  color: #6b7280;
}

.admin-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
}

.admin-btn:not(:disabled) {
  background: #10b981;
  color: white;
}

.admin-btn:disabled {
  background: #d1d5db;
  color: #6b7280;
  cursor: not-allowed;
}
</style>
