<template>
  <div class="management">
    <el-row :gutter="20">
      <el-col :span="5" class="sidebar-col">
        <el-card shadow="never" class="sidebar-card">
          <el-menu :default-active="activeMenu" @select="onSelect" class="sidebar-menu">
            <el-menu-item index="messages">
              <el-icon><Bell /></el-icon>
              <span>消息中心</span>
              <el-badge v-if="unreadCount > 0" :value="unreadCount" class="sidebar-badge" />
            </el-menu-item>
            <el-menu-item index="comments">
              <el-icon><ChatDotRound /></el-icon>
              <span>评论管理</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>

      <el-col :span="19">
        <el-card v-if="activeMenu === 'messages'" shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <div class="panel-title">消息通知</div>
              <el-button @click="markAllRead" :disabled="unreadCount === 0">全部标记已读</el-button>
            </div>
          </template>

          <div v-if="notificationLoading" class="state-text">正在加载通知...</div>
          <div v-else-if="notifications.length === 0" class="state-text">暂无通知</div>
          <div v-else class="notification-list">
            <div v-for="n in notifications" :key="n.id" class="notification-item" :class="{ unread: !n.is_read }">
              <div class="n-left">
                <el-avatar :size="34" :src="n.actor_avatar || '/uploads/avatars/default.png'" />
                <div class="n-body">
                  <div class="n-title">
                    <span class="actor">{{ n.actor_name }}</span>
                    <template v-if="n.type === 'reply'">
                      <span class="type">回复了你在文章</span>
                      <span class="article" @click="router.push(`/article/${n.article_id}`)">{{ n.article_title }}</span>
                      <span class="type">中的评论</span>
                      <span class="parent-content">（{{ n.parent_comment_content }}）</span>
                    </template>
                    <template v-else>
                      <span class="type">{{ typeText(n.type) }}</span>
                      <span class="article" @click="router.push(`/article/${n.article_id}`)">{{ n.article_title }}</span>
                    </template>
                  </div>
                  <div class="n-content">{{ n.comment_content }}</div>
                  <div class="n-time">{{ formatDate(n.created_at) }}</div>
                </div>
              </div>
              <div class="n-actions">
                <el-button v-if="!n.is_read" text type="primary" @click="markRead(n.id)">标记已读</el-button>
                <el-button text type="danger" @click="deleteNotification(n.id)">删除</el-button>
              </div>
            </div>
          </div>

          <div class="pagination" v-if="notificationTotal > pageSize">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="notificationTotal"
              :page-size="pageSize"
              v-model:current-page="notificationPage"
              @current-change="fetchNotifications"
            />
          </div>
        </el-card>

        <el-card v-if="activeMenu === 'comments'" shadow="never" class="panel-card comment-panel">
          <template #header>
            <div class="panel-header">
              <el-tabs v-model="activeCommentTab" @tab-change="handleCommentTabChange">
                <el-tab-pane label="收到评论" name="received" />
                <el-tab-pane label="我的评论" name="sent" />
                <el-tab-pane v-if="isAdmin" label="隐藏评论" name="hidden" />
              </el-tabs>
            </div>
          </template>

          <div v-if="commentLoading" class="state-text">正在加载评论...</div>
          <div v-else-if="activeCommentTab === 'received'">
            <div v-if="receivedComments.length === 0" class="state-text">暂无收到评论</div>
            <div v-else class="comment-list">
              <div v-for="c in receivedComments" :key="c.id" class="u-comment">
                <div class="u-top">
                  <div class="u-article" @click="router.push(`/article/${c.article_id}`)">{{ c.article_title }}</div>
                  <div class="u-time">{{ formatDate(c.created_at) }}</div>
                </div>
                <div class="u-body">
                  <el-avatar :size="28" :src="c.avatar || '/uploads/avatars/default.png'" />
                  <div class="u-right">
                    <div class="u-user">
                      {{ c.username }}
                      <template v-if="c.parent_id > 0">
                        <span class="u-reply-text">回复了</span>
                        {{ c.parent_username || '用户' }}
                      </template>
                    </div>
                    <div class="u-content">{{ c.content }}</div>
                    <div v-if="replyingCommentId === c.id" class="u-reply-form">
                      <el-input
                        v-model="replyContent"
                        type="textarea"
                        :rows="2"
                        placeholder="输入回复内容..."
                        maxlength="500"
                        show-word-limit
                      />
                      <div class="u-reply-actions">
                        <el-button size="small" @click="replyingCommentId = null">取消</el-button>
                        <el-button type="primary" size="small" :loading="replySubmitting" :disabled="!replyContent.trim()" @click="submitReply(c)">提交回复</el-button>
                      </div>
                    </div>
                  </div>
                  <div class="u-actions">
                    <el-button text type="primary" @click="toggleReply(c)">{{ replyingCommentId === c.id ? '取消回复' : '回复' }}</el-button>
                    <el-button text type="danger" @click="hideReceivedComment(c.id)">隐藏</el-button>
                  </div>
                </div>
              </div>
            </div>
            <div class="pagination" v-if="receivedTotal > pageSize">
              <el-pagination
                background
                layout="prev, pager, next"
                :total="receivedTotal"
                :page-size="pageSize"
                v-model:current-page="receivedPage"
                @current-change="fetchReceivedComments"
              />
            </div>
          </div>

          <div v-else-if="activeCommentTab === 'sent'">
            <div v-if="sentComments.length === 0" class="state-text">暂无我的评论</div>
            <div v-else class="comment-list">
              <div v-for="c in sentComments" :key="c.id" class="u-comment">
                <div class="u-top">
                  <div class="u-article" @click="router.push(`/article/${c.article_id}`)">{{ c.article_title }}</div>
                  <div class="u-time">{{ formatDate(c.created_at) }}</div>
                </div>
                <div class="u-body">
                  <el-avatar :size="28" :src="c.avatar || '/uploads/avatars/default.png'" />
                  <div class="u-right">
                    <div class="u-user">
                      {{ c.username }}
                      <template v-if="c.parent_id > 0">
                        <span class="u-reply-text">回复了</span>
                        {{ c.parent_username || '用户' }}
                      </template>
                    </div>
                    <div class="u-content">{{ c.content }}</div>
                  </div>
                  <el-popconfirm title="确定撤回该评论吗？" @confirm="recallSentComment(c.id)">
                    <template #reference>
                      <el-button text type="danger">撤回</el-button>
                    </template>
                  </el-popconfirm>
                </div>
              </div>
            </div>
            <div class="pagination" v-if="sentTotal > pageSize">
              <el-pagination
                background
                layout="prev, pager, next"
                :total="sentTotal"
                :page-size="pageSize"
                v-model:current-page="sentPage"
                @current-change="fetchSentComments"
              />
            </div>
          </div>

          <div v-else-if="activeCommentTab === 'hidden'">
            <div v-if="hiddenComments.length === 0" class="state-text">暂无隐藏评论</div>
            <div v-else class="comment-list">
              <div v-for="c in hiddenComments" :key="c.id" class="u-comment">
                <div class="u-top">
                  <div class="u-article" @click="router.push(`/article/${c.article_id}`)">{{ c.article_title }}</div>
                  <div class="u-time">{{ formatDate(c.created_at) }}</div>
                </div>
                <div class="u-body">
                  <el-avatar :size="28" :src="c.avatar || '/uploads/avatars/default.png'" />
                  <div class="u-right">
                    <div class="u-user">
                      {{ c.username }}
                      <template v-if="c.parent_id > 0">
                        <span class="u-reply-text">回复了</span>
                        {{ c.parent_username || '用户' }}
                      </template>
                    </div>
                    <div class="u-content">{{ c.content }}</div>
                  </div>
                  <el-button text type="primary" @click="unhideComment(c.id)">取消隐藏</el-button>
                </div>
              </div>
            </div>
            <div class="pagination" v-if="hiddenTotal > pageSize">
              <el-pagination
                background
                layout="prev, pager, next"
                :total="hiddenTotal"
                :page-size="pageSize"
                v-model:current-page="hiddenPage"
                @current-change="fetchHiddenComments"
              />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import {ref, onMounted, watch, computed} from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

