<template>
  <div class="admin">
    <el-row :gutter="20">
      <el-col :span="5">
        <el-card shadow="never" class="sidebar-card">
          <el-menu :default-active="activeMenu" @select="handleSelect" class="sidebar-menu">
            <el-menu-item index="users">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
            <el-menu-item index="articles">
              <el-icon><Document /></el-icon>
              <span>文章管理</span>
            </el-menu-item>
            <el-menu-item index="categories">
              <el-icon><Menu /></el-icon>
              <span>分类管理</span>
            </el-menu-item>
            <el-menu-item index="tags">
              <el-icon><CollectionTag /></el-icon>
              <span>标签管理</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>

      <el-col :span="19">
        <el-card v-if="activeMenu === 'users'" shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <div class="filters">
                <el-input v-model="userQuery" placeholder="搜索用户名或邮箱..." clearable @keyup.enter="fetchUsers">
                  <template #prefix><el-icon><Search /></el-icon></template>
                </el-input>
                <el-select v-model="userStatus" placeholder="全部状态" clearable class="status-select" @change="fetchUsers">
                  <el-option label="正常" value="1" />
                  <el-option label="禁用" value="0" />
                </el-select>
              </div>
              <el-button type="primary" @click="openUserDialog()">
                <el-icon><Plus /></el-icon>添加用户
              </el-button>
            </div>
          </template>

          <el-table :data="users" v-loading="userLoading" style="width: 100%">
            <el-table-column label="用户" min-width="240">
              <template #default="{ row }">
                <div class="user-cell">
                  <el-avatar :size="32" :src="row.avatar || '/uploads/avatars/default.png'" />
                  <div class="user-body">
                    <div class="user-name">{{ row.username }}</div>
                    <div class="user-id">ID: {{ row.id }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="email" label="邮箱" min-width="220" />
            <el-table-column prop="role" label="角色" width="120">
              <template #default="{ row }">
                <el-tag :type="row.role === 'admin' ? 'warning' : 'info'">{{ row.role }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '正常' : '禁用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="注册时间" width="180">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openUserDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该用户吗？" @confirm="deleteUser(row.id)">
                  <template #reference>
                    <el-button size="small" type="danger" :disabled="row.role === 'admin'">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="userTotal > pageSize">
            <el-pagination
              background
              layout="prev, pager, next, jumper"
              :total="userTotal"
              :page-size="pageSize"
              v-model:current-page="userPage"
              :pager-count="5"
              @current-change="fetchUsers"
            />
          </div>
        </el-card>

        <el-card v-if="activeMenu === 'articles'" shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <div class="filters">
                <el-input v-model="adminArticleQuery" placeholder="搜索我的文章..." clearable @keyup.enter="fetchAdminArticles">
                  <template #prefix><el-icon><Search /></el-icon></template>
                </el-input>
                <el-select v-model="adminArticleStatus" placeholder="全部状态" clearable class="status-select" @change="fetchAdminArticles">
                  <el-option label="已发布" value="published" />
                  <el-option label="草稿" value="draft" />
                </el-select>
              </div>
              <el-button type="primary" @click="router.push('/article/edit')">
                <el-icon><Edit /></el-icon>写文章
              </el-button>
            </div>
          </template>

          <el-table :data="adminArticles" v-loading="adminArticleLoading" style="width: 100%">
            <el-table-column label="文章" min-width="420">
              <template #default="{ row }">
                <div class="article-cell">
                  <div class="mini-cover">
                    <img :src="row.cover_image || defaultCover" alt="cover" />
                  </div>
                  <div class="mini-body">
                    <div class="mini-title" @click="router.push(`/article/${row.id}`)">{{ row.title }}</div>
                    <div class="mini-sub">{{ row.summary || '暂无摘要' }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.status === 'published' ? 'success' : 'info'">
                  {{ row.status === 'published' ? '已发布' : '草稿' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="浏览" width="90">
              <template #default="{ row }">{{ row.view_count || 0 }}</template>
            </el-table-column>
            <el-table-column label="点赞" width="90">
              <template #default="{ row }">{{ row.like_count || 0 }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="router.push(`/article/edit/${row.id}`)">编辑</el-button>
                <el-popconfirm title="确定删除该文章吗？" @confirm="deleteArticle(row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="adminArticleTotal > pageSize">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="adminArticleTotal"
              :page-size="pageSize"
              v-model:current-page="adminArticlePage"
              @current-change="fetchAdminArticles"
            />
          </div>
        </el-card>

        <el-card v-if="activeMenu === 'categories'" shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <div class="filters">
                <el-input v-model="categoryQuery" placeholder="搜索分类名称..." clearable @keyup.enter="fetchCategories">
                  <template #prefix><el-icon><Search /></el-icon></template>
                </el-input>
              </div>
              <el-button type="primary" @click="openCategoryDialog()"><el-icon><Plus /></el-icon>添加分类</el-button>
            </div>
          </template>

          <el-table :data="categories" v-loading="categoryLoading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="updated_at" label="更新时间" width="180">
              <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openCategoryDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该分类吗？" @confirm="deleteCategory(row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="categoryTotal > pageSize">
            <el-pagination
              background
              layout="prev, pager, next, jumper"
              :total="categoryTotal"
              :page-size="pageSize"
              v-model:current-page="categoryPage"
              :pager-count="5"
              @current-change="fetchCategories"
            />
          </div>
        </el-card>

        <el-card v-if="activeMenu === 'tags'" shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <div class="filters">
                <el-input v-model="tagQuery" placeholder="搜索标签名称..." clearable @keyup.enter="fetchTags">
                  <template #prefix><el-icon><Search /></el-icon></template>
                </el-input>
              </div>
              <el-button type="primary" @click="openTagDialog()"><el-icon><Plus /></el-icon>添加标签</el-button>
            </div>
          </template>

          <el-table :data="tags" v-loading="tagLoading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="updated_at" label="更新时间" width="180">
              <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openTagDialog(row)">编辑</el-button>
                <el-popconfirm title="确定删除该标签吗？" @confirm="deleteTag(row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination" v-if="tagTotal > pageSize">
            <el-pagination
              background
              layout="prev, pager, next, jumper"
              :total="tagTotal"
              :page-size="pageSize"
              v-model:current-page="tagPage"
              :pager-count="5"
              @current-change="fetchTags"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="userDialogVisible" :title="userForm.id ? '编辑用户' : '添加用户'" width="520px">
      <el-form :model="userForm" label-width="90px">
        <el-form-item label="用户名">
          <el-input v-model="userForm.username" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="userForm.email" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="userForm.role" class="w-100" :disabled="!!userForm.id || true">
            <el-option label="admin" value="admin" />
            <el-option label="user" value="user" />
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

    <el-dialog v-model="categoryDialogVisible" :title="categoryForm.id ? '编辑分类' : '添加分类'" width="420px">
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCategory">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="tagDialogVisible" :title="tagForm.id ? '编辑标签' : '添加标签'" width="420px">
      <el-form :model="tagForm" label-width="80px">
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
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { encryptText, getPublicKey } from '@/utils/rsa'

const route = useRoute()
const router = useRouter()
const pageSize = 10
const defaultCover = 'https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=300'

const activeMenu = ref('users')

const syncTabFromRoute = () => {
  const tab = route.query.tab
  if (['users', 'articles', 'categories', 'tags'].includes(tab)) {
    activeMenu.value = tab
  }
}

const handleSelect = (key) => {
  activeMenu.value = key
  router.replace({ path: '/admin', query: { tab: key } })
}

watch(() => route.query.tab, () => {
  syncTabFromRoute()
})

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
const userForm = reactive({
  id: null,
  username: '',
  email: '',
  role: 'user',
  status: 1,
  password: '',
  confirmPassword: ''
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

const fetchUsers = async () => {
  userLoading.value = true
  try {
    const res = await request.get('/api/admin/list', {
      params: { query: userQuery.value || '', status: userStatus.value || '', page: userPage.value, pageSize }
    })
    if (res.code === 200) {
      users.value = res.data.list
      userTotal.value = res.data.total
    }
  } finally {
    userLoading.value = false
  }
}

const openUserDialog = (row) => {
  if (row) {
    userForm.id = row.id
    userForm.username = row.username
    userForm.email = row.email
    userForm.role = row.role
    userForm.status = row.status
  } else {
    userForm.id = null
    userForm.username = ''
    userForm.email = ''
    userForm.role = 'user'
    userForm.status = 1
  }
  userForm.password = ''
  userForm.confirmPassword = ''
  userDialogVisible.value = true
}

const saveUser = async () => {
  if (!userForm.username || !userForm.email) {
    ElMessage.warning('用户名和邮箱不能为空')
    return
  }
  if (userForm.password || userForm.confirmPassword) {
    if (userForm.password !== userForm.confirmPassword) {
      ElMessage.warning('两次密码不一致')
      return
    }
  }

  userSaving.value = true
  try {
    let payload = {
      username: userForm.username,
      email: userForm.email,
      role: userForm.role,
      status: userForm.status
    }
    if (userForm.password) {
      const { keyId, publicKey } = await getPublicKey()
      payload = {
        ...payload,
        key_id: keyId,
        password: encryptText(publicKey, userForm.password),
        confirm_password: encryptText(publicKey, userForm.confirmPassword)
      }
    }

    let res
    if (userForm.id) {
      res = await request.put(`/api/admin/${userForm.id}`, payload)
    } else {
      if (!userForm.password) {
        ElMessage.warning('新增用户必须填写密码')
        userSaving.value = false
        return
      }
      res = await request.post('/api/admin', payload)
    }
    if (res.code === 200) {
      ElMessage.success('保存成功')
      userDialogVisible.value = false
      fetchUsers()
    }
  } finally {
    userSaving.value = false
  }
}

const deleteUser = async (id) => {
  const res = await request.delete(`/api/admin/${id}`)
  if (res.code === 200) {
    ElMessage.success('删除成功')
    const index = users.value.findIndex(u => u.id === id)
    if (index !== -1) {
      users.value.splice(index, 1)
      userTotal.value--
    }
  }
}

const fetchAdminArticles = async () => {
  adminArticleLoading.value = true
  try {
    const res = await request.get('/api/article/list', {
      params: {
        query: adminArticleQuery.value || '',
        status: adminArticleStatus.value || '',
        page: adminArticlePage.value,
        pageSize
      }
    })
    if (res.code === 200) {
      adminArticles.value = res.data.list
      adminArticleTotal.value = res.data.total
    }
  } finally {
    adminArticleLoading.value = false
  }
}

const deleteArticle = async (id) => {
  const res = await request.delete(`/api/article/${id}`)
  if (res.code === 200) {
    ElMessage.success('删除成功')
    const index = adminArticles.value.findIndex(a => a.id === id)
    if (index !== -1) {
      adminArticles.value.splice(index, 1)
      adminArticleTotal.value--
    }
  }
}

const fetchCategories = async () => {
  categoryLoading.value = true
  try {
    const res = await request.get('/api/category/pageList', { 
      params: { 
        query: categoryQuery.value || '',
        page: categoryPage.value, 
        pageSize 
      } 
    })
    if (res.code === 200) {
      categories.value = res.data.list
      categoryTotal.value = res.data.total
    }
  } finally {
    categoryLoading.value = false
  }
}

const openCategoryDialog = (row) => {
  categoryForm.id = row ? row.id : null
  categoryForm.name = row ? row.name : ''
  categoryDialogVisible.value = true
}

const saveCategory = async () => {
  if (!categoryForm.name) {
    ElMessage.warning('名称不能为空')
    return
  }
  let res
  if (categoryForm.id) {
    res = await request.put(`/api/category/${categoryForm.id}`, { name: categoryForm.name })
  } else {
    res = await request.post('/api/category', { name: categoryForm.name })
  }
  if (res.code === 200) {
    ElMessage.success('保存成功')
    categoryDialogVisible.value = false
    fetchCategories()
  }
}

const deleteCategory = async (id) => {
  const res = await request.delete(`/api/category/${id}`)
  if (res.code === 200) {
    ElMessage.success('删除成功')
    const index = categories.value.findIndex(c => c.id === id)
    if (index !== -1) {
      categories.value.splice(index, 1)
      categoryTotal.value--
    }
  }
}

const fetchTags = async () => {
  tagLoading.value = true
  try {
    const res = await request.get('/api/tag/pageList', { 
      params: { 
        query: tagQuery.value || '',
        page: tagPage.value, 
        pageSize 
      } 
    })
    if (res.code === 200) {
      tags.value = res.data.list
      tagTotal.value = res.data.total
    }
  } finally {
    tagLoading.value = false
  }
}

const openTagDialog = (row) => {
  tagForm.id = row ? row.id : null
  tagForm.name = row ? row.name : ''
  tagDialogVisible.value = true
}

const saveTag = async () => {
  if (!tagForm.name) {
    ElMessage.warning('名称不能为空')
    return
  }
  let res
  if (tagForm.id) {
    res = await request.put(`/api/tag/${tagForm.id}`, { name: tagForm.name })
  } else {
    res = await request.post('/api/tag', { name: tagForm.name })
  }
  if (res.code === 200) {
    ElMessage.success('保存成功')
    tagDialogVisible.value = false
    fetchTags()
  }
}

const deleteTag = async (id) => {
  const res = await request.delete(`/api/tag/${id}`)
  if (res.code === 200) {
    ElMessage.success('删除成功')
    const index = tags.value.findIndex(t => t.id === id)
    if (index !== -1) {
      tags.value.splice(index, 1)
      tagTotal.value--
    }
  }
}

watch(userQuery, () => { userPage.value = 1; fetchUsers() })
watch(userStatus, () => { userPage.value = 1; fetchUsers() })
watch(adminArticleQuery, () => { adminArticlePage.value = 1; fetchAdminArticles() })
watch(adminArticleStatus, () => { adminArticlePage.value = 1; fetchAdminArticles() })
watch(categoryQuery, () => { categoryPage.value = 1; fetchCategories() })
watch(tagQuery, () => { tagPage.value = 1; fetchTags() })

watch(activeMenu, (newVal) => {
  if (newVal === 'users') fetchUsers()
  else if (newVal === 'articles') fetchAdminArticles()
  else if (newVal === 'categories') fetchCategories()
  else if (newVal === 'tags') fetchTags()
})

onMounted(() => {
  syncTabFromRoute()
  if (activeMenu.value === 'users') fetchUsers()
  else if (activeMenu.value === 'articles') fetchAdminArticles()
  else if (activeMenu.value === 'categories') fetchCategories()
  else if (activeMenu.value === 'tags') fetchTags()
})
</script>

<style scoped>
.admin {
  padding: 20px;
  background: #f5f7fb;
  min-height: 100vh;
}

/* 卡片 */
.sidebar-card,
.panel-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* 顶部渐变线 */
.panel-card::before {
  content: '';
  display: block;
  height: 3px;
  background: linear-gradient(90deg, #409eff, #67c23a);
}

/* 侧边栏 */
.sidebar-menu {
  border: none;
}

.sidebar-menu .el-menu-item {
  border-radius: 10px;
  margin: 4px 0;
  transition: all 0.2s;
}

.sidebar-menu .el-menu-item:hover {
  background: #f0f7ff;
}

.sidebar-menu .el-menu-item.is-active {
  background: linear-gradient(135deg, #409eff, #66b1ff);
  color: #fff;
}

/* 头部 */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.filters {
  display: flex;
  gap: 12px;
  flex: 1;
  max-width: 720px;
}

.filters .el-input,
.filters .el-select {
  border-radius: 10px;
}

/* 按钮 */
.el-button--primary {
  border-radius: 10px;
  padding: 10px 18px;
  font-weight: 600;
}

/* 表格 */
.el-table {
  border-radius: 12px;
  overflow: hidden;
}

.el-table th {
  background: #f9fafb;
  font-weight: 600;
  color: #374151;
}

.el-table__row:hover td {
  background: #f5f7ff !important;
}

/* 标签 */
.el-tag {
  border-radius: 8px;
  padding: 0 10px;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.el-pagination.is-background .el-pager li {
  border-radius: 8px;
}

/* 用户 */
.user-cell {
  display: flex;
  gap: 12px;
  align-items: center;
}

.user-cell:hover {
  transform: translateX(2px);
  transition: 0.2s;
}

.user-name {
  font-weight: 600;
  font-size: 14px;
  color: #111827;
}

.user-id {
  font-size: 12px;
  color: #9ca3af;
}

/* 文章 */
.article-cell {
  display: flex;
  gap: 12px;
  align-items: center;
}

.mini-cover {
  width: 90px;
  height: 56px;
  border-radius: 10px;
  overflow: hidden;
  background: #f3f4f6;
  border: 1px solid #eef2f7;
  flex-shrink: 0;
}

.mini-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mini-body {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mini-title {
  font-weight: 600;
  font-size: 15px;
  color: #111827;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.mini-title:hover {
  color: #409EFF;
  transform: translateX(2px);
  transition: 0.2s;
}

.mini-sub {
  font-size: 12px;
  color: #9aa0a6;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 弹窗 */
.el-dialog {
  border-radius: 16px;
}

.el-dialog__header {
  font-weight: 600;
  font-size: 16px;
}

/* 通用 */
.w-100 {
  width: 100%;
}
</style>
