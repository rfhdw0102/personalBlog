<template>
  <div class="article-detail">
    <el-card v-if="article" shadow="never" class="article-card">
      <div class="top">
        <div class="badges">
          <span v-if="article.category_name" class="badge category">{{ article.category_name }}</span>
          <span v-for="tag in article.tags || []" :key="tag.id" class="badge tag">{{ tag.name }}</span>
        </div>
        <h1 class="title">{{ article.title }}</h1>
        <div class="meta">
          <div class="author">
            <el-avatar :size="34" :src="avatarUrl(article.avatar)" />
            <span class="name">{{ article.username || '匿名' }}</span>
          </div>
          <span class="meta-item">{{ formatDate(article.updated_at) }}</span>
          <span class="meta-item"><el-icon><View /></el-icon> {{ article.view_count || 0 }}</span>
          <span class="meta-item"><el-icon><Star /></el-icon> {{ article.like_count || 0 }}</span>
        </div>
      </div>


      <el-divider />

      <div class="content ql-snow">
        <div class="ql-editor" v-html="article.content"></div>
      </div>

      <el-divider />

      <div class="actions">
        <el-button type="primary" :plain="!isLiked" @click="toggleLike">
          <el-icon><Star /></el-icon>
          {{ isLiked ? '已点赞' : '点赞' }}
        </el-button>
        <el-button v-if="isAuthor" @click="goToEdit">编辑</el-button>
      </div>
    </el-card>

    <div class="adjacent-nav">
      <div v-if="prevArticle" class="adjacent-item prev" @click="goToArticle(prevArticle.id)">
        <span class="adjacent-label">上一篇</span>
        <span class="adjacent-title">{{ prevArticle.title }}</span>
      </div>
      <div v-else class="adjacent-item prev disabled">
        <span class="adjacent-label">上一篇</span>
        <span class="adjacent-title">暂无内容</span>
      </div>
      <div v-if="nextArticle" class="adjacent-item next" @click="goToArticle(nextArticle.id)">
        <span class="adjacent-label">下一篇</span>
        <span class="adjacent-title">{{ nextArticle.title }}</span>
      </div>
      <div v-else class="adjacent-item next disabled">
        <span class="adjacent-label">下一篇</span>
        <span class="adjacent-title">暂无内容</span>
      </div>
    </div>

    <el-card v-if="article" shadow="never" class="comment-card">
      <template #header>
        <div class="comment-header">
          <div class="header-left">
            <span class="comment-title">评论</span>
            <el-switch
              v-if="isAuthor"
              v-model="showHiddenComments"
              class="ml-4"
              inline-prompt
              active-text="隐藏"
              inactive-text="正常"
              @change="onCommentTypeChange"
            />
          </div>
          <el-button text type="primary" @click="refreshComments">刷新</el-button>
        </div>
      </template>

      <div v-if="commentsLoading" class="state-text">正在加载评论...</div>
      <div v-else-if="comments.length === 0" class="state-text">暂无评论</div>
      <div v-else class="comment-list">
        <comment-item
          v-for="c in comments"
          :key="c.id"
          :comment="c"
          :is-author="isAuthor"
          :show-hidden-comments="showHiddenComments"
          :article-id="article.id"
          :user-id="user?.id"
          @reply-success="refreshComments"
          @hide-success="refreshComments"
          @unhide-success="refreshComments"
          @recall-success="refreshComments"
        />
      </div>

      <el-divider />

      <div class="comment-form">
        <el-input
          v-model="commentContent"
          type="textarea"
          :rows="3"
          placeholder="写下你的评论..."
          maxlength="500"
          show-word-limit
        />
        <div class="comment-actions">
          <el-button type="primary" :disabled="!commentContent.trim()" :loading="commentSubmitting" @click="submitComment">
            发表评论
          </el-button>
        </div>
      </div>
    </el-card>

    <div v-if="loading" class="loading">正在加载文章...</div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import CommentItem from '@/components/CommentItem.vue'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const isLiked = ref(false)
const loading = ref(false)
const comments = ref([])
const commentsLoading = ref(false)
const commentContent = ref('')
const commentSubmitting = ref(false)
const showHiddenComments = ref(false)
const prevArticle = ref(null)
const nextArticle = ref(null)

const userStr = localStorage.getItem('user')
const user = userStr ? JSON.parse(userStr) : null

const isAuthor = computed(() => {
  return user && article.value && user.id === article.value.user_id
})

const avatarUrl = (path) => {
  if (!path) return '/uploads/avatars/default.png'
  if (path.startsWith('http')) return path
  return path
}