const pageSize = 10

const activeMenu = ref('messages')
const unreadCount = ref(0)

const notifications = ref([])
const notificationTotal = ref(0)
const notificationPage = ref(1)
const notificationLoading = ref(false)

const activeCommentTab = ref('received')
const receivedComments = ref([])
const receivedTotal = ref(0)
const receivedPage = ref(1)
const sentComments = ref([])
const sentTotal = ref(0)
const sentPage = ref(1)
const hiddenComments = ref([])
const hiddenTotal = ref(0)
const hiddenPage = ref(1)
const commentLoading = ref(false)

const replyingCommentId = ref(null)
const replyContent = ref('')
const replySubmitting = ref(false)

const toggleReply = (comment) => {
  if (replyingCommentId.value === comment.id) {
    replyingCommentId.value = null
    replyContent.value = ''
  } else {
    replyingCommentId.value = comment.id
    replyContent.value = ''
  }
}

const submitReply = async (comment) => {
  if (!replyContent.value.trim()) return
  replySubmitting.value = true
  try {
    const res = await request.post('/comment', {
      article_id: comment.article_id,
      content: replyContent.value,
      parent_id: comment.id
    })
    if (res.code === 200) {
      ElMessage.success('回复成功')
      replyingCommentId.value = null
      replyContent.value = ''
      // 切换到“我的评论”标签查看回复，或者刷新当前列表
      fetchReceivedComments()
    }
  } catch (e) {
    console.error(e)
  } finally {
    replySubmitting.value = false
  }
}

