<template>
  <div class="about-container">
    <el-card class="about-card" shadow="never">
      <div class="about-header">
        <h1>关于本站</h1>
        <p class="subtitle">分享技术，记录生活</p>
      </div>

      <el-divider />

      <div class="about-content">
        <section>
          <h2><el-icon><InfoFilled /></el-icon> 博客简介</h2>
          <p>这是一个基于 Vue 3 和 Go (Gin+Gorm) 开发的个人博客系统。旨在提供一个简洁、高效的平台来发布文章、交流想法。</p>
        </section>

        <section>
          <h2><el-icon><Tools /></el-icon> 我的的技术栈</h2>
          <div class="tech-stack">
            <el-tag size="large">Vue 3</el-tag>
            <el-tag size="large" type="success">Git</el-tag>
            <el-tag size="large" type="warning">Go (Gin)</el-tag>
            <el-tag size="large" type="danger">GORM</el-tag>
            <el-tag size="large" type="info">MySQL</el-tag>
            <el-tag size="large">Redis</el-tag>
            <el-tag size="large">Docker</el-tag>
            <el-tag size="large">Eino</el-tag>
          </div>
        </section>

        <section>
          <h2><el-icon><Message /></el-icon> 联系我</h2>
          <p>如果你有任何建议或合作意向，欢迎通过以下方式联系我：</p>
          <div class="contact-list">
            <div class="contact-item" @click="copyEmail">
              <el-icon class="contact-icon"><Message /></el-icon>
              <span>QQ邮箱：</span>
              <span class="contact-link copyable">{{ email }}</span>
              <el-icon v-if="copied" class="copy-success"><CircleCheck /></el-icon>
            </div>
            <div class="contact-item">
              <el-icon class="contact-icon"><Platform /></el-icon>
              <span>GitHub：</span>
              <a href="https://github.com/rfhdw0102" target="_blank" class="contact-link">rfhdw0102</a>
            </div>
          </div>
        </section>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { InfoFilled, Tools, Message, Platform, CircleCheck } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const email = ref('2129905621@qq.com')
const copied = ref(false)

const copyEmail = async () => {
  try {
    await navigator.clipboard.writeText(email.value)
    copied.value = true
    ElMessage.success('邮箱已复制到剪贴板')

    // 2秒后隐藏对勾图标
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    ElMessage.error('复制失败，请手动复制')
  }
}
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

/* 联系样式 */
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

.email {
  font-family: monospace;
  font-size: 16px;
  background: #f1f5f9;
  padding: 8px 16px;
  border-radius: 8px;
  display: inline-block;
}
</style>