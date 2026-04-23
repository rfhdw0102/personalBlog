<template>
  <div class="forgot-wrapper">
    <div class="forgot-background"></div>
    <div class="forgot-container">
      <el-card class="forgot-card" shadow="always">
        <div class="forgot-header">
          <h2 class="forgot-title">重置密码</h2>
          <p class="forgot-subtitle">请输入您的邮箱以获取验证码</p>
        </div>
        <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
          <el-form-item prop="email">
            <el-input 
              v-model="form.email" 
              placeholder="邮箱" 
              size="large"
              :prefix-icon="Message"
            ></el-input>
          </el-form-item>
          <el-form-item prop="code">
            <div class="code-input-group">
              <el-input 
                v-model="form.code" 
                placeholder="验证码" 
                size="large"
                :prefix-icon="Key"
              ></el-input>
              <el-button 
                @click="sendCode" 
                :disabled="countdown > 0" 
                class="code-btn"
                size="large"
              >
                {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
              </el-button>
            </div>
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              v-model="form.password" 
              type="password" 
              show-password 
              placeholder="新密码" 
              size="large"
              :prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input 
              v-model="form.confirmPassword" 
              type="password" 
              show-password 
              placeholder="确认新密码" 
              size="large"
              :prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="submit" :loading="loading" class="forgot-btn" size="large">重 置 密 码</el-button>
          </div>
        </el-form>
        <div class="forgot-footer">
          <router-link to="/login" class="link">返回登录</router-link>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Message, Lock, Key } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { encryptText, getPublicKey } from '@/utils/rsa'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const countdown = ref(0)

const form = reactive({
  email: '',
  code: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
  password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认新密码', trigger: 'blur' }]
}

const sendCode = async () => {
  if (!form.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  try {
    const res = await request.post('/api/auth/code', { email: form.email })
    if (res.code === 200) {
      ElMessage.success('验证码已发送')
      countdown.value = 60
      const timer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) clearInterval(timer)
      }, 1000)
    }
  } catch (e) {
    console.error(e)
  }
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    if (form.password !== form.confirmPassword) {
      ElMessage.warning('两次密码不一致')
      return
    }
    loading.value = true
    try {
      const { keyId, publicKey } = await getPublicKey()
      const res = await request.post('/api/auth/password', {
        email: form.email,
        code: form.code,
        key_id: keyId,
        password: encryptText(publicKey, form.password),
        confirm_password: encryptText(publicKey, form.confirmPassword)
      })
      if (res.code === 200) {
        ElMessage.success('密码重置成功')
        router.push('/login')
      }
    } catch (e) {
      console.error(e)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.forgot-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;  
}

.forgot-container {
  width: 100%;
  max-width: 450px;
  padding: 20px;
  z-index: 1;
}

.forgot-card {
  border: none;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.forgot-header {
  text-align: center;
  margin-bottom: 30px;
}

.forgot-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.forgot-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
  background-color: #f5f7fa;
  box-shadow: none !important;
  border: 1px solid transparent;
  transition: all 0.3s;
}

:deep(.el-input__wrapper.is-focus) {
  background-color: #fff;
  border-color: #409eff;
}

.code-input-group {
  display: flex;
  gap: 12px;
  width: 100%;
}

.code-btn {
  width: 140px;
  height: 48px;
  border-radius: 8px;
}

.form-actions {
  margin-top: 30px;
}

.forgot-btn {
  width: 100%;
  height: 48px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 2px;
  transition: all 0.3s;
}

.forgot-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.forgot-footer {
  text-align: center;
  margin-top: 24px;
}

.link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.link:hover {
  color: #409eff;
}

@media (max-width: 480px) {
  .forgot-container {
    padding: 15px;
  }
}
</style>
