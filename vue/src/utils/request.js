import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
    baseURL: process.env.VUE_APP_API_BASE_URL || 'http://localhost:8082', // 设置你的API地址
    timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        // 添加 token 到请求头
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data

        // 处理业务错误
        if (res.code !== 200) {
            ElMessage.error(res.message || '请求失败')

            // 如果是401未授权，清除登录状态
            if (res.code === 401) {
                const token = localStorage.getItem('token')
                if (token) {
                    localStorage.removeItem('token')
                    localStorage.removeItem('user')
                    window.dispatchEvent(new CustomEvent('user-logout'))

                    const currentRoute = router.currentRoute.value
                    // 只在需要认证的页面才跳转到登录页
                    if (currentRoute.meta?.requiresAuth) {
                        router.push('/login')
                    }
                }
            }

            return Promise.reject(new Error(res.message || '请求失败'))
        }

        return res
    },
    error => {
        console.error('响应错误:', error)

        // 处理HTTP错误
        if (error.response) {
            const { status, data } = error.response

            switch (status) {
                case 401: {
                    ElMessage.error('登录已过期，请重新登录')
                    localStorage.removeItem('token')
                    localStorage.removeItem('user')
                    window.dispatchEvent(new CustomEvent('user-logout'))

                    const currentRoute = router.currentRoute.value
                    if (currentRoute.meta?.requiresAuth) {
                        router.push('/login')
                    }
                    break
                }
                case 403: {
                    ElMessage.error('没有权限访问')
                    break
                }
                case 404: {
                    ElMessage.error('请求的资源不存在')
                    break
                }
                case 500: {
                    ElMessage.error('服务器错误，请稍后重试')
                    break
                }
                default: {
                    ElMessage.error(data.msg || data.message || '网络错误，请稍后重试')
                }
            }
        } else if (error.request) {
            ElMessage.error('网络连接失败，请检查网络')
        } else {
            ElMessage.error(error.message || '请求失败')
        }

        return Promise.reject(error)
    }
)

export default request