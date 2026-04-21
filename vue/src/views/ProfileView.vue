<template>
  <div class="profile-view">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>个人资料</span>
        </div>
      </template>

      <el-form :model="userForm" :rules="rules" ref="userFormRef" label-width="100px">
        <el-form-item label="头像">
          <el-upload
            class="avatar-uploader"
            action="/api/user/avatar"
            name="avatar"
            :headers="headers"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
          >
            <img v-if="userForm.avatar" :src="avatarUrl" class="avatar" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>

        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="输入新的用户名"></el-input>
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" disabled></el-input>
        </el-form-item>

        <el-form-item label="账户状态">
          <el-tag :type="userForm.status === 1 ? 'success' : 'danger'">
            {{ userForm.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-form-item>

        <el-form-item label="修改密码" prop="password">
          <el-input type="password" v-model="userForm.password" placeholder="留空则不修改" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input type="password" v-model="userForm.confirmPassword" placeholder="确认新密码" show-password></el-input>
        </el-form-item>

        <template v-if="userForm.role === 'admin'">
          <el-divider content-position="left">管理员专有设置</el-divider>
          
          <el-form-item label="个人简介">
            <el-input
              type="textarea"
              :rows="4"
              v-model="userForm.introduction"
              placeholder="输入个人简介，将展示在关于页面或作者栏"
            ></el-input>
          </el-form-item>

          <el-form-item label="收款二维码">
            <el-upload
              class="qr-uploader"
              action="/api/admin/qr"
              name="qr"
              :headers="headers"
              :show-file-list="false"
              :on-success="handleQrSuccess"
              :before-upload="beforeQrUpload"
            >
              <img v-if="userForm.qr" :src="qrUrl" class="qr-code" />
              <el-icon v-else class="qr-uploader-icon"><Plus /></el-icon>
            </el-upload>
            <div class="qr-tip">上传您的收款二维码（打赏功能使用）</div>
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="loading">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { encryptText, getPublicKey } from '@/utils/rsa'

const userFormRef = ref(null)
const loading = ref(false)

const userForm = reactive({
  avatar: '',
  username: '',
  email: '',
  status: 1,
  role: 'user',
  introduction: '',
  qr: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
}

const token = localStorage.getItem('token')
const headers = {
  Authorization: `Bearer ${token}`
}

const avatarUrl = computed(() => {
  if (userForm.avatar) {
    return userForm.avatar.startsWith('http') ? userForm.avatar : 'http://localhost:8082' + userForm.avatar
  }
  return ''
})

const qrUrl = computed(() => {
  if (userForm.qr) {
    return userForm.qr.startsWith('http') ? userForm.qr : 'http://localhost:8082' + userForm.qr
  }
  return ''
})

const fetchUserInfo = async () => {
  try {
    const res = await request.get('/api/user/info')
    if (res.code === 200) {
      const data = res.data
      userForm.avatar = data.avatar
      userForm.username = data.username
      userForm.email = data.email
      userForm.status = data.status
      userForm.role = data.role
      userForm.introduction = data.introduction || ''
      userForm.qr = data.qr || ''
      // sync to localstorage
      localStorage.setItem('user', JSON.stringify(data))
      window.dispatchEvent(new Event('user-login'))
    }
  } catch (error) {
    console.error(error)
  }
}

const handleAvatarSuccess = (res) => {
  if (res.code === 200) {
    userForm.avatar = res.data
    ElMessage.success('上传头像成功')
    fetchUserInfo()
  } else {
    ElMessage.error(res.msg || '上传失败')
  }
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png' || file.type === 'image/gif'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('上传头像图片只能是 JPG/PNG/GIF 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

const handleQrSuccess = (res) => {
  if (res.code === 200) {
    userForm.qr = res.data 
    ElMessage.success('上传二维码成功，请点击"保存修改"完成设置')
  } else {
    ElMessage.error(res.msg || '上传失败')
  }
}

const beforeQrUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png' || file.type === 'image/gif'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('上传二维码只能是 JPG/PNG/GIF 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传二维码大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

const submitForm = async () => {
  if (!userFormRef.value) return
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const payload = {
          username: userForm.username,
          email: userForm.email,
          status: userForm.status,
          introduction: userForm.introduction,
          qr: userForm.qr
        }
        if (userForm.password) {
          if (userForm.password !== userForm.confirmPassword) {
            ElMessage.warning('两次密码不一致')
            loading.value = false
            return
          }
          const { keyId, publicKey } = await getPublicKey()
          payload.key_id = keyId
          payload.password = encryptText(publicKey, userForm.password)
          payload.confirm_password = encryptText(publicKey, userForm.confirmPassword)
        }

        const res = await request.put('/api/user', payload)
        if (res.code === 200) {
          ElMessage.success('修改成功')
          userForm.password = ''
          userForm.confirmPassword = ''
          fetchUserInfo()
        }
      } catch (error) {
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.profile-view {
  max-width: 800px;
  margin: 0 auto;
}
.card-header {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
}
.avatar-uploader :deep(.el-upload) {
  border: 2px dashed #d1d5db;
  border-radius: 50%;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.2s;
  width: 100px;
  height: 100px;
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: var(--primary-color);
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #94a3b8;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.avatar {
  width: 100px;
  height: 100px;
  display: block;
  object-fit: cover;
}
.qr-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.qr-uploader :deep(.el-upload:hover) {
  border-color: var(--primary-color);
}
.qr-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 148px;
  height: 148px;
  line-height: 148px;
  text-align: center;
}
.qr-code {
  width: 148px;
  height: 148px;
  display: block;
  object-fit: contain;
}
.qr-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}
.el-form {
  padding: 20px 0;
}
</style>
