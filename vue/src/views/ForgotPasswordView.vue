<template>
  <div class="forgot-container">
    <el-card class="forgot-card">
      <template #header>
        <h2 class="forgot-title">重置密码</h2>
      </template>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="验证码" prop="code">
          <div class="code-container">
            <el-input v-model="form.code" placeholder="验证码"></el-input>
            <el-button @click="sendCode" :disabled="countdown > 0" class="code-btn">
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入新密码"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password placeholder="请再次输入新密码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit" :loading="loading" class="w-100">提交</el-button>
        </el-form-item>
      </el-form>
      <div class="forgot-links">
        <router-link to="/login">返回登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
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
.forgot-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 120px);
}
.forgot-card {
  width: 420px;
}
.forgot-title {
  margin: 0;
  text-align: center;
}
.w-100 {
  width: 100%;
}
.forgot-links {
  text-align: right;
  margin-top: 10px;
}
.forgot-links a {
  color: #409EFF;
  text-decoration: none;
  font-size: 14px;
}
.code-container {
  display: flex;
  gap: 10px;
  width: 100%;
}
.code-btn {
  width: 110px;
}
</style>
