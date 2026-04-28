<template>
  <div class="article-edit">
    <el-card shadow="never" class="edit-card">
      <template #header>
        <div class="card-header">
          <span class="title">{{ isEdit ? '编辑文章' : '写文章' }}</span>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px" class="form">
        <el-form-item label="文章标题" prop="title">
          <el-input v-model="form.title" placeholder="输入文章标题（最多 100 字）" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item label="文章概述" prop="summary">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="输入文章概述（最多 200 字）" maxlength="200" show-word-limit />
        </el-form-item>

        <el-form-item label="文章分类">
          <el-select v-model="form.category_id" placeholder="请选择分类" class="w-100" clearable>
            <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="文章标签">
          <el-select v-model="form.tag_ids" multiple placeholder="请选择标签" class="w-100" clearable>
            <el-option v-for="item in tags" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
          <div class="custom-tag-line">
            <el-input v-model="customTagInput" placeholder="输入自定义标签，按回车添加" @keyup.enter="addCustomTag" />
            <el-button type="primary" @click="addCustomTag">添加</el-button>
          </div>
          <div class="custom-tags" v-if="customTags.length">
            <el-tag v-for="t in customTags" :key="t" closable @close="removeCustomTag(t)">{{ t }}</el-tag>
          </div>
        </el-form-item>

        <el-form-item label="文章封面图">
          <div class="cover-row">
            <div v-if="form.cover_image" class="cover-preview">
              <img :src="form.cover_image" alt="cover" />
              <div class="cover-actions">
                <el-button size="small" type="danger" @click="form.cover_image = ''">移除</el-button>
              </div>
            </div>
            <el-upload
              v-else
              class="cover-uploader"
              action="/api/article/coverImage"
              name="coverImage"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleCoverSuccess"
            >
              <div class="upload-box">
                <el-icon><UploadFilled /></el-icon>
                <div class="upload-text">点击上传封面</div>
              </div>
            </el-upload>
          </div>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="draft">草稿</el-radio>
            <el-radio label="published">已发布</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="正文内容" prop="content">
           <RichTextEditor v-model="form.content" />
        </el-form-item>

        <el-form-item class="footer">
          <el-button @click="router.back()">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="submit('published')">发布文章</el-button>
          <el-button type="primary" plain :loading="submitting" @click="submit('draft')">存为草稿</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import RichTextEditor from '@/components/RichTextEditor.vue'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)

const isEdit = ref(false)
const articleId = ref(null)

const categories = ref([])
const tags = ref([])
const customTags = ref([])
const customTagInput = ref('')

const form = reactive({
  title: '',
  summary: '',
  content: '',
  category_id: null,
  tag_ids: [],
  cover_image: '',
  status: 'published'
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
})

const addCustomTag = () => {
  const name = (customTagInput.value || '').trim()
  if (!name) return
  if (customTags.value.includes(name)) {
    customTagInput.value = ''
    return
  }
  const exists = tags.value.some(t => t.name === name)
  if (exists) {
    ElMessage.info('该标签已存在于候选列表中')
    customTagInput.value = ''
    return
  }
  customTags.value = [...customTags.value, name]
  customTagInput.value = ''
}

const removeCustomTag = (name) => {
  customTags.value = customTags.value.filter(t => t !== name)
}

const handleCoverSuccess = (res) => {
  if (res.code === 200) {
    form.cover_image = res.data
  } else {
    ElMessage.error(res.msg || '上传失败')
  }
}

const fetchData = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([
      request.get('/category/list'),
      request.get('/tag/list')
    ])
    if (catRes.code === 200) categories.value = catRes.data
    if (tagRes.code === 200) tags.value = tagRes.data
  } catch (error) {
    console.error(error)
  }
}

const fetchArticle = async (id) => {
  try {
    const res = await request.get(`/article/${id}`)
    if (res.code === 200) {
      const data = res.data
      form.title = data.title
      form.summary = data.summary
      form.content = data.content
      form.category_id = data.category_id || null
      form.status = data.status
      form.cover_image = data.cover_image || ''
      form.tag_ids = (data.tags || []).map(t => t.id)
      customTags.value = []
    }
  } catch (error) {
    console.error(error)
  }
}

const submit = async (status) => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        let res
        const payload = {
          title: form.title,
          summary: form.summary,
          content: form.content,
          cover_image: form.cover_image,
          category_id: form.category_id || 0,
          tag_ids: form.tag_ids,
          tag_names: customTags.value,
          status: status
        }
        if (isEdit.value) {
          res = await request.put('/article', { ...payload, id: Number(articleId.value) })
        } else {
          res = await request.post('/article', payload)
        }
        if (res.code === 200) {
          ElMessage.success('保存成功')
          router.push('/admin?tab=articles')
        }
      } catch (error) {
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  fetchData()
  const id = route.params.id
  if (id) {
    isEdit.value = true
    articleId.value = id
    fetchArticle(id)
  }
})
</script>

<style scoped>
.w-100 {
  width: 100%;
}
.edit-card {
  border-radius: 12px;
}
.card-header .title {
  font-weight: 700;
  color: #303133;
}
.custom-tag-line {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  width: 100%;
}
.custom-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}
.cover-row {
  display: flex;
  align-items: center;
  gap: 16px;
}
.cover-preview {
  position: relative;
  width: 260px;
  height: 146px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  background: #f3f4f6;
}
.cover-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.cover-actions {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.35);
  opacity: 0;
  transition: opacity 0.2s ease;
}
.cover-preview:hover .cover-actions {
  opacity: 1;
}
.cover-uploader {
  width: 260px;
  height: 146px;
}
.upload-box {
  width: 260px;
  height: 146px;
  border: 2px dashed #d1d5db;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  gap: 10px;
  transition: all 0.2s ease;
}
.upload-box:hover {
  border-color: #409EFF;
  background: #ecf5ff;
  color: #409EFF;
}
.upload-text {
  font-size: 14px;
  font-weight: 600;
}
.footer :deep(.el-form-item__content) {
  justify-content: flex-end;
  gap: 10px;
}
</style>
