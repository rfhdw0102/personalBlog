<template>
  <el-config-provider>
    <div id="app-layout">
      <header class="app-header">
        <div class="header-content">
          <div class="logo-container" @click="router.push('/')">
            <span class="logo-text">Blogs</span>
          </div>

          <nav class="nav-links">
            <el-button link @click="router.push('/')" :class="{ active: currentPath === '/' }">首页</el-button>
            <el-button link @click="router.push('/about')" :class="{ active: currentPath === '/about' }">关于</el-button>
          </nav>

          <div class="header-center">
            <el-input
              v-model="searchQuery"
              placeholder="搜索文章..."
              class="search-input"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <div class="header-right">
            <template v-if="!isLoggedIn">
              <el-button type="primary" link @click="router.push('/login')" class="nav-btn">登录</el-button>
            </template>
            <template v-else>
              <div class="nav-actions">
                <div class="notification-wrapper" @click="router.push('/management?tab=messages')">
                  <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="badge">
                    <el-icon><Bell /></el-icon>
                  </el-badge>
                </div>

                <el-dropdown trigger="click" @command="handleCommand">
                  <div class="user-info">
                    <el-avatar :size="32" :src="userAvatar" />
                    <span class="username">{{ username }}</span>
                    <el-icon><ArrowDown /></el-icon>
                  </div>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="/profile">
                        <el-icon><User /></el-icon>个人信息
                      </el-dropdown-item>
                      <el-dropdown-item command="/management">
                        <el-icon><Document /></el-icon>内容管理
                      </el-dropdown-item>
                      <el-dropdown-item v-if="isAdmin" command="/admin">
                        <el-icon><Management /></el-icon>后台管理
                      </el-dropdown-item>
                      <el-dropdown-item divided command="logout" class="logout-item">
                        <el-icon><SwitchButton /></el-icon>退出登录
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
          </div>
        </div>
      </header>

      <main class="main-content">
        <router-view />
      </main>
    </div>
  </el-config-provider>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Bell, User, Document, Management, SwitchButton, ArrowDown, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()

const isLoggedIn = ref(!!localStorage.getItem('token'))
const userStr = ref(localStorage.getItem('user'))
const unreadCount = ref(0)
const searchQuery = ref('')

const currentPath = computed(() => route.path)
const user = computed(() => userStr.value ? JSON.parse(userStr.value) : null)
const username = computed(() => user.value ? user.value.username : '')
const userAvatar = computed(() => {
  if (user.value && user.value.avatar) {
    return user.value.avatar.startsWith('http') ? user.value.avatar : 'http://localhost:8082' + user.value.avatar
  }
  return 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
})
const isAdmin = computed(() => user.value && user.value.role === 'admin')

// Update login state when storage changes (optional basic reactivity for simplicity)
window.addEventListener('storage', () => {
  isLoggedIn.value = !!localStorage.getItem('token')
  userStr.value = localStorage.getItem('user')
})
// Custom event to handle login locally
window.addEventListener('user-login', () => {
  isLoggedIn.value = !!localStorage.getItem('token')
  userStr.value = localStorage.getItem('user')
  fetchUnread()
})

// Listen for custom event to refresh unread count
window.addEventListener('refresh-unread', (event) => {
  if (typeof event.detail === 'number') {
    unreadCount.value = event.detail
  } else {
    fetchUnread()
  }
})

const fetchUnread = async () => {
  if (!localStorage.getItem('token')) {
    unreadCount.value = 0
    return
  }
  try {
    const res = await request.get('/api/notification/unread-count')
    if (res.code === 200) {
      unreadCount.value = res.data.count || 0
    }
  } catch (e) {
    unreadCount.value = 0
  }
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/', query: { q: searchQuery.value.trim() } })
  } else {
    router.push('/')
  }
}

// Watch route query for search box sync
watch(() => route.query.q, (newQ) => {
  searchQuery.value = newQ || ''
}, { immediate: true })

const handleCommand = (command) => {
  if (command === 'logout') {
    logout()
  } else {
    router.push(command)
  }
}

const logout = () => {
  request.post('/api/auth/logout').finally(() => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    isLoggedIn.value = false
    userStr.value = null
    unreadCount.value = 0
    ElMessage.success('已退出登录')
    router.push('/')
  })
}

onMounted(() => {
  fetchUnread()
})
</script>

<style>
.badge {
  transform: translate(50%, -50%);
}
.app-header {
  height: 64px;
  border-bottom: 1px solid var(--border-color);
  background-color: #fff;
  display: flex;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 1000;
}
.header-content {
  max-width: 1100px;
  width: 100%;
  margin: 0 auto;
  padding: 0 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-sizing: border-box;
}
.logo-container {
  display: flex;
  align-items: center;
  cursor: pointer;
  outline: none;
}
.logo-text {
  font-size: 24px;
  font-weight: 800;
  color: #409eff;
  letter-spacing: -0.5px;
}
.nav-links {
  display: flex;
  gap: 16px;
  margin-left: 32px;
}
.nav-links .el-button {
  font-size: 16px;
  font-weight: 500;
  color: var(--secondary-text-color);
}
.nav-links .el-button.active {
  color: var(--primary-color);
  font-weight: 600;
}
.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 40px;
}
.search-input {
  max-width: 400px;
}
.search-input {
  border-radius: 20px;
  background-color: #f1f5f9;
  box-shadow: none !important;
  border: 1px solid transparent;
  transition: all 0.3s;
}
.search-input {
  background-color: #fff;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 1px var(--primary-color) !important;
}
.header-right {
  display: flex;
  align-items: center;
}
.nav-btn {
  font-size: 20px !important;
  font-weight: 600;
}
.nav-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}
.notification-wrapper {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: var(--secondary-text-color);
  font-size: 20px;
  transition: color 0.2s;
}
.notification-wrapper:hover {
  color: var(--primary-color);
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  outline: none;
}
.username {
  font-weight: 600;
  color: var(--text-color);
  font-size: 14px;
}
.logout-item {
  color: #f56c6c !important;
}
</style>
