<template>
  <div class="about-container">
    <el-card v-if="about" class="about-card" shadow="never">
      <div class="about-header">
        <h1>关于本站</h1>
        <p class="subtitle">分享技术，记录生活</p>
        <el-button v-if="isAdmin" type="warning" size="small" class="edit-btn" @click="toggleEdit">
          {{ editing ? '取消编辑' : '编辑' }}
        </el-button>
      </div>

      <el-divider />

      <div v-if="!editing" class="about-content">
        <section>
          <h2><el-icon><InfoFilled /></el-icon> 关于我</h2>
          <p>{{ about.content || '暂无内容' }}</p>
        </section>
        <section>
          <h2><el-icon><InfoFilled /></el-icon> 博客简介</h2>
          <p>{{ about.intro || '暂无内容' }}</p>
        </section>

        <section>
          <h2><el-icon><Tools /></el-icon> 我的技术栈</h2>
          <div class="tech-stack">
            <el-tag v-for="(item, i) in techStackList" :key="i" size="large">{{ item }}</el-tag>
            <span v-if="techStackList.length === 0" class="empty-hint">暂无内容</span>
          </div>
        </section>

        <section>
          <h2><el-icon><Message /></el-icon> 联系我</h2>
          <p>如果你对我的文章感兴趣，或者有任何问题与建议，欢迎通过以下方式联系我：</p>
          <div class="contact-list">
            <div v-if="about.email" class="contact-item" @click="copyEmail">
              <el-icon class="contact-icon"><Message /></el-icon>
              <span>QQ邮箱：</span>
              <span class="contact-link copyable">{{ about.email }}</span>
              <el-icon v-if="copied" class="copy-success"><CircleCheck /></el-icon>
            </div>
            <div v-if="about.github" class="contact-item">
              <el-icon class="contact-icon"><Platform /></el-icon>
              <span>GitHub：</span>
              <a :href="githubUrl" target="_blank" class="contact-link">{{ about.github }}</a>
            </div>
          </div>
        </section>
        <p>期待与你在技术道路上共同进步！</p>
      </div>

      <div v-else class="about-edit">
        <el-form label-position="top">
          <el-form-item label="关于我">
            <el-input v-model="form.content" type="textarea" :rows="4" maxlength="2000" show-word-limit />
          </el-form-item>
          <el-form-item label="博客简介">
            <el-input v-model="form.intro" type="textarea" :rows="4" maxlength="2000" show-word-limit />
          </el-form-item>
          <el-form-item label="技术栈（逗号分隔）">
            <el-input v-model="form.tech_stack" placeholder="Vue 3, Go, MySQL" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="form.email" placeholder="example@qq.com" />
          </el-form-item>
          <el-form-item label="GitHub 用户名或地址">
            <el-input v-model="form.github" placeholder="rfhdw0102" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="saving" @click="saveAbout">保存</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <div v-if="loading" class="loading">正在加载...</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { InfoFilled, Tools, Message, Platform, CircleCheck } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const about = ref(null)
const loading = ref(false)
const editing = ref(false)
const saving = ref(false)
const copied = ref(false)

const userStr = localStorage.getItem('user')
const user = userStr ? JSON.parse(userStr) : null
const isAdmin = computed(() => user?.role === 'admin')

const form = ref({
  content: '',
  intro: '',
  tech_stack: '',
  email: '',
  github: ''
})

const techStackList = computed(() => {
  if (!about.value?.tech_stack) return []
  try {
    const arr = JSON.parse(about.value.tech_stack)
    return Array.isArray(arr) ? arr : []
  } catch {
    return about.value.tech_stack.split(',').map(s => s.trim()).filter(Boolean)
  }
})

const githubUrl = computed(() => {
  const g = about.value?.github || ''
  if (!g) return ''
  if (g.startsWith('http')) return g
  return `https://github.com/${g}`
})

const fetchAbout = async () => {
  loading.value = true
  try {
    const res = await request.get('/about')
    if (res.code === 200) {
      about.value = res.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const saveAbout = async () => {
  saving.value = true
  try {
    let techStack = form.value.tech_stack
    if (techStack) {
      const arr = techStack.split(',').map(s => s.trim()).filter(Boolean)
      techStack = JSON.stringify(arr)
    } else {
      techStack = '[]'
    }
    const payload = {
      content: form.value.content,
      intro: form.value.intro,
      tech_stack: techStack,
      email: form.value.email,
      github: form.value.github
    }
    const res = await request.put('/about', payload)
    if (res.code === 200) {
      about.value = res.data
      editing.value = false
      ElMessage.success('保存成功')
    }
  } catch (e) {
    console.error(e)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const copyEmail = async () => {
  try {
    await navigator.clipboard.writeText(about.value.email)
    copied.value = true
    ElMessage.success('邮箱已复制到剪贴板')
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

const toggleEdit = () => {
  editing.value = !editing.value
  if (editing.value) {
    form.value = {
      content: about.value?.content || '',
      intro: about.value?.intro || '',
      tech_stack: techStackList.value.join(', '),
      email: about.value?.email || '',
      github: about.value?.github || ''
    }
  }
}

onMounted(async () => {
  await fetchAbout()
})
</script>

<style scoped>
.about-container {
  max-width: 800px;
  margin: 40px auto;
  padding: 0 20px;
}

.about-card {
  border-radius: 16px;
  padding: 20px;
}

.about-header {
  text-align: center;
  margin-bottom: 30px;
  position: relative;
}

.about-header h1 {
  font-size: 32px;
  font-weight: 800;
  color: var(--text-color);
  margin-bottom: 8px;
}

.subtitle {
  color: var(--secondary-text-color);
  font-size: 16px;
}

.edit-btn {
  position: absolute;
  right: 0;
  top: 0;
}

.about-content section {
  margin-bottom: 32px;
}

.about-content h2 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--primary-color);
}

.about-content p {
  line-height: 1.8;
  color: var(--secondary-text-color);
}

.tech-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.empty-hint {
  color: var(--secondary-text-color);
}

.about-content ul {
  list-style: none;
  padding-left: 0;
}

.about-content li {
  position: relative;
  padding-left: 20px;
  margin-bottom: 10px;
  color: var(--secondary-text-color);
}

.about-content li::before {
  content: "•";
  position: absolute;
  left: 0;
  color: var(--primary-color);
  font-weight: bold;
}

.contact-list {
  margin-top: 16px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #f8fafc;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.3s;
}

.contact-item:first-child {
  cursor: pointer;
}

.contact-item:first-child:hover {
  background: #e6f4ff;
  transform: translateX(4px);
}

.contact-item:last-child:hover {
  background: #f1f5f9;
  transform: translateX(4px);
}

.contact-icon {
  font-size: 20px;
  color: var(--primary-color);
}

.contact-item span {
  color: var(--secondary-text-color);
  font-size: 15px;
}

.contact-link {
  color: var(--primary-color);
  text-decoration: none;
  font-family: monospace;
  font-size: 15px;
  font-weight: 500;
  transition: all 0.3s;
  border-bottom: 1px dashed transparent;
}

.contact-link.copyable {
  cursor: pointer;
  user-select: all;
}

a.contact-link:hover {
  border-bottom-color: var(--primary-color);
}

.copy-success {
  color: #67c23a;
  font-size: 18px;
  margin-left: 8px;
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.about-edit {
  padding: 0 20px;
}

.loading {
  text-align: center;
  padding: 100px 0;
  color: var(--secondary-text-color);
}
</style>
