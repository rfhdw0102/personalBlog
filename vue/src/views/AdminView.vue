<template>
  <div class="admin">
    <div class="admin-header">
      <div class="header-content">
        <h1 class="admin-title">后台管理</h1>
        <div class="header-stats">
          <div class="stat-card">
            <div class="stat-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-body">
              <span class="stat-value">{{ statsData.user_count }}</span>
              <span class="stat-label">用户</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon articles">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-body">
              <span class="stat-value">{{ statsData.article_count }}</span>
              <span class="stat-label">文章</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon categories">
              <el-icon><Menu /></el-icon>
            </div>
            <div class="stat-body">
              <span class="stat-value">{{ statsData.category_count }}</span>
              <span class="stat-label">分类</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon tags">
              <el-icon><CollectionTag /></el-icon>
            </div>
            <div class="stat-body">
              <span class="stat-value">{{ statsData.tag_count }}</span>
              <span class="stat-label">标签</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="admin-body">
      <aside class="admin-sidebar">
        <nav class="sidebar-nav">
          <button
            v-for="item in menuItems"
            :key="item.key"
            class="nav-item"
            :class="{ active: activeMenu === item.key }"
            @click="handleSelect(item.key)"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </button>
        </nav>
      </aside>

      <main class="admin-main">
        <!-- 用户管理 -->
        <div v-if="activeMenu === 'users'" class="panel">
          <div class="panel-toolbar">
            <div class="toolbar-filters">
              <el-input v-model="userQuery" placeholder="搜索用户名或邮箱..." clearable @keyup.enter="fetchUsers" class="filter-input">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="userStatus" placeholder="全部状态" clearable class="filter-select" @change="fetchUsers">
                <el-option label="正常" value="1" />
                <el-option label="禁用" value="0" />
              </el-select>
            </div>
            <el-button type="primary" @click="openUserDialog()">
              <el-icon><Plus /></el-icon>添加用户
            </el-button>
          </div>

          <el-table :data="users" v-loading="userLoading" class="data-table">
            <el-table-column label="用户" min-width="200">
              <template #default="{ row }">
                <div class="user-cell">
                  <el-avatar :size="36" :src="row.avatar || '/uploads/avatars/default.png'" />
                  <div>
                    <div class="cell-primary">{{ row.username }}</div>
                    <div class="cell-secondary">ID: {{ row.id }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column label="角色" width="90">
              <template #default="{ row }">
                <span class="role-badge" :class="row.role">{{ row.role === 'admin' ? '管理员' : '用户' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <span class="status-dot" :class="{ on: row.status === 1 }" />
                {{ row.status === 1 ? '正常' : '禁用' }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="注册时间" width="160">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="openUserDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该用户吗？" @confirm="deleteUser(row.id)">
                  <template #reference>
                    <el-button link type="danger" size="small" :disabled="row.role === 'admin'">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="userTotal > pageSize">
            <el-pagination background layout="prev, pager, next" :total="userTotal" :page-size="pageSize" v-model:current-page="userPage" @current-change="fetchUsers" />
          </div>
        </div>

        <!-- 文章管理 -->
        <div v-if="activeMenu === 'articles'" class="panel">
          <div class="panel-toolbar">
            <div class="toolbar-filters">
              <el-input v-model="adminArticleQuery" placeholder="搜索文章..." clearable @keyup.enter="fetchAdminArticles" class="filter-input">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
              <el-select v-model="adminArticleStatus" placeholder="全部状态" clearable class="filter-select" @change="fetchAdminArticles">
                <el-option label="已发布" value="published" />
                <el-option label="草稿" value="draft" />
              </el-select>
            </div>
            <el-button type="primary" @click="router.push('/article/edit')">
              <el-icon><Edit /></el-icon>写文章
            </el-button>
          </div>

          <el-table :data="adminArticles" v-loading="adminArticleLoading" class="data-table">
            <el-table-column label="文章" min-width="360">
              <template #default="{ row }">
                <div class="article-cell">
                  <div class="thumb">
                    <img :src="row.cover_image || defaultCover" alt="cover" />
                  </div>
                  <div class="article-info">
                    <div class="cell-primary link" @click="router.push(`/article/${row.id}`)">{{ row.title }}</div>
                    <div class="cell-secondary">{{ row.summary || '暂无摘要' }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <span class="article-status-tag" :class="row.status">
                  {{ row.status === 'published' ? '已发布' : '草稿' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="浏览" width="70" prop="view_count" />
            <el-table-column label="点赞" width="70" prop="like_count" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="router.push(`/article/edit/${row.id}`)">编辑</el-button>
                <el-popconfirm title="确定删除该文章吗？" @confirm="deleteArticle(row.id)">
                  <template #reference>
                    <el-button link type="danger" size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="adminArticleTotal > pageSize">
            <el-pagination background layout="prev, pager, next" :total="adminArticleTotal" :page-size="pageSize" v-model:current-page="adminArticlePage" @current-change="fetchAdminArticles" />
          </div>
        </div>

        <!-- 分类管理 -->
        <div v-if="activeMenu === 'categories'" class="panel">
          <div class="panel-toolbar">
            <div class="toolbar-filters">
              <el-input v-model="categoryQuery" placeholder="搜索分类..." clearable @keyup.enter="fetchCategories" class="filter-input">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
            </div>
            <el-button type="primary" @click="openCategoryDialog()">
              <el-icon><Plus /></el-icon>添加分类
            </el-button>
          </div>

          <el-table :data="categories" v-loading="categoryLoading" class="data-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="updated_at" label="更新时间" width="170">
              <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="openCategoryDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该分类吗？" @confirm="deleteCategory(row.id)">
                  <template #reference>
                    <el-button link type="danger" size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="categoryTotal > pageSize">
            <el-pagination background layout="prev, pager, next" :total="categoryTotal" :page-size="pageSize" v-model:current-page="categoryPage" @current-change="fetchCategories" />
          </div>
        </div>

        <!-- 标签管理 -->
        <div v-if="activeMenu === 'tags'" class="panel">
          <div class="panel-toolbar">
            <div class="toolbar-filters">
              <el-input v-model="tagQuery" placeholder="搜索标签..." clearable @keyup.enter="fetchTags" class="filter-input">
                <template #prefix><el-icon><Search /></el-icon></template>
              </el-input>
            </div>
            <el-button type="primary" @click="openTagDialog()">
              <el-icon><Plus /></el-icon>添加标签
            </el-button>
          </div>

          <el-table :data="tags" v-loading="tagLoading" class="data-table">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="updated_at" label="更新时间" width="170">
              <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="openTagDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该标签吗？" @confirm="deleteTag(row.id)">
                  <template #reference>
                    <el-button link type="danger" size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="tagTotal > pageSize">
            <el-pagination background layout="prev, pager, next" :total="tagTotal" :page-size="pageSize" v-model:current-page="tagPage" @current-change="fetchTags" />
          </div>
        </div>
      </main>
    </div>

    <!-- 用户弹窗 -->
    <el-dialog v-model="userDialogVisible" :title="userForm.id ? '编辑用户' : '添加用户'" width="480px" class="modern-dialog">
      <el-form :model="userForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="userForm.username" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="userForm.email" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="userForm.role" class="w-100" :disabled="!!userForm.id || true">
            <el-option label="管理员" value="admin" />
            <el-option label="用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="userForm.status" class="w-100" :disabled="!userForm.id">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="userForm.password" type="password" show-password :placeholder="userForm.id ? '留空则不修改' : '请输入密码'" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="userForm.confirmPassword" type="password" show-password :placeholder="userForm.id ? '留空则不修改' : '请确认密码'" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="userSaving" @click="saveUser">保存</el-button>
      </template>
    </el-dialog>

    <!-- 分类弹窗 -->
    <el-dialog v-model="categoryDialogVisible" :title="categoryForm.id ? '编辑分类' : '添加分类'" width="420px" class="modern-dialog">
      <el-form :model="categoryForm" label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCategory">保存</el-button>
      </template>
    </el-dialog>

    <!-- 标签弹窗 -->
    <el-dialog v-model="tagDialogVisible" :title="tagForm.id ? '编辑标签' : '添加标签'" width="420px" class="modern-dialog">
      <el-form :model="tagForm" label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="tagForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTag">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Document, Menu, CollectionTag, Search, Plus, Edit } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { encryptText, getPublicKey } from '@/utils/rsa'

const route = useRoute()
const router = useRouter()
const pageSize = 10
const defaultCover = 'https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=300'

const menuItems = [
  { key: 'users', label: '用户管理', icon: markRaw(User) },
  { key: 'articles', label: '文章管理', icon: markRaw(Document) },
  { key: 'categories', label: '分类管理', icon: markRaw(Menu) },
  { key: 'tags', label: '标签管理', icon: markRaw(CollectionTag) }
]

const activeMenu = ref('users')

const statsData = ref({ user_count: 0, article_count: 0, category_count: 0, tag_count: 0 })

const fetchStats = async () => {
  try {
    const res = await request.get('/admin/stats')
    if (res.code === 200) statsData.value = res.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const syncTabFromRoute = () => {
  const tab = route.query.tab
  if (['users', 'articles', 'categories', 'tags'].includes(tab)) activeMenu.value = tab
}

const handleSelect = (key) => {
  activeMenu.value = key
  router.replace({ path: '/admin', query: { tab: key } })
}

watch(() => route.query.tab, () => syncTabFromRoute())

const users = ref([])
const userTotal = ref(0)
const userPage = ref(1)
const userQuery = ref('')
const userStatus = ref('')
const userLoading = ref(false)

const adminArticles = ref([])
const adminArticleTotal = ref(0)
const adminArticlePage = ref(1)
const adminArticleQuery = ref('')
const adminArticleStatus = ref('')
const adminArticleLoading = ref(false)

const categories = ref([])
const categoryTotal = ref(0)
const categoryPage = ref(1)
const categoryQuery = ref('')
const categoryLoading = ref(false)
const categoryDialogVisible = ref(false)
const categoryForm = reactive({ id: null, name: '' })

const tags = ref([])
const tagTotal = ref(0)
const tagPage = ref(1)
const tagQuery = ref('')
const tagLoading = ref(false)
const tagDialogVisible = ref(false)
const tagForm = reactive({ id: null, name: '' })

const userDialogVisible = ref(false)
const userSaving = ref(false)
const userForm = reactive({ id: null, username: '', email: '', role: 'user', status: 1, password: '', confirmPassword: '' })

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

const fetchUsers = async () => {
  userLoading.value = true
  try {
    const res = await request.get('/admin/list', { params: { query: userQuery.value || '', status: userStatus.value || '', page: userPage.value, pageSize } })
    if (res.code === 200) { users.value = res.data.list; userTotal.value = res.data.total }
  } finally { userLoading.value = false }
}

const openUserDialog = (row) => {
  if (row) {
    Object.assign(userForm, { id: row.id, username: row.username, email: row.email, role: row.role, status: row.status, password: '', confirmPassword: '' })
  } else {
    Object.assign(userForm, { id: null, username: '', email: '', role: 'user', status: 1, password: '', confirmPassword: '' })
  }
  userDialogVisible.value = true
}

const saveUser = async () => {
  if (!userForm.username || !userForm.email) { ElMessage.warning('用户名和邮箱不能为空'); return }
  if (userForm.password || userForm.confirmPassword) {
    if (userForm.password !== userForm.confirmPassword) { ElMessage.warning('两次密码不一致'); return }
  }
  userSaving.value = true
  try {
    let payload = { username: userForm.username, email: userForm.email, role: userForm.role, status: userForm.status }
    if (userForm.password) {
      const { keyId, publicKey } = await getPublicKey()
      payload = { ...payload, key_id: keyId, password: encryptText(publicKey, userForm.password), confirm_password: encryptText(publicKey, userForm.confirmPassword) }
    }
    let res
    if (userForm.id) {
      res = await request.put(`/admin/${userForm.id}`, payload)
    } else {
      if (!userForm.password) { ElMessage.warning('新增用户必须填写密码'); userSaving.value = false; return }
      res = await request.post('/admin', payload)
    }
    if (res.code === 200) { ElMessage.success('保存成功'); userDialogVisible.value = false; fetchUsers() }
  } catch (e) {
    // interceptor 已显示具体错误信息
  } finally { userSaving.value = false }
}

const deleteUser = async (id) => {
  const res = await request.delete(`/admin/${id}`)
  if (res.code === 200) { ElMessage.success('删除成功'); const i = users.value.findIndex(u => u.id === id); if (i !== -1) { users.value.splice(i, 1); userTotal.value-- } }
}

const fetchAdminArticles = async () => {
  adminArticleLoading.value = true
  try {
    const res = await request.get('/article/list', { params: { query: adminArticleQuery.value || '', status: adminArticleStatus.value || '', page: adminArticlePage.value, pageSize } })
    if (res.code === 200) { adminArticles.value = res.data.list; adminArticleTotal.value = res.data.total }
  } finally { adminArticleLoading.value = false }
}

const deleteArticle = async (id) => {
  const res = await request.delete(`/article/${id}`)
  if (res.code === 200) { ElMessage.success('删除成功'); const i = adminArticles.value.findIndex(a => a.id === id); if (i !== -1) { adminArticles.value.splice(i, 1); adminArticleTotal.value-- } }
}

const fetchCategories = async () => {
  categoryLoading.value = true
  try {
    const res = await request.get('/category/pageList', { params: { query: categoryQuery.value || '', page: categoryPage.value, pageSize } })
    if (res.code === 200) { categories.value = res.data.list; categoryTotal.value = res.data.total }
  } finally { categoryLoading.value = false }
}

const openCategoryDialog = (row) => { categoryForm.id = row ? row.id : null; categoryForm.name = row ? row.name : ''; categoryDialogVisible.value = true }

const saveCategory = async () => {
  if (!categoryForm.name) { ElMessage.warning('名称不能为空'); return }
  const res = categoryForm.id ? await request.put(`/category/${categoryForm.id}`, { name: categoryForm.name }) : await request.post('/category', { name: categoryForm.name })
  if (res.code === 200) { ElMessage.success('保存成功'); categoryDialogVisible.value = false; fetchCategories() }
}

const deleteCategory = async (id) => {
  const res = await request.delete(`/category/${id}`)
  if (res.code === 200) { ElMessage.success('删除成功'); const i = categories.value.findIndex(c => c.id === id); if (i !== -1) { categories.value.splice(i, 1); categoryTotal.value-- } }
}

const fetchTags = async () => {
  tagLoading.value = true
  try {
    const res = await request.get('/tag/pageList', { params: { query: tagQuery.value || '', page: tagPage.value, pageSize } })
    if (res.code === 200) { tags.value = res.data.list; tagTotal.value = res.data.total }
  } finally { tagLoading.value = false }
}

const openTagDialog = (row) => { tagForm.id = row ? row.id : null; tagForm.name = row ? row.name : ''; tagDialogVisible.value = true }

const saveTag = async () => {
  if (!tagForm.name) { ElMessage.warning('名称不能为空'); return }
  const res = tagForm.id ? await request.put(`/tag/${tagForm.id}`, { name: tagForm.name }) : await request.post('/tag', { name: tagForm.name })
  if (res.code === 200) { ElMessage.success('保存成功'); tagDialogVisible.value = false; fetchTags() }
}

const deleteTag = async (id) => {
  const res = await request.delete(`/tag/${id}`)
  if (res.code === 200) { ElMessage.success('删除成功'); const i = tags.value.findIndex(t => t.id === id); if (i !== -1) { tags.value.splice(i, 1); tagTotal.value-- } }
}

watch(userQuery, () => { userPage.value = 1; fetchUsers() })
watch(userStatus, () => { userPage.value = 1; fetchUsers() })
watch(adminArticleQuery, () => { adminArticlePage.value = 1; fetchAdminArticles() })
watch(adminArticleStatus, () => { adminArticlePage.value = 1; fetchAdminArticles() })
watch(categoryQuery, () => { categoryPage.value = 1; fetchCategories() })
watch(tagQuery, () => { tagPage.value = 1; fetchTags() })

watch(activeMenu, (val) => {
  if (val === 'users') fetchUsers()
  else if (val === 'articles') fetchAdminArticles()
  else if (val === 'categories') fetchCategories()
  else if (val === 'tags') fetchTags()
})

onMounted(() => {
  syncTabFromRoute()
  fetchStats()
  if (activeMenu.value === 'users') fetchUsers()
  else if (activeMenu.value === 'articles') fetchAdminArticles()
  else if (activeMenu.value === 'categories') fetchCategories()
  else if (activeMenu.value === 'tags') fetchTags()
})
</script>

<style scoped>
.admin {
  min-height: 100vh;
  background: #f8f9fb;
}

/* ── Header ── */
.admin-header {
  background: #fff;
  border-bottom: 1px solid #eef0f2;
  padding: 28px 32px 24px;
}
.header-content {
  max-width: 1400px;
  margin: 0 auto;
}
.admin-title {
  font-size: 22px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 24px 0;
  letter-spacing: -0.5px;
}

.header-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  border-radius: 12px;
  border: 1px solid #eef0f2;
  transition: box-shadow 0.2s;
}
.stat-card:hover {
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
}
.stat-icon {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.stat-icon.users    { background: #eff6ff; color: #3b82f6; }
.stat-icon.articles { background: #fef3c7; color: #f59e0b; }
.stat-icon.categories { background: #ecfdf5; color: #10b981; }
.stat-icon.tags    { background: #fce7f3; color: #ec4899; }

.stat-body {
  display: flex;
  flex-direction: column;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}
.stat-label {
  font-size: 13px;
  color: #99a0aa;
  margin-top: 2px;
}

/* ── Body Layout ── */
.admin-body {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  gap: 24px;
  padding: 24px 32px 40px;
}

/* ── Sidebar ── */
.admin-sidebar {
  width: 200px;
  flex-shrink: 0;
}
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: sticky;
  top: 16px;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #64748b;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  width: 100%;
  text-align: left;
}
.nav-item:hover {
  background: #f1f5f9;
  color: #334155;
}
.nav-item.active {
  background: #1a1a2e;
  color: #fff;
}
.nav-item .el-icon {
  font-size: 18px;
}

/* ── Main Panel ── */
.admin-main {
  flex: 1;
  min-width: 0;
}
.panel {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #eef0f2;
  overflow: hidden;
}

.panel-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 18px 20px;
  border-bottom: 1px solid #f1f5f9;
}
.toolbar-filters {
  display: flex;
  gap: 12px;
  flex: 1;
}
.filter-input {
  max-width: 280px;
}
.filter-select {
  width: 120px;
}

/* ── Table ── */
.data-table {
  width: 100%;
}
.data-table :deep(.el-table__header th) {
  background: #fafbfc;
  color: #64748b;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #eef0f2;
  padding: 12px 0;
}
.data-table :deep(.el-table__body td) {
  padding: 14px 0;
  border-bottom: 1px solid #f6f8fa;
}
.data-table :deep(.el-table__row:hover td) {
  background: #fafcfd;
}

/* ── Cell Content ── */
.cell-primary {
  font-weight: 600;
  font-size: 14px;
  color: #1a1a2e;
}
.cell-primary.link {
  cursor: pointer;
  transition: color 0.15s;
}
.cell-primary.link:hover {
  color: #3b82f6;
}
.cell-secondary {
  font-size: 12px;
  color: #99a0aa;
  margin-top: 2px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.article-cell {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}
.thumb {
  width: 100px;
  height: 62px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: #f1f5f9;
}
.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.article-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.article-status-tag {
  display: inline-block;
  font-size: 12px;
  font-weight: 500;
  padding: 3px 10px;
  border-radius: 6px;
  white-space: nowrap;
}
.article-status-tag.published {
  background: #ecfdf5;
  color: #059669;
}
.article-status-tag.draft {
  background: #f1f5f9;
  color: #64748b;
}

.role-badge {
  font-size: 12px;
  font-weight: 500;
  padding: 3px 10px;
  border-radius: 6px;
}
.role-badge.admin {
  background: #fef3c7;
  color: #b45309;
}
.role-badge.user {
  background: #f1f5f9;
  color: #64748b;
}

.status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-right: 6px;
  background: #d1d5db;
}
.status-dot.on {
  background: #10b981;
}

/* ── Pagination ── */
.pagination {
  display: flex;
  justify-content: center;
  padding: 18px 0;
}

/* ── Dialogs ── */
.modern-dialog :deep(.el-dialog) {
  border-radius: 16px;
}
.modern-dialog :deep(.el-dialog__header) {
  padding: 22px 24px 16px;
  font-weight: 700;
  font-size: 17px;
}
.modern-dialog :deep(.el-dialog__body) {
  padding: 8px 24px 24px;
}
.modern-dialog :deep(.el-dialog__footer) {
  padding: 0 24px 20px;
}

/* ── Utility ── */
.w-100 { width: 100%; }

/* ── Responsive ── */
@media (max-width: 1024px) {
  .header-stats { grid-template-columns: repeat(2, 1fr); }
  .admin-sidebar { width: 160px; }
}
@media (max-width: 768px) {
  .admin-header { padding: 20px 16px; }
  .admin-body { flex-direction: column; padding: 16px; }
  .admin-sidebar { width: 100%; }
  .sidebar-nav { flex-direction: row; flex-wrap: wrap; position: static; }
  .nav-item { flex: 1; min-width: 80px; justify-content: center; font-size: 12px; }
  .header-stats { grid-template-columns: repeat(2, 1fr); }
  .panel-toolbar { flex-direction: column; align-items: stretch; }
  .toolbar-filters { flex-direction: column; }
  .filter-input { max-width: none; }
  .filter-select { width: 100%; }
}
</style>
