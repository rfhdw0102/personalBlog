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
import { ref, onMounted, computed } from 'vue'
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
    const res = await request.get(`/article/${id}/like-status`)
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

onMounted(() => {
  const id = route.params.id
  if (!id) return
  incrementView(id)
  fetchArticle(id)
  fetchLikeStatus(id)
  refreshComments()
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
</style>
