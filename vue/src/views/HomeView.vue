<template>
  <div class="home">
    <div class="content-wrapper">
      <!-- 左侧：排序导航 -->
      <div class="left-sidebar">
        <div class="sidebar-sticky">
          <div class="sidebar-title">浏览模式</div>
          <div class="sort-section-vertical">
            <div
                class="sort-item-v"
                :class="{ active: sortType === 0 }"
                @click="changeSort(0)"
            >
              <el-icon><Clock /></el-icon>
              <span>最近发布</span>
            </div>
            <div
                class="sort-item-v"
                :class="{ active: sortType === 1 }"
                @click="changeSort(1)"
            >
              <el-icon><TrendCharts /></el-icon>
              <span>热门推荐</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 中间：文章列表 -->
      <div class="main-wrapper">
        <el-card class="list-card" shadow="never">
          <template #header>
            <div class="category-section">
              <div class="section-label">文章分类</div>
              <div class="category-pills">
                <div
                    class="category-pill"
                    :class="{ active: activeCategory === '' }"
                    @click="changeCategory('')"
                >
                  全部
                </div>
                <div
                    v-for="cat in categories"
                    :key="cat.id"
                    class="category-pill"
                    :class="{ active: activeCategory === cat.id }"
                    @click="changeCategory(cat.id)"
                >
                  {{ cat.name }}
                </div>
              </div>
            </div>
          </template>
          <div v-if="loading" class="state-text">正在加载文章...</div>
          <div v-else-if="articles.length === 0" class="state-text">暂无文章</div>
          <div v-else>
            <div
                v-for="article in articles"
                :key="article.id"
                class="article-item"
                @click="goToDetail(article.id)"
            >
              <div class="cover">
                <img :src="article.cover_image || defaultCover" alt="cover" />
              </div>
              <div class="body">
                <div class="badges">
                  <span v-if="article.category_name" class="badge category">{{ article.category_name }}</span>
                  <span v-for="tag in article.tags || []" :key="tag.id" class="badge tag">{{ tag.name }}</span>
                </div>
                <div class="title">{{ article.title }}</div>
                <div class="summary">{{ article.summary || '暂无摘要' }}</div>
                <div class="meta">
                  <div class="meta-left">
                    <span class="meta-item">
                      <el-icon><User /></el-icon>
                      {{ article.username || '匿名' }}
                    </span>
                    <span class="meta-item">
                      <el-icon><View /></el-icon>
                      {{ article.view_count || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon><Star /></el-icon>
                      {{ article.like_count || 0 }}
                    </span>
                  </div>
                  <div class="meta-right">{{ formatDate(article.updated_at) }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <div class="pagination-container" v-if="total > pageSize">
          <el-pagination
              background
              layout="prev, pager, next"
              :total="total"
              :page-size="pageSize"
              v-model:current-page="currentPage"
              @current-change="fetchArticles"
          />
        </div>
      </div>

      <!-- 右侧：作者信息 -->
      <div class="sidebar-wrapper">
        <el-card class="author-card" shadow="never">
          <div class="author-header">
            <el-avatar :size="80" :src="authorInfo.avatar || defaultAvatar" />
            <h3 class="author-name">{{ authorInfo.username || '博主' }}</h3>
          </div>

          <div class="author-bio" v-if="authorInfo.introduction">
            {{ authorInfo.introduction }}
          </div>

          <div class="author-stats">
            <div class="stat-item">
              <div class="stat-value">{{ authorInfo.articleCount || 0 }}</div>
              <div class="stat-label">文章</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ categories.length }}</div>
              <div class="stat-label">分类</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ tags.length }}</div>
              <div class="stat-label">标签</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ authorInfo.viewCount || 0 }}</div>
              <div class="stat-label">浏览</div>
            </div>
          </div>

          <el-divider />

          <div class="donate-section" v-if="authorInfo.qr">
            <el-button type="primary" plain @click="showDonateDialog = true" class="donate-btn">
              <el-icon><Present /></el-icon>
              打赏支持
            </el-button>
          </div>
        </el-card>

        <!-- 标签云 -->
        <el-card class="tag-cloud-card mt-6" shadow="never">
          <template #header>
            <div class="sidebar-title" style="padding-left: 0; margin-bottom: 0;">标签云</div>
          </template>
          <div class="tag-cloud">
            <el-tag
                v-for="tag in tags"
                :key="tag.id"
                class="tag-item"
                :type="getTagType(tag.id)"
                :effect="activeTag === tag.id ? 'dark' : 'plain'"
                round
                @click="changeTag(tag.id)"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 打赏弹窗 -->
    <el-dialog v-model="showDonateDialog" title="赞赏支持" width="400px" center>
      <div class="donate-dialog-content">
        <p class="donate-text">如果我的文章对你有帮助，欢迎打赏支持~</p>
        <img :src="qrUrl" alt="收款二维码" class="donate-qr" />
      </div>
    </el-dialog>

    <!--    页脚-->
    <footer class="page-footer">
      <div class="footer-content">
        <div class="footer-divider">
          <span class="divider-line"></span>
          <span class="divider-text">已经阅读到底了</span>
          <span class="divider-line"></span>
        </div>
        <p class="footer-copyright">
          © {{ currentYear }} {{ authorInfo.username || '博主' }} · 记录学习与成长
        </p>
        <p class="footer-subtitle">
          感谢你的阅读与陪伴 ❤️
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const showDonateDialog = ref(false)
const categories = ref([])
const tags = ref([])
const activeCategory = ref('')
const activeTag = ref(0)
const sortType = ref(0) // 0: 最近发布, 1: 热门推荐
const defaultCover = 'https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=300'
const defaultAvatar = '/uploads/avatars/default.png'

