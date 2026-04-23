<template>
  <div class="register-wrapper">
    <div class="register-background"></div>
    <div class="register-container">
      <el-card class="register-card" shadow="always">
        <div class="register-header">
          <h2 class="register-title">创建账号</h2>
          <p class="register-subtitle">加入我们的社区，开始您的创作之旅</p>
        </div>
        <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-position="top">
          <el-form-item prop="username">
            <el-input 
              v-model="registerForm.username" 
              placeholder="用户名" 
              size="large"
              :prefix-icon="User"
            ></el-input>
          </el-form-item>
          <el-form-item prop="email">
            <el-input 
              v-model="registerForm.email" 
              placeholder="邮箱" 
              size="large"
              :prefix-icon="Message"
            ></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              type="password" 
              v-model="registerForm.password" 
              placeholder="密码" 
              size="large"
              show-password 
              :prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input 
              type="password" 
              v-model="registerForm.confirmPassword" 
              placeholder="确认密码" 
              size="large"
              show-password 
              :prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <el-form-item prop="code">
            <div class="code-input-group">
              <el-input 
                v-model="registerForm.code" 
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
          <div class="form-actions">
            <el-button type="primary" @click="handleRegister" :loading="loading" class="register-btn" size="large">注 册</el-button>
          </div>
        </el-form>
        <div class="register-footer">
          <router-link to="/login" class="link">已有账号？立即登录</router-link>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Key } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { encryptText, getPublicKey } from '@/utils/rsa'

const router = useRouter()
const registerFormRef = ref(null)
const loading = ref(false)
const countdown = ref(0)

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认密码', trigger: 'blur' }],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const sendCode = async () => {
  if (!registerForm.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  try {
    const res = await request.post('/api/auth/code', { email: registerForm.email })
    if (res.code === 200) {
      ElMessage.success('验证码已发送')
      countdown.value = 60
      const timer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          clearInterval(timer)
        }
      }, 1000)
    }
  } catch (error) {
    console.error(error)
  }
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      if (registerForm.password !== registerForm.confirmPassword) {
        ElMessage.warning('两次密码不一致')
        return
      }
      loading.value = true
      try {
        const { keyId, publicKey } = await getPublicKey()
        const res = await request.post('/api/auth/register', {
          username: registerForm.username,
          email: registerForm.email,
          code: registerForm.code,
          key_id: keyId,
          password: encryptText(publicKey, registerForm.password),
          confirm_password: encryptText(publicKey, registerForm.confirmPassword)
        })
        if (res.code === 200) {
          ElMessage.success('注册成功')
          router.push('/login')
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
.register-wrapper {
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

.register-container {
  width: 100%;
  max-width: 450px;
  padding: 20px;
  z-index: 1;
}

.register-card {
  border: none;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.register-subtitle {
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

.register-btn {
  width: 100%;
  height: 48px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 2px;
  transition: all 0.3s;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.register-footer {
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
  .register-container {
    padding: 15px;
  }
}
</style>
