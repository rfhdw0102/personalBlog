<template>
  <div class="login-wrapper">
    <div class="login-background"></div>
    <div class="login-container">
      <el-card class="login-card" shadow="always">
        <div class="login-header">
          <h2 class="login-title">欢迎回来</h2>
          <p class="login-subtitle">请登录您的账号</p>
        </div>
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-position="top">
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="用户名" 
              size="large"
              :prefix-icon="User"
            ></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              type="password" 
              v-model="loginForm.password" 
              placeholder="密码" 
              size="large"
              show-password 
              @keyup.enter="handleLogin"
              :prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="handleLogin" :loading="loading" class="login-btn" size="large">登 录</el-button>
          </div>
        </el-form>
        <div class="login-footer">
          <router-link to="/forgot" class="link">忘记密码？</router-link>
          <div class="divider"></div>
          <router-link to="/register" class="link">没有账号？立即注册</router-link>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { encryptWithServerKey } from '@/utils/rsa'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const { encrypted, keyId } = await encryptWithServerKey(loginForm.password)
        const res = await request.post('/auth/login', {
          username: loginForm.username,
          password: encrypted,
          key_id: keyId
        })
        if (res.code === 200) {
          ElMessage.success('登录成功')
          localStorage.setItem('user', JSON.stringify(res.data))
          // 单独保存 token
          localStorage.setItem('token', res.data.token)
          window.dispatchEvent(new Event('user-login'))
          router.push('/')
        }
      } catch (error) {
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-wrapper {
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

.login-container {
  width: 100%;
  max-width: 420px;
  padding: 20px;
  z-index: 1;
}

.login-card {
  border: none;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
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

.form-actions {
  margin-top: 30px;
}

.login-btn {
  width: 100%;
  height: 48px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 2px;
  transition: all 0.3s;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.login-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 24px;
  font-size: 14px;
}

.link {
  color: #666;
  text-decoration: none;
  transition: color 0.3s;
}

.link:hover {
  color: #409eff;
}

.divider {
  width: 1px;
  height: 14px;
  background-color: #dcdfe6;
  margin: 0 16px;
}

@media (max-width: 480px) {
  .login-container {
    padding: 15px;
  }
  .login-card {
    border-radius: 12px;
  }
}
</style>
