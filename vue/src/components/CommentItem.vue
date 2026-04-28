<template>
  <div class="comment-item">
    <div class="comment-meta">
      <div class="comment-user">
        <el-avatar :size="28" :src="avatarUrl(comment.avatar)" />
        <div class="user-info">
          <span class="comment-username">{{ comment.username || '匿名用户' }}</span>
          <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
        </div>
      </div>
      <div class="comment-actions-wrapper">
        <el-button v-if="!showHiddenComments" text type="primary" size="small" @click="showReplyInput = !showReplyInput">
          {{ showReplyInput ? '取消回复' : '回复' }}
        </el-button>
        <el-button v-if="isOwner && !showHiddenComments" text type="warning" size="small" @click="handleRecall">撤回</el-button>
        <el-button v-if="isAuthor && !showHiddenComments" text type="danger" size="small" @click="handleHide">隐藏</el-button>
        <el-button v-if="isAuthor && showHiddenComments" text type="primary" size="small" @click="handleUnhide">取消隐藏</el-button>
      </div>
    </div>
    <div class="comment-content">{{ comment.content }}</div>

    <div v-if="showReplyInput" class="reply-form">
      <el-input
        v-model="replyContent"
        type="textarea"
        :rows="2"
        :placeholder="`回复 @${comment.username}...`"
        maxlength="500"
        show-word-limit
      />
      <div class="reply-actions">
        <el-button size="small" @click="showReplyInput = false">取消</el-button>
        <el-button type="primary" size="small" :disabled="!replyContent.trim()" :loading="replySubmitting" @click="submitReply">
          发表回复
        </el-button>
      </div>
    </div>

    <div v-if="comment.children && comment.children.length > 0" class="comment-children">
      <comment-item
        v-for="child in comment.children"
        :key="child.id"
        :comment="child"
        :is-author="isAuthor"
        :show-hidden-comments="showHiddenComments"
        :article-id="articleId"
        :user-id="userId"
        @reply-success="$emit('reply-success')"
        @hide-success="(id) => $emit('hide-success', id)"
        @unhide-success="(id) => $emit('unhide-success', id)"
        @recall-success="(id) => $emit('recall-success', id)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  isAuthor: {
    type: Boolean,
    default: false
  },
  showHiddenComments: {
    type: Boolean,
    default: false
  },
  articleId: {
    type: Number,
    required: true
  },
  userId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['reply-success', 'hide-success', 'unhide-success', 'recall-success'])

const showReplyInput = ref(false)
const replyContent = ref('')
const replySubmitting = ref(false)

const isOwner = computed(() => {
  return props.userId && props.comment.user_id === props.userId
})

const avatarUrl = (path) => {
  if (!path) return '/uploads/avatars/default.png'
  if (path.startsWith('http')) return path
  return path
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}

const submitReply = async () => {
  if (!localStorage.getItem('token')) {
    ElMessage.warning('请先登录')
    return
  }
  
  replySubmitting.value = true
  try {
    const res = await request.post('/comment', {
      article_id: props.articleId,
      content: replyContent.value,
      parent_id: props.comment.id
    })
    if (res.code === 200) {
      ElMessage.success('回复成功')
      replyContent.value = ''
      showReplyInput.value = false
      emit('reply-success')
    }
  } catch (e) {
    console.error(e)
  } finally {
    replySubmitting.value = false
  }
}

const handleHide = async () => {
  try {
    const res = await request.put(`/comment/${props.comment.id}/hide`)
    if (res.code === 200) {
      ElMessage.success('已隐藏')
      emit('hide-success', props.comment.id)
    }
  } catch (e) {
    console.error(e)
  }
}

const handleUnhide = async () => {
  try {
    const res = await request.put(`/comment/${props.comment.id}/unhide`)
    if (res.code === 200) {
      ElMessage.success('已取消隐藏')
      emit('unhide-success', props.comment.id)
    }
  } catch (e) {
    console.error(e)
  }
}

const handleRecall = async () => {
  try {
    const res = await request.delete(`/comment/${props.comment.id}`)
    if (res.code === 200) {
      ElMessage.success('评论已撤回')
      emit('recall-success', props.comment.id)
    }
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.comment-item {
  margin-top: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f1f5f9;
}
.comment-item:last-child {
  border-bottom: none;
}
.comment-meta {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}
.comment-user {
  display: flex;
  align-items: center;
  gap: 10px;
}
.user-info {
  display: flex;
  flex-direction: column;
}
.comment-username {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-color);
}
.comment-time {
  font-size: 12px;
  color: #94a3b8;
}
.comment-content {
  font-size: 15px;
  line-height: 1.6;
  color: var(--secondary-text-color);
  padding-left: 38px;
  word-break: break-all;
}
.comment-actions-wrapper {
  display: flex;
  gap: 4px;
}
.reply-form {
  margin-top: 12px;
  padding-left: 38px;
}
.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}
.comment-children {
  margin-top: 8px;
  padding-left: 38px;
  border-left: 2px solid #f1f5f9;
}
</style>