// 作者信息
const authorInfo = ref({
  username: '',
  avatar: '',
  introduction: '',
  qr: '',
  articleCount: 0,
  viewCount: 0
})

// 二维码完整URL
const qrUrl = computed(() => {
  return authorInfo.value.qr || ''
})

// 获取作者信息
const fetchAuthorInfo = async () => {
  try {
    const res = await request.get('/user/author')
    if (res.code === 200) {
      authorInfo.value = {
        username: res.data.username || '博主',
        avatar: res.data.avatar || '',
        introduction: res.data.introduction || '这个人很懒，什么都没写~',
        qr: res.data.qr || '',
        articleCount: res.data.article_count || 0,
        viewCount: res.data.view_count || 0
      }
    }
  } catch (error) {
    console.error('获取作者信息失败:', error)
  }
}

// 添加当前年份
const currentYear = ref(new Date().getFullYear())

// 获取分类列表
const fetchCategories = async () => {
  try {
    const res = await request.get('/category/list')
    if (res.code === 200) {
      categories.value = res.data || []
    }
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const res = await request.get('/tag/list')
    if (res.code === 200) {
      tags.value = res.data || []
    }
  } catch (error) {
    console.error('获取标签失败:', error)
  }
}

// 切换分类
const changeCategory = (categoryId) => {
  router.push({ path: '/', query: { ...route.query, category: categoryId || undefined, tag: undefined } })
}

// 切换标签
const changeTag = (tagId) => {
  const newTagId = activeTag.value === tagId ? 0 : tagId
  router.push({ path: '/', query: { ...route.query, tag: newTagId || undefined } })
}