const fetchArticle = async (id) => {
  loading.value = true
  try {
    const res = await request.get(`/article/${id}`)
    if (res.code === 200) {
      article.value = res.data
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const incrementView = async (id) => {
  try {
    await request.post(`/article/${id}/view`)
  } catch (e) {
    console.error(e)
  }
}

const goToEdit = () => {
  router.push(`/article/edit/${article.value.id}`)
}

const fetchLikeStatus = async (id) => {
  try {
    const res = await request.get(`/article/${id}/like-status`, {
      params: { uid: user?.id || 0 }
    })
    if (res.code === 200) {
      isLiked.value = !!res.data.is_liked
    }
  } catch (e) {
    console.error(e)
  }
}

const toggleLike = async () => {
  if (!localStorage.getItem('token')) {
    ElMessage.warning('请先登录')
    return
  }
  const id = route.params.id
  try {
    if (!isLiked.value) {
      const res = await request.post(`/article/like/${id}`)
      if (res.code === 200) {
        isLiked.value = true
        if (article.value) article.value.like_count = (article.value.like_count || 0) + 1
        ElMessage.success('点赞成功')
      }
    } else {
      const res = await request.post(`/article/unlike/${id}`)
      if (res.code === 200) {
        isLiked.value = false
        if (article.value) article.value.like_count = Math.max((article.value.like_count || 0) - 1, 0)
        ElMessage.success('取消点赞成功')
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const refreshComments = async () => {
  const id = route.params.id
  if (!id) return
  commentsLoading.value = true
  try {
    const url = showHiddenComments.value 
      ? `/comment/article/${id}/hide`
      : `/comment/article/${id}`
    const res = await request.get(url)
    if (res.code === 200) {
      comments.value = res.data || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    commentsLoading.value = false
  }
}

const onCommentTypeChange = () => {
  refreshComments()
}

const submitComment = async () => {
  if (!localStorage.getItem('token')) {
    ElMessage.warning('请先登录')
    return
  }
  const id = Number(route.params.id)
  if (!id) return
  commentSubmitting.value = true
  try {
    const res = await request.post('/comment', { article_id: id, content: commentContent.value })
    if (res.code === 200) {
      ElMessage.success('评论成功')
      commentContent.value = ''
      refreshComments()
    }
  } catch (e) {
    console.error(e)
  } finally {
    commentSubmitting.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}

const fetchAdjacent = async (id) => {
  try {
    const res = await request.get(`/article/${id}/adjacent`, {
      params: { sort: Number(route.query.sort) || 0 }
    })
    if (res.code === 200) {
      prevArticle.value = res.data.prev
      nextArticle.value = res.data.next
    }
  } catch (e) {
    console.error(e)
  }
}

const goToArticle = (id) => {
  router.push(`/article/${id}`)
}

const loadArticle = (id) => {
  incrementView(id)
  fetchArticle(id)
  fetchLikeStatus(id)
  fetchAdjacent(id)
  refreshComments()
}

onMounted(() => {
  const id = route.params.id
  if (!id) return
  loadArticle(id)
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    window.scrollTo(0, 0)
    loadArticle(newId)
  }
})
</script>

<style scoped>
.article-detail {
  max-width: 900px;
  margin: 0 auto;
}
.article-card {
  margin-bottom: 24px;
  padding: 20px;
}
.top {
  margin-bottom: 32px;
}
.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}
.badge {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 13px;
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
  font-size: 32px;
  font-weight: 800;
  color: var(--text-color);
  margin-bottom: 20px;
  line-height: 1.3;
}
.meta {
  display: flex;
  align-items: center;
  gap: 24px;
  color: #94a3b8;
  font-size: 14px;
}
.author {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-color);
  font-weight: 600;
}
.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}
.cover {
  width: 100%;
  max-height: 400px;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 32px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}
.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.content {
  font-size: 17px;
  line-height: 1.8;
  color: var(--text-color);
}
.actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 40px;
}

.comment-card {
  margin-bottom: 40px;
}
.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.ml-4 {
  margin-left: 1rem;
}
.comment-title {
  font-size: 18px;
  font-weight: 700;
}
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.comment-form {
  margin-top: 24px;
}
.comment-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
.state-text {
  text-align: center;
  padding: 30px 0;
  color: var(--secondary-text-color);
}
.loading {
  text-align: center;
  padding: 100px 0;
  color: var(--secondary-text-color);
}
.adjacent-nav {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
}
.adjacent-item {
  flex: 1;
  padding: 16px 20px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.2s;
}
.adjacent-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.12);
}
.adjacent-item.disabled {
  cursor: default;
  color: #94a3b8;
}
.adjacent-item.disabled:hover {
  border-color: #e5e7eb;
  box-shadow: none;
}
.adjacent-item.next {
  text-align: right;
}
.adjacent-label {
  display: block;
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 6px;
}
.adjacent-title {
  font-size: 14px;
  color: var(--text-color);
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
