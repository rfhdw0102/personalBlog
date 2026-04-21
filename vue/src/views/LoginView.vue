<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="login-title">登录</h2>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input type="password" v-model="loginForm.password" placeholder="请输入密码" show-password @keyup.enter="handleLogin"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" class="w-100">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="login-links">
        <router-link to="/forgot">忘记密码</router-link>
        <span class="sep">|</span>
        <router-link to="/register">没有账号？去注册</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
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
        const res = await request.post('/api/auth/login', {
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
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 120px);
}
.login-card {
  width: 400px;
}
.login-title {
  margin: 0;
  text-align: center;
}
.w-100 {
  width: 100%;
}
.login-links {
  text-align: right;
  margin-top: 10px;
}
.sep {
  margin: 0 8px;
  color: #c0c4cc;
}
.login-links a {
  color: #409EFF;
  text-decoration: none;
  font-size: 14px;
}
</style>