// 切换排序方式
const changeSort = (sort) => {
  router.push({ path: '/', query: { ...route.query, sort: sort || undefined } })
}

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await request.get('/article/list', {
      params: {
        status: 'published',
        query: route.query.q || '',
        tag_id: activeTag.value || 0,
        page: currentPage.value,
        pageSize: pageSize.value,
        category_id: activeCategory.value || 0,
        sort: sortType.value
      }
    })
    if (res.code === 200) {
      articles.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Watch for all query changes in URL to prevent duplicate requests
watch(() => route.query, (newQuery) => {
  activeTag.value = Number(newQuery.tag) || 0
  activeCategory.value = Number(newQuery.category) || ''
  sortType.value = newQuery.sort !== undefined ? Number(newQuery.sort) : 0
  currentPage.value = 1
  fetchArticles()
}, { deep: true })


const goToDetail = (id) => {
  const query = {}
  if (sortType.value !== 0) query.sort = sortType.value
  router.push({ path: `/article/${id}`, query })
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
}

const getTagType = (tagId) => {
  if (activeTag.value === tagId) return ''
  const types = ['', 'success', 'info', 'warning', 'danger']
  return types[tagId % types.length]
}

onMounted(() => {
  const queryTag = Number(route.query.tag)
  if (queryTag) activeTag.value = queryTag
  const queryCat = Number(route.query.category)
  if (queryCat) activeCategory.value = queryCat
  if (route.query.sort !== undefined) sortType.value = Number(route.query.sort)

  fetchArticles()
  fetchAuthorInfo()
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.home {
  max-width: 1300px;
  padding: 20px 20px 20px 16px;
  margin-left: -8px;
}

/* 新的布局结构 */
.content-wrapper {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

/* 左侧侧边栏 */
.left-sidebar {
  width: 160px;
  margin-left: -120px;
}

.sidebar-sticky {
  position: sticky;
  top: 84px;
}

.sidebar-title {
  font-size: 14px;
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 12px;
  padding-left: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.mt-6 {
  margin-top: 24px;
}

.sort-section-vertical, .category-list-vertical {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sort-item-v, .category-item-v {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  transition: all 0.2s;
}

.sort-item-v:hover, .category-item-v:hover {
  background-color: #f1f5f9;
  color: var(--primary-color);
}

.sort-item-v.active, .category-item-v.active {
  background-color: #ecf5ff;
  color: var(--primary-color);
  font-weight: 600;
}

.sort-item-v .el-icon {
  font-size: 18px;
}

/* 分类部分样式 */
.category-section {
  padding: 4px 0;
}

.section-label {
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.category-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-pill {
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: #64748b;
  background-color: #f8fafc;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.category-pill:hover {
  background-color: #f1f5f9;
  color: var(--primary-color);
}

.category-pill.active {
  background-color: #fff;
  color: var(--primary-color);
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
}

/* 中间文章列表 */
.main-wrapper {
  flex: 1;
  min-width: 0;
}

/* 右侧作者卡片 */
.sidebar-wrapper {
  width: 280px;
  flex-shrink: 0;
  position: sticky;
  top: 84px;
  margin-right: -60px;
}

.author-card {
  border-radius: 12px;
}

.author-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0 16px;
}

.author-name {
  margin: 12px 0 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.author-bio {
  padding: 0 16px 16px;
  font-size: 14px;
  color: var(--secondary-text-color);
  line-height: 1.6;
  text-align: center;
}

.author-stats {
  display: flex;
  justify-content: space-around;
  padding: 16px;
  background-color: #f8fafc;
  border-radius: 8px;
  margin: 0 16px 8px;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
}

.stat-label {
  font-size: 12px;
  color: var(--secondary-text-color);
  margin-top: 4px;
}

.donate-section {
  padding: 8px 16px 16px;
  text-align: center;
}

.donate-btn {
  width: 100%;
  border-radius: 20px;
}

.tag-cloud-card {
  border-radius: 12px;
}

.tag-cloud-card :deep(.el-card__header) {
  padding: 12px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 4px 0;
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s;
}

.tag-item:hover {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.donate-dialog-content {
  text-align: center;
  padding: 20px 0;
}

.donate-text {
  margin-bottom: 20px;
  color: var(--secondary-text-color);
}

.donate-qr {
  width: 200px;
  height: 200px;
  object-fit: contain;
  border-radius: 8px;
}

.list-card {
  border-radius: 12px;
}

.state-text {
  text-align: center;
  padding: 60px 0;
  color: var(--secondary-text-color);
  font-size: 15px;
}

.article-item {
  display: flex;
  gap: 20px;
  padding: 24px 16px;
  border-bottom: 1px solid #f1f5f9;
  cursor: pointer;
  transition: all 0.25s ease;
}

.article-item:last-child {
  border-bottom: none;
}

.article-item:hover {
  background-color: #f8fafc;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.cover {
  width: 200px;
  height: 130px;
  border-radius: 10px;
  overflow: hidden;
  background: #f1f5f9;
  flex-shrink: 0;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.article-item:hover .cover img {
  transform: scale(1.05);
}

.body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.badge {
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.badge.category {
  background-color: #ecf5ff;
  color: #409eff;
}

.badge.tag {
  background-color: #f0f9eb;
  color: #67c23a;
}

.title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 8px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.summary {
  font-size: 14px;
  color: var(--secondary-text-color);
  margin-bottom: 12px;
  line-height: 1.6;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #94a3b8;
}

.meta-left {
  display: flex;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-container {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}

/* 响应式布局 */
@media (max-width: 1024px) {
  .left-sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    flex-direction: column;
  }

  .sidebar-wrapper {
    width: 100%;
    position: static;
  }

  .article-item {
    flex-direction: column;
  }

  .cover {
    width: 100%;
    height: 200px;
  }
}
/* 页脚样式 */
.page-footer {
  margin-top: 60px;
  padding: 40px 20px 30px;
  text-align: center;
}

.footer-content {
  max-width: 600px;
  margin: 0 auto;
}

.footer-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 20px;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(
      to right,
      transparent,
      #e2e8f0 20%,
      #e2e8f0 80%,
      transparent
  );
}

.divider-text {
  font-size: 14px;
  color: #94a3b8;
  white-space: nowrap;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.footer-copyright {
  font-size: 13px;
  color: #94a3b8;
  margin: 0 0 8px;
  letter-spacing: 0.3px;
}

.footer-subtitle {
  font-size: 12px;
  color: #cbd5e1;
  margin: 0;
  font-style: italic;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .divider-line {
    background: linear-gradient(
        to right,
        transparent,
        #475569 20%,
        #475569 80%,
        transparent
    );
  }

  .divider-text {
    color: #64748b;
  }

  .footer-copyright {
    color: #64748b;
  }

  .footer-subtitle {
    color: #475569;
  }
}

/* 响应式 */
@media (max-width: 768px) {
  .page-footer {
    margin-top: 40px;
    padding: 30px 16px 24px;
  }

  .footer-divider {
    gap: 12px;
  }

  .divider-text {
    font-size: 13px;
  }
}
</style>