const isAdmin = computed(() => {
  const userInfo = localStorage.getItem('user')
  if (userInfo) {
    try {
      const user = JSON.parse(userInfo)
      return user.role === 'admin'
    } catch (e) {
      return false
    }
  }
  return false
})
const onSelect = (key) => {
  activeMenu.value = key
  router.replace({ path: '/management', query: { tab: key } })
}

const syncTabFromRoute = () => {
  const tab = route.query.tab
  if (tab === 'messages' || tab === 'comments') {
    activeMenu.value = tab
  } else {
    activeMenu.value = 'messages'
  }
}

watch(
  () => route.query.tab,
  () => {
    syncTabFromRoute()
    if (activeMenu.value === 'messages') fetchNotifications()
    if (activeMenu.value === 'comments') {
      if (activeCommentTab.value === 'received') fetchReceivedComments()
      else if (activeCommentTab.value === 'sent') fetchSentComments()
      else fetchHiddenComments()
    }
  }
)

const fetchUnread = async () => {
  try {
    const res = await request.get('/notification/unread-count')
    if (res.code === 200) {
      unreadCount.value = res.data.count || 0
      window.dispatchEvent(new CustomEvent('refresh-unread', { detail: unreadCount.value }))
    }
  } catch (e) {
    unreadCount.value = 0
  }
}

const fetchNotifications = async () => {
  notificationLoading.value = true
  try {
    const res = await request.get('/notification/list', {
      params: { page: notificationPage.value, pageSize }
    })
    if (res.code === 200) {
      notifications.value = res.data.list
      notificationTotal.value = res.data.total
      fetchUnread()
    }
  } finally {
    notificationLoading.value = false
  }
}

const markRead = async (id) => {
  const res = await request.put(`/notification/${id}/read`)
  if (res.code === 200) {
    const item = notifications.value.find(n => n.id === id)
    if (item && !item.is_read) {
      item.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
      window.dispatchEvent(new CustomEvent('refresh-unread', { detail: unreadCount.value }))
    }
  }
}

const markAllRead = async () => {
  const res = await request.put('/notification/read')
  if (res.code === 200) {
    ElMessage.success('已全部标记已读')
    notifications.value.forEach(n => n.is_read = true)
    unreadCount.value = 0
    window.dispatchEvent(new CustomEvent('refresh-unread', { detail: 0 }))
  }
}

const deleteNotification = async (id) => {
  const res = await request.delete(`/notification/${id}`)
  if (res.code === 200) {
    ElMessage.success('已删除')
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      const wasUnread = !notifications.value[index].is_read
      notifications.value.splice(index, 1)
      notificationTotal.value--
      if (wasUnread) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
        window.dispatchEvent(new CustomEvent('refresh-unread', { detail: unreadCount.value }))
      }
    }
  }
}

const handleCommentTabChange = (tab) => {
  if (tab === 'received') fetchReceivedComments()
  else if (tab === 'sent') fetchSentComments()
  else if (tab === 'hidden' && isAdmin.value) fetchHiddenComments()
}

const fetchReceivedComments = async () => {
  commentLoading.value = true
  try {
    const res = await request.get('/comment/user', {
      params: { page: receivedPage.value, pageSize }
    })
    if (res.code === 200) {
      receivedComments.value = res.data.list
      receivedTotal.value = res.data.total
    }
  } finally {
    commentLoading.value = false
  }
}

const fetchSentComments = async () => {
  commentLoading.value = true
  try {
    const res = await request.get('/comment/userTake', {
      params: { page: sentPage.value, pageSize }
    })
    if (res.code === 200) {
      sentComments.value = res.data.list
      sentTotal.value = res.data.total
    }
  } finally {
    commentLoading.value = false
  }
}

const fetchHiddenComments = async () => {
  commentLoading.value = true
  try {
    const res = await request.get('/comment/user/hide', {
      params: { page: hiddenPage.value, pageSize }
    })
    if (res.code === 200) {
      hiddenComments.value = res.data.list
      hiddenTotal.value = res.data.total
    }
  } finally {
    commentLoading.value = false
  }
}

const hideReceivedComment = async (id) => {
  const res = await request.put(`/comment/${id}/hide`)
  if (res.code === 200) {
    ElMessage.success('已隐藏')
    const index = receivedComments.value.findIndex(c => c.id === id)
    if (index !== -1) {
      receivedComments.value.splice(index, 1)
      receivedTotal.value--
    }
  }
}

const unhideComment = async (id) => {
  const res = await request.put(`/comment/${id}/unhide`)
  if (res.code === 200) {
    ElMessage.success('已取消隐藏')
    const index = hiddenComments.value.findIndex(c => c.id === id)
    if (index !== -1) {
      hiddenComments.value.splice(index, 1)
      hiddenTotal.value--
    }
  }
}

const recallSentComment = async (id) => {
  const res = await request.delete(`/comment/${id}`)
  if (res.code === 200) {
    ElMessage.success('评论已撤回')
    const index = sentComments.value.findIndex(c => c.id === id)
    if (index !== -1) {
      sentComments.value.splice(index, 1)
      sentTotal.value--
    }
  }
}

const typeText = (type) => {
  if (type === 'comment') return '评论了'
  if (type === 'like') return '点赞了'
  if (type === 'reply') return '回复了'
  return '通知了'
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  syncTabFromRoute()
  fetchUnread()
  if (activeMenu.value === 'messages') fetchNotifications()
  if (activeMenu.value === 'comments') {
    if (activeCommentTab.value === 'received') fetchReceivedComments()
    else if (activeCommentTab.value === 'sent') fetchSentComments()
    else fetchHiddenComments()
  }
})
</script>

<style scoped>
.management {
  padding-top: 12px;
}
.sidebar-col {
  position: sticky;
  top: 24px;
}
.sidebar-card {
  padding: 0;
  overflow: hidden;
}
.sidebar-menu {
  border-right: none;
}
.sidebar-menu :deep(.el-menu-item) {
  height: 54px;
  line-height: 54px;
  margin: 4px 0;
  border-radius: 8px;
  display: flex;
  align-items: center;
}
.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: #ecf5ff;
  font-weight: 600;
}
.sidebar-badge {
  margin-left: 8px;
  display: inline-flex;
  align-items: center;
}
.sidebar-badge :deep(.el-badge__content) {
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.panel-card {
  min-height: 600px;
}
.panel-card :deep(.el-card__header) {
  padding: 0 20px;
  height: 70px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
}
.comment-panel :deep(.el-card__header) {
  display: block;
  height: auto;
  min-height: 70px;
}
.comment-panel :deep(.el-tabs__header) {
  margin: 0;
}
.comment-panel :deep(.el-tabs__nav-wrap::after) {
  display: none;
}
.u-right.no-avatar {
  padding-left: 0;
}
.panel-card :deep(.el-card__body) {
  min-height: 540px;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 100%;
}
.panel-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
  line-height: 1;
  padding: 0;
  display: flex;
  align-items: center;
}
.panel-header .el-button {
  flex-shrink: 0;
  margin-left: auto;
}
.filters {
  display: flex;
  gap: 12px;
  flex: 1;
  max-width: 500px;
}
.status-select {
  width: 120px;
}

.state-text {
  text-align: center;
  padding: 60px 0;
  color: var(--secondary-text-color);
}

.notification-list {
  display: flex;
  flex-direction: column;
}
.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.2s;
}
.notification-item:hover {
  background-color: #f8fafc;
}
.notification-item.unread {
  background-color: #f0f9ff;
}
.n-left {
  display: flex;
  gap: 12px;
  flex: 1;
}
.n-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.n-title {
  font-size: 14px;
}
.actor {
  font-weight: 600;
  color: var(--text-color);
}
.type {
  color: var(--secondary-text-color);
  margin: 0 4px;
}
.article {
  color: var(--primary-color);
  font-weight: 500;
  cursor: pointer;
}
.parent-content {
  color: #94a3b8;
  font-size: 13px;
  font-style: italic;
}
.n-content {
  font-size: 13px;
  color: var(--secondary-text-color);
  background: #f1f5f9;
  padding: 8px 12px;
  border-radius: 6px;
  margin: 4px 0;
}
.n-time {
  font-size: 12px;
  color: #94a3b8;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.u-comment {
  padding: 16px;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
}
.u-top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
.u-article {
  font-weight: 600;
  color: var(--primary-color);
  cursor: pointer;
  font-size: 14px;
}
.u-time {
  font-size: 12px;
  color: #94a3b8;
}
.u-body {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}
.u-right {
  flex: 1;
}
.u-user {
  font-weight: 600;
  font-size: 13px;
  margin-bottom: 4px;
}
.u-reply-text {
  color: #94a3b8;
  font-weight: 400;
  margin: 0 4px;
}
.u-content {
  font-size: 14px;
  color: var(--secondary-text-color);
  line-height: 1.5;
}
.u-reply-form {
  margin-top: 12px;
  background: #f8fafc;
  padding: 12px;
  border-radius: 8px;
}
.u-reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}
.u-actions {
  display: flex;
  flex-direction: row;  /* 改为水平排列 */
  gap: 8px;
  align-items: center;
  margin-left: auto;  /* 推到最右边 */
  flex-shrink: 0;
}
.u-actions .el-button {
  padding: 5px 8px;
}
.u-right {
  flex: 1;
}
.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
</style